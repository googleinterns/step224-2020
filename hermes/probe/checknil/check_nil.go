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
	"reflect"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/logger"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
	"google.golang.org/api/iterator"

	cpmetrics "github.com/google/cloudprober/metrics"
	sjpb "github.com/googleinterns/step224-2020/hermes/proto"
)

// nilFile ...
// TODO(evanSpendlove): Add doc comment
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
//	- ctx: pass the context so this probe can be cancelled if needed.
//	- target: pass the target run information.
//	- client: pass an initialised storage client for this target system.
//	- logger: pass the logger associated with the probe calling this function.
// Returns:
//	- err:
// TODO(evanSpendlove): Add err details to return comment
// TODO(evanSpendlove): Add exit status to return values.
func CheckNilFile(ctx context.Context, target *probe.Target, client stiface.Client, logger *logger.Logger) (metrics.ExitStatus, error) {
	bucket := target.Target.GetBucketName()

	file, status, err := checkNilFileExists(ctx, bucket, target, client, logger)
	if err != nil {
		return nil, status, err
	}

	file, status, err = readNilFile(ctx, file, bucket, target, client, logger)
	if err != nil {
		return status, err
	}

	if status, err = verifyNilFile(ctx, file, target, client, logger); err != nil {
		return status, err
	}

	logger.Infof("CheckNilFile completed successfully, nil file is valid.")
	return metrics.Success, nil
}

// TODO(evanSpendlove): Add doc comment.
func checkNilFileExists(ctx context.Context, bucket string, target *probe.Target, client stiface.Client, logger *logger.Logger) (*nilFile, metrics.ExitStatus, error) {
	latency, filename, status, err := RecordLatency(func() (interface{}, metrics.ExitStatus, error) {
		objIter := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: "Hermes_Nil"})
		for {
			// TODO(evanSpendlove): Review how Nil file is updated and consider multiple nil files being present.
			obj, err := objIter.Next()
			if err == iterator.Done {
				return nil, metrics.FileMissing, fmt.Errorf("checkNilFile.list_files_failed: unable to find nil file: %w", err)
			}
			if err != nil {
				return nil, metrics.APICallFailed, fmt.Errorf("checkNilFile.list_files_failed: unable to list files: %w", err)
			}
			return obj.Name, metrics.Success, nil
		}
	})

	RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APIListFiles, status, latency)
	name, ok := filename.(string)
	if !ok {
		return nil, metrics.ProbeFailed, fmt.Errorf("failed to assert filename to string")
	}

	logger.Infof("checkNilFile: nil file located on target storage system.")
	return &nilFile{name: name}, status, err
}

// TODO(evanSpendlove): Add doc comment
func readNilFile(ctx context.Context, file *nilFile, bucket string, target *probe.Target, client stiface.Client, logger *logger.Logger) (*nilFile, metrics.ExitStatus, error) {
	fileObject := client.Bucket(bucket).Object(file.name)

	start := time.Now()
	reader, err := fileObject.NewReader(ctx)
	end := time.Now().Sub(start).Seconds()
	defer reader.Close()
	if err != nil {
		RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APIGetFile, metrics.FileReadFailure, end)
		return nil, metrics.FileReadFailure, fmt.Errorf("checkNilFile.read_file_failed: unable to read nil file: %w", err)
	}

	logger.Infof("checkNilFile: nil file retrieved from target storage system.")

	file.contents, err = ioutil.ReadAll(reader)
	if err != nil {
		RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APIGetFile, metrics.FileReadFailure, end)
		return nil, metrics.FileReadFailure, fmt.Errorf("checkNilFile.read_file_failed: unable to read nil file: %w", err)
	}

	RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APIGetFile, metrics.Success, end)
	logger.Infof("checkNilFile: nil file read complete.")

	return file, metrics.Success, nil
}

// TODO(evanSpendlove: Add doc comment
func verifyNilFile(ctx context.Context, file *nilFile, target *probe.Target, client stiface.Client, logger *logger.Logger) (metrics.ExitStatus, error) {
	hash := sha1.New()
	if _, err := io.Copy(hash, file.reader()); err != nil {
		return metrics.ProbeFailed, fmt.Errorf("checkNilFile.compute_checksum_failed: unable to compute checksum of file contents: %w", err)
	}
	file.hash = fmt.Sprintf("%x", hash.Sum(nil))

	if !file.validateChecksum() {
		return metrics.FileCorrupted, fmt.Errorf("checkNilFile.checksum of file did not match filename")
	}

	journal := &sjpb.StateJournal{}
	if err := proto.Unmarshal(file.contents, journal); err != nil {
		return metrics.FileReadFailure, fmt.Errorf("checkNilFile.unmarshal_failed: unable to unmarshal nil file proto: %w", err)
	}

	logger.Infof("checkNilFile: unmarshal complete for nil file proto.")

	if status, err := resolveIntent(ctx, journal, target, client, logger); err != nil {
		return status, err
	}

	if !reflect.DeepEqual(journal, target.Journal) {
		return metrics.FileCorrupted, fmt.Errorf("checkNilFile.nil file contents do not match expected file contents")
	}

	logger.Infof("checkNilFile: nil file contents verified against in-memory struct.")

	return metrics.Success, nil
}

// TODO(evanSpendlove): Add doc comment
func resolveIntent(ctx context.Context, journal *sjpb.StateJournal, target *probe.Target, client stiface.Client, logger *logger.Logger) (metrics.ExitStatus, error) {
	// Read intent, if not nil, check intent validity and update nil file in memory and defer nil file write.
	if journal.Intent != nil {
		fileID, err := strconv.Atoi(journal.Intent.Filename[len("Hermes_"):len("Hermes_ID")])
		if err != nil {
			return metrics.ProbeFailed, err
		}
		switch journal.Intent.FileOperation {
		case sjpb.Intent_CREATE:
			// Check if the file with this filename has already been created
			// Should we check the file system for the file existing?
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
	return metrics.Success, nil
}

// TODO(evanSpendlove): Add doc comment
// TODO(evanSpendlove): Add existing nil file deletion to here
// What do we do if the existing Nil file isn't deleted when Hermes restarts?
func UpdateNilFile(ctx context.Context, target *probe.Target, client stiface.Client, logger *logger.Logger) (metrics.ExitStatus, error) {
	file := &nilFile{}
	var err error

	file.contents, err = proto.Marshal(target.Journal)
	if err != nil {
		return metrics.ProbeFailed, err
	}

	hash := sha1.New()
	if _, err := io.Copy(hash, file.reader()); err != nil {
		return metrics.ProbeFailed, fmt.Errorf("checkNilFile.compute_checksum_failed: unable to compute checksum of file contents: %w", err)
	}

	file.name = fmt.Sprintf("Hermes_Nil_%x", hash)

	bucket := target.Target.GetBucketName()

	logger.Infof("Nil file object created, proceeding to upload.")

	start := time.Now()
	wc := client.Bucket(bucket).Object(file.name).NewWriter(ctx)
	_, err = io.Copy(wc, file.reader())
	end := time.Now().Sub(start).Seconds()
	if err != nil {
		var status metrics.ExitStatus
		switch err {
		case storage.ErrBucketNotExist:
			status = metrics.BucketMissing
		default:
			status = metrics.ProbeFailed
		}
		RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APICreateFile, status, end)
		return metrics.APICallFailed, err
	}

	if err := wc.Close(); err != nil {
		return metrics.APICallFailed, err
	}

	RecordAPILatency(target.LatencyMetrics.APICallLatency, metrics.APICreateFile, metrics.Success, end)

	target.Journal.Intent = nil

	logger.Infof("Nil file update complete.")

	return metrics.Success, nil
}

// TODO(evanSpendlove): Add doc comment
func RecordLatency(f func() (interface{}, metrics.ExitStatus, error)) (float64, interface{}, metrics.ExitStatus, error) {
	start := time.Now()
	obj, status, err := f()
	end := time.Now().Sub(start).Seconds()
	return end, obj, status, err
}

// TODO(evanSpendlove): Add doc comment
func RecordAPILatency(m map[metrics.APICall]map[metrics.ExitStatus]*cpmetrics.EventMetrics, call metrics.APICall, status metrics.ExitStatus, latency float64) {
	m[call][status].Metric("hermes_api_latency_seconds").AddFloat64(latency)
}
