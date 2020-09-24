// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Evan Spendlove, GitHub: evanSpendlove

// Package journal implements the probe operation for checking if the Journal file
// on the target storage system is present. It also verifies the contents of
// the file. The Journal stores the state of Hermes as a map of file IDs to
// filenames. It also stores the intent to complete the next operation.
// This allows Hermes to recover its state if it is interrupted during its
// algorithm.
package journal

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/logger"
	"github.com/google/go-cmp/cmp"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/testing/protocmp"

	cpmetrics "github.com/google/cloudprober/metrics"
	pb "github.com/googleinterns/step224-2020/hermes/proto"
)

const (
	journalFilePrefix = "Hermes_Journal_" // Prefix of the journal filename.
	apiLatency        = "hermes_api_latency_seconds"
)

// UpdateJournal writes the current in-memory journal proto to a new journal file on
// the target system. It then deletes the old journal file located on the system.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
// TODO(evanSpendlove): Add existing journal file deletion to here
func UpdateJournal(ctx context.Context, target *probe.Target, client stiface.Client, logger *logger.Logger) *probe.ProbeError {
	contents, err := proto.Marshal(target.Journal)
	if err != nil {
		return probe.NewProbeError(metrics.ProbeFailed, fmt.Errorf("UpdateJournal(): failed to marshal journal contents: %w", err))
	}

	file := &journalFile{
		contents: contents,
	}

	if err := file.checksumContents(); err != nil {
		return probe.NewProbeError(metrics.ProbeFailed, err)
	}

	file.name = fmt.Sprintf("%s%x", journalFilePrefix, file.checksum)
	bucket := target.Target.GetBucketName()
	logger.Debug("UpdateJournal(): in-memory journal successfully marshalled, proceeding to upload to target storage system.")

	var writer stiface.Writer
	if err := RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APICreateFile, func() *probe.ProbeError {
		writer = client.Bucket(bucket).Object(file.name).NewWriter(ctx)
		length, err := io.Copy(writer, file.reader())
		if writerErr := writer.Close(); writerErr != nil {
			return probe.NewProbeError(metrics.APICallFailed, fmt.Errorf("UpdateJournal(): failed to close writer: %w", writerErr))
		}
		if err != nil {
			switch err {
			case storage.ErrBucketNotExist:
				return probe.NewProbeError(metrics.BucketMissing, err)
			default:
				return probe.NewProbeError(metrics.ProbeFailed, err)
			}
		}
		if length != int64(len(file.contents)) {
			return probe.NewProbeError(metrics.APICallFailed, fmt.Errorf("UpdateJournal(): failed to write correct number of bytes to file, expected %d, got %d", len(file.contents), length))
		}
		return nil
	}); err != nil {
		return err
	}

	logger.Debug("UpdateJournal(): journal file update complete.")
	return nil
}

// CheckJournal checks the existence and validity of the Journal file on the target
// storage system. It also reads the intent from the Journal file to identify if
// Hermes failed to complete the last operation, and completes it if necessary.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
func CheckJournal(ctx context.Context, target *probe.Target, client stiface.Client, logger *logger.Logger) *probe.ProbeError {
	fileName, err := findJournal(ctx, target, client, logger)
	if err != nil {
		return err
	}

	contents, err := readJournal(ctx, fileName, target, client, logger)
	if err != nil {
		return err
	}

	j := &journalFile{
		name:     fileName,
		contents: contents,
	}

	if err := j.checksumContents(); err != nil {
		return probe.NewProbeError(metrics.ProbeFailed, err)
	}

	journal, err := j.unmarshalContents()
	if err != nil {
		return err
	}

	if err := resolveIntent(ctx, journal, target, client, logger); err != nil {
		return err
	}

	if err = verifyJournal(ctx, j, journal, target, client, logger); err != nil {
		return err
	}

	logger.Debug("CheckJournal() completed successfully, journal file is valid.")
	return nil
}

// findJournal checks that the journal file exists on the target system
// and returns its filename if it does.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- fileName: the file name of the Journal file on the target system.
//	- err: error exit statuses are well defined in metrics.go
func findJournal(ctx context.Context, target *probe.Target, client stiface.Client, logger *logger.Logger) (string, *probe.ProbeError) {
	bucket := target.Target.GetBucketName()
	var fileName *string

	err := RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APIListFiles, func() *probe.ProbeError {
		objIter := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: journalFilePrefix})
		// TODO(evanSpendlove): Review how journal file is updated and consider multiple journal files being present.
		// 		If one, success
		// 		If two -> verify whih one reflects the current state of the system and delete the other.
		// 		If two and only differ by intent, prefer the one with intent.
		// TODO(evanSpendlove): Look into adding a datetime to the journal proto.
		obj, err := objIter.Next()
		switch err {
		case iterator.Done:
			return probe.NewProbeError(metrics.FileMissing, fmt.Errorf("findJournal(): unable to find journal file: %w", err))
		case nil:
			break
		default:
			return probe.NewProbeError(metrics.APICallFailed, fmt.Errorf("findJournal(): unable to list files: %w", err))
		}
		*fileName = obj.Name
		return nil
	})

	logger.Debug("findJournal(): journal file located on target storage system.")
	return *fileName, err
}

// readJournal returns the file contents of the journal on the target system.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- fileName: the file name of the Journal file on the target system.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- []byte: file contents read from the journal on the target system.
//	- err: error exit statuses are well defined in metrics.go.
func readJournal(ctx context.Context, filename string, target *probe.Target, client stiface.Client, logger *logger.Logger) ([]byte, *probe.ProbeError) {
	bucket := target.Target.GetBucketName()
	fileObject := client.Bucket(bucket).Object(filename)

	var reader stiface.Reader
	err := RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APIGetFile, func() *probe.ProbeError {
		r, err := fileObject.NewReader(ctx)
		if err != nil {
			return probe.NewProbeError(metrics.FileReadFailure, fmt.Errorf("readJournal(): unable to read journal file: %w", err))
		}
		reader = r
		return nil
	})
	defer reader.Close()
	if err != nil {
		return nil, err
	}

	logger.Debug("readJournal(): journal file retrieved from target storage system.")

	contents, readErr := ioutil.ReadAll(reader)
	if readErr != nil {
		return nil, probe.NewProbeError(metrics.FileReadFailure, fmt.Errorf("readJournal(): unable to read journal file: %w", readErr))
	}

	logger.Debug("readJournal(): journal file read complete.")
	return contents, nil
}

// resolveIntent verifies and completes the intent stored in the state journal passed.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- journal: the StateJournal read from the storage system.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
func resolveIntent(ctx context.Context, journal *pb.StateJournal, target *probe.Target, client stiface.Client, logger *logger.Logger) *probe.ProbeError {
	const (
		hermesPrefix   = "Hermes_"
		hermesIDPrefix = "Hermes_ID"
	)

	if journal.Intent == nil {
		return nil
	}

	fileID, err := strconv.Atoi(journal.Intent.Filename[len(hermesPrefix):len(hermesIDPrefix)])
	if err != nil {
		return probe.NewProbeError(metrics.ProbeFailed, fmt.Errorf("resolveIntent(): failed to parse int ID from file name %q: %w", journal.Intent.Filename, err))
	}

	switch journal.Intent.FileOperation {
	case pb.Intent_CREATE:
		// TODO(evanSpendlove): Check if the file with this filename has already been created. Check the filesystem for the file existing and check the filename is valid.
		if journal.Filenames[int32(fileID)] != journal.Intent.Filename {
			// TODO(evanSpendlove): Call createFile() passing the fileID and return
			logger.Debug("File create not complete yet")
			journal.Filenames[int32(fileID)] = journal.Intent.Filename
		}
	case pb.Intent_DELETE:
		if name := journal.Filenames[int32(fileID)]; name != "" {
			// TODO(evanSpendlove): Call  DeleteFile() passing the fileID and return
			logger.Debug("File delete not complete yet")
			delete(journal.Filenames, int32(fileID))
		}
	default:
		return probe.NewProbeError(metrics.ProbeFailed, fmt.Errorf("resolveIntent(): unexpected FileOperation parsed as part of intent: %v", journal.Intent.FileOperation))
	}
	journal.Intent = nil
	return nil
}

// verifyJournal verifies that the contents of the journal file match its checksum
// and that the journal matches the existing in-memory proto.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- file: a journalFile struct with the file name, contents and checksum recorded.
//  - journal: unmarshalled journal struct, to be verified, read from the target storage system.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
func verifyJournal(ctx context.Context, file *journalFile, journal *pb.StateJournal, target *probe.Target, client stiface.Client, logger *logger.Logger) *probe.ProbeError {
	if err := file.validateChecksum(); err != nil {
		return err
	}
	logger.Debug("verifyJournal(): journal contents checksum verified against filename.")

	if diff := cmp.Diff(target.Journal, journal, protocmp.Transform()); diff != "" {
		return probe.NewProbeError(metrics.FileCorrupted, fmt.Errorf("verifyJournal(): in-memory journal does not match the journal on the target system (-want +got):\n%s", diff))
	}
	logger.Debug("verifyJournal(): journal struct contents verified against in-memory struct.")
	return nil
}

// RecordAPILatency records the time taken for the function passed to complete.
// It then stores this in the metrics map passed using the provided labels.
// Arguments:
//	- m: API latency metric map for this target.
//	- call: metric label for this API call
//	- f: target function to be timed.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
// TODO(evanSpendlove): Migrate this into metrics.go and add a RecordProbeLatency() func there.
func RecordAPILatency(m map[metrics.APICall]map[metrics.ExitStatus]*cpmetrics.EventMetrics, call metrics.APICall, f func() *probe.ProbeError) *probe.ProbeError {
	start := time.Now()
	err := f()
	m[call][err.Status].Metric(apiLatency).AddFloat64(time.Now().Sub(start).Seconds())
	return err
}

// journalFile stores the name, checksum and contents of the Journal.
type journalFile struct {
	name     string
	checksum []byte
	contents []byte
}

// checksumContents generates the checksum of the contents of the file object and stores
// the checksum in the checksum field of the object.
// Returns:
// 	- err: returned if there is an error when checksumming the contents.
func (j *journalFile) checksumContents() error {
	h := sha1.New()

	length, err := io.Copy(h, j.reader())
	if err != nil {
		return fmt.Errorf("checksumContents(): unable to compute checksum of file contents: %w", err)
	}
	if length != int64(len(j.contents)) {
		return fmt.Errorf("checksumContents(): failed to generate checksum of all of the journal file contents")
	}

	j.checksum = h.Sum(nil)
	return nil
}

// validateChecksum validates that the checksum matches the checksum
// in the file name.
// The checksum in the filename is being validated against the file contents.
// This will return false if the file contents have been updated
// without the file name being changed accordingly.
// Returns:
//	- err: returns err if the checksum does not match or the checksum or filename is blank, else nil.
func (j *journalFile) validateChecksum() *probe.ProbeError {
	if len(j.checksum) == 0 {
		return probe.NewProbeError(metrics.ProbeFailed, fmt.Errorf("validateChecksum(): journal checksum cannot be validated as checksum, %q, is blank", j.checksum))
	}
	if j.name == "" {
		return probe.NewProbeError(metrics.ProbeFailed, fmt.Errorf("validateChecksum(): journal checksum cannot be validated as name, %q, is blank", j.name))
	}

	hexChecksum := fmt.Sprintf("%x", j.checksum)
	if hexChecksum != j.name[len(journalFilePrefix):] {
		return probe.NewProbeError(metrics.ProbeFailed, fmt.Errorf("validateChecksum(): checksum %q does not match checksum in name %q", j.checksum, j.name[len(journalFilePrefix):]))
	}
	return nil
}

// reader returns an io.Reader that in turn returns the file contents.
// Returns:
// - reader: returns an io.Reader that in turn returns the file contents.
func (j *journalFile) reader() io.Reader {
	return bytes.NewReader(j.contents)
}

// unmarshalContents unmarshals the contents of the journal file and returns
// the completed proto struct.
// Returns:
// 	- journal: proto struct containing the unmarshalled contents of the journal.
//	- err: error if once occurs during the unmarshalling process.
func (j *journalFile) unmarshalContents() (*pb.StateJournal, *probe.ProbeError) {
	journal := &pb.StateJournal{}
	if err := proto.Unmarshal(j.contents, journal); err != nil {
		return nil, probe.NewProbeError(metrics.FileCorrupted, fmt.Errorf("unmarshalContents(): unable to unmarshal journal file proto: %w", err))
	}
	return journal, nil
}
