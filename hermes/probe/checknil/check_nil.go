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
//
// Check_nil implements the probe operation for checking if the NIL file is
// present and reading verifying its contents.

package checknil

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
	sjpb "github.com/googleinterns/step224-2020/hermes/proto"
)

// nilFile stores the name, hash and contents of the nil file.
// It is used when checking or updating the nil file.
type nilFile struct {
	name     string
	hash     string
	contents []byte
}

// validateChecksum validates that the hash matches the checksum
// in the file name.
// Returns:
//	- bool: returns true if the hash matches, else false.
func (n *nilFile) validateChecksum() bool {
	if n.hash == "" || n.name == "" {
		return false
	}

	// Here the checksum in the filename is being validated against the file contents
	// This will return false if the file contents have been updated
	// without the file name being changed accordingly.
	return n.hash == n.name[len("Hermes_NIL_"):]
}

// reader returns the file contents as a reader.
// Returns:
// - reader: returns the file contents as a reader.
func (n *nilFile) reader() io.Reader {
	return bytes.NewReader(n.contents)
}

// CheckNilFile checks the existence and validity of the Nil file on the target
// storage system. It also reads the intent from the Nil file to identify if
// Hermes failed to complete the last operation, and completes it if necessary.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
func CheckNilFile(ctx context.Context, target *probe.Target, client stiface.Client, logger *logger.Logger) *probe.ProbeError {
	bucket := target.Target.GetBucketName()

	file, err := checkNilFileExists(ctx, bucket, target, client, logger)
	if err.Status != metrics.Success {
		return err
	}

	if file, err = readNilFile(ctx, file, bucket, target, client, logger); err.Status != metrics.Success {
		return err
	}

	if err = verifyNilFile(ctx, file, target, client, logger); err.Status != metrics.Success {
		return err
	}

	logger.Infof("CheckNilFile() completed successfully, nil file is valid.")
	return probe.NewProbeError(metrics.Success, nil)
}

// checkNilFileExists checks that the nil file exists on the target system.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- nilFile: a nilFile struct with the file name recorded
//	- err: error exit statuses are well defined in metrics.go
func checkNilFileExists(ctx context.Context, bucket string, target *probe.Target, client stiface.Client, logger *logger.Logger) (*nilFile, *probe.ProbeError) {
	file := &nilFile{}

	err := RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APIListFiles, func() *probe.ProbeError {
		objIter := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: "Hermes_Nil"})
		for {
			// TODO(evanSpendlove): Review how Nil file is updated and consider multiple nil files being present.
			// If one, success
			// If two -> verify whih one reflects the current state of the system and delete the other.
			// If two and only differ by intent, prefer the one with intent.
			obj, err := objIter.Next()
			if err == iterator.Done {
				return probe.NewProbeError(metrics.FileMissing, fmt.Errorf("checkNilFileExists(): unable to find nil file: %w", err))
			}
			if err != nil {
				return probe.NewProbeError(metrics.APICallFailed, fmt.Errorf("checkNilFileExists(): unable to list files: %w", err))
			}
			file.name = obj.Name
			return probe.NewProbeError(metrics.Success, nil)
		}
	})

	logger.Infof("checkNilFileExists(): nil file located on target storage system.")
	return file, err
}

// checkNilFileExists checks that the nil file exists on the target system.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- nilFile: a nilFile struct with the file name recorded
//	- err: error exit statuses are well defined in metrics.go
func readNilFile(ctx context.Context, file *nilFile, bucket string, target *probe.Target, client stiface.Client, logger *logger.Logger) (*nilFile, *probe.ProbeError) {
	fileObject := client.Bucket(bucket).Object(file.name)
	var reader stiface.Reader

	err := RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APIGetFile, func() *probe.ProbeError {
		var err error
		reader, err = fileObject.NewReader(ctx)
		if err != nil {
			return probe.NewProbeError(metrics.FileReadFailure, fmt.Errorf("readNilFile(): unable to read nil file: %w", err))

		}
		return probe.NewProbeError(metrics.Success, nil)
	})
	defer reader.Close()

	if err.Status != metrics.Success {
		return nil, err
	}

	logger.Infof("readNilFile(): nil file retrieved from target storage system.")

	var readErr error
	file.contents, readErr = ioutil.ReadAll(reader)
	if readErr != nil {
		return nil, probe.NewProbeError(metrics.FileReadFailure, fmt.Errorf("readNilFile(): unable to read nil file: %w", readErr))
	}

	logger.Infof("readNilFile(): nil file read complete.")

	return file, probe.NewProbeError(metrics.Success, nil)
}

// verifyNilFile verifies that the contents of the nil file match its checksum
// and that the journal matches the existing in-memory proto.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- nilFile: a nilFile struct with the file name and contents recorded
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
func verifyNilFile(ctx context.Context, file *nilFile, target *probe.Target, client stiface.Client, logger *logger.Logger) *probe.ProbeError {
	hash := sha1.New()
	if _, err := io.Copy(hash, file.reader()); err != nil {
		return probe.NewProbeError(metrics.ProbeFailed, fmt.Errorf("verifyNilFile(): unable to compute checksum of file contents: %w", err))
	}
	file.hash = fmt.Sprintf("%x", hash.Sum(nil))

	if !file.validateChecksum() {
		return probe.NewProbeError(metrics.FileCorrupted, fmt.Errorf("verifyNilFile(): checksum of nil file did not match filename"))
	}

	journal := &sjpb.StateJournal{}
	if err := proto.Unmarshal(file.contents, journal); err != nil {
		return probe.NewProbeError(metrics.FileReadFailure, fmt.Errorf("verifyNilFile(): unable to unmarshal nil file proto: %w", err))
	}

	logger.Infof("verifyNilFile(): unmarshal complete for nil file proto.")

	if err := resolveIntent(ctx, journal, target, client, logger); err != nil {
		return err
	}

	if diff := cmp.Diff(target.Journal, journal, protocmp.Transform()); diff != "" {
		return probe.NewProbeError(metrics.FileCorrupted, fmt.Errorf("verifyNilFile(): journal contents mismatch (-want +got):\n%s", diff))
	}

	logger.Infof("verifyNilFile(): nil file contents verified against in-memory struct.")

	return probe.NewProbeError(metrics.Success, nil)
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
func resolveIntent(ctx context.Context, journal *sjpb.StateJournal, target *probe.Target, client stiface.Client, logger *logger.Logger) *probe.ProbeError {
	if journal.Intent != nil {
		fileID, err := strconv.Atoi(journal.Intent.Filename[len("Hermes_"):len("Hermes_ID")])
		if err != nil {
			return probe.NewProbeError(metrics.ProbeFailed, err)
		}
		switch journal.Intent.FileOperation {
		case sjpb.Intent_CREATE:
			// TODO(evanSpendlove): Check if the file with this filename has already been created. Should we check the filesystem for the file existing?
			if journal.Filenames[int32(fileID)] != journal.Intent.Filename {
				// TODO(evanSpendlove): Call createFile() passing the fileID and return
				logger.Info("File create not complete yet")
				journal.Filenames[int32(fileID)] = journal.Intent.Filename
			}
		case sjpb.Intent_DELETE:
			if _, ok := journal.Filenames[int32(fileID)]; ok {
				// TODO(evanSpendlove): Call  DeleteFile() passing the fileID and return
				logger.Info("File delete not complete yet")
				delete(journal.Filenames, int32(fileID))
			}
		}
		journal.Intent = nil
	}
	return probe.NewProbeError(metrics.Success, nil)
}

// WriteNilFile writes the current in-memory journal proto a new nil file on
// the target system. It then deletes the old nil file located on the system.
// Arguments:
//	- ctx: context so this probe can be cancelled if needed.
//	- target: target run information.
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
// TODO(evanSpendlove): Add existing nil file deletion to here
func UpdateNilFile(ctx context.Context, target *probe.Target, client stiface.Client, logger *logger.Logger) *probe.ProbeError {
	file := &nilFile{}
	var err error

	file.contents, err = proto.Marshal(target.Journal)
	if err != nil {
		return probe.NewProbeError(metrics.ProbeFailed, err)
	}

	hash := sha1.New()
	if _, err := io.Copy(hash, file.reader()); err != nil {
		return probe.NewProbeError(metrics.ProbeFailed, fmt.Errorf("UpdateNilFile: unable to compute checksum of file contents: %w", err))
	}

	file.name = fmt.Sprintf("Hermes_Nil_%x", hash.Sum(nil))

	bucket := target.Target.GetBucketName()

	logger.Infof("UpdateNilFile(): Nil file object created, proceeding to upload.")

	var writer stiface.Writer

	if err := RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APICreateFile, func() *probe.ProbeError {
		writer = client.Bucket(bucket).Object(file.name).NewWriter(ctx)
		_, err = io.Copy(writer, file.reader())
		if err != nil {
			switch err {
			case storage.ErrBucketNotExist:
				return probe.NewProbeError(metrics.BucketMissing, err)
			default:
				return probe.NewProbeError(metrics.ProbeFailed, err)
			}
			return probe.NewProbeError(metrics.APICallFailed, err)
		}
		return probe.NewProbeError(metrics.Success, nil)
	}); err.Status != metrics.Success {
		return err
	}

	if err := writer.Close(); err != nil {
		return probe.NewProbeError(metrics.APICallFailed, err)
	}

	target.Journal.Intent = nil

	logger.Infof("UpdateNilFile(): Nil file update complete.")

	return probe.NewProbeError(metrics.Success, nil)
}

// RecordAPILatency records the time taken for the function passed to complete.
// It then stores this in the metrics map passed using the provided labels.
// Arguments:
//	- m: api latency metric map for this target.
//	- call: metric label for this API call
//	- f: target function to be timed.
// Returns:
//	- err: error exit statuses are well defined in metrics.go
// TODO(evanSpendlove): Migrate this into metrics.go and add a RecordProbeLatency() func there.
func RecordAPILatency(m map[metrics.APICall]map[metrics.ExitStatus]*cpmetrics.EventMetrics, call metrics.APICall, f func() *probe.ProbeError) *probe.ProbeError {
	start := time.Now()
	err := f()
	m[call][err.Status].Metric("hermes_api_latency_seconds").AddFloat64(time.Now().Sub(start).Seconds())
	return err
}
