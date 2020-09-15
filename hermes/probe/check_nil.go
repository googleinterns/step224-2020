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
// Package probe implements the probe that Hermes uses to monitor
// a storage system.

package probe

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/logger"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"

	m "github.com/googleinterns/step224-2020/hermes/probe/metrics"
	sjpb "github.com/googleinterns/step224-2020/hermes/proto"
)

// TODO(evanSpendlove): Update comment here
// deleteRandomFile deletes a random file in the target storage system bucket.
// It then checks that the file has been deleted by trying to get the object.
// Arguments:
//	- ctx: pass the context so this probe can be cancelled if needed.
//	- target: pass the target run information.
//	- client: pass an initialised storage client for this target system.
//	- logger: pass the logger associated with the probe calling this function.
// Returns:
//	- err:
//		- [...].hermes_file_missing: the file to be deleted does not exist in Hermes' StateJournal.
//		- [...].file_missing: the file to be deleted could not be found in the target bucket.
//		- [...].bucket_missing: the target bucket on this target system was not found.
//		- [...].probe_failed: there was an error during one of the API calls and the probe failed.
//		- [...].deleted_file_found: the file was deleted but it is still found in the target bucket.
//		- [...].list_bucket_failed: the listBucket operation failed when checking if the target file was deleted.
// TODO(evanSpendlove): Add exit status to return values.
func CheckNilFile(ctx context.Context, target *Target, client *stiface.Client, logger *logger.Logger) (m.ExitStatus, error) {
	bucket := target.Target.GetBucketName()

	// TODO(evanSpendlove): Find the Nil file on the storage system, alert if not present.
	query := &storage.Query{Prefix: "Hermes_NIL"}
	start := time.Now()
	objIter := (*client).Bucket(bucket).Objects(ctx, query) // Can add API call timeouts here.
	end := time.Now()
	var nilFilename string
	for {
		obj, err := objIter.Next()
		if err != nil {
			target.LatencyMetrics.ApiCallLatency[m.ApiCalls[m.API_LIST_FILES]][m.ExitStatuses[m.FILE_MISSING]].Metric("latency").AddFloat64(end.Sub(start).Seconds())
			return fmt.Errorf("checkNilFile.list_files_failed: unable to find nil file: %w", err)
		}
		nilFilename = obj.Name
	}
	target.LatencyMetrics.ApiCallLatency[m.ApiCalls[m.API_LIST_FILES]][m.ExitStatuses[m.SUCCESS]].Metric("latency").AddFloat64(end.Sub(start).Seconds())
	logger.Infof("checkNilFile: nil file located on target storage system.")

	// TODO(evanSpendlove): Get the contents of the NIL file.
	nilFile := (*client).Bucket(bucket).Object(nilFilename)

	start = time.Now()
	reader, err := nilFile.NewReader(ctx)
	end = time.Now()
	if err != nil {
		target.LatencyMetrics.ApiCallLatency[m.ApiCalls[m.API_GET_FILE]][m.ExitStatuses[m.FILE_READ_FAILURE]].Metric("latency").AddFloat64(end.Sub(start).Seconds())
		return fmt.Errorf("checkNilFile.read_file_failed: unable to read nil file: %w", err)
	}
	target.LatencyMetrics.ApiCallLatency[m.ApiCalls[m.API_GET_FILE]][m.ExitStatuses[m.SUCCESS]].Metric("latency").AddFloat64(end.Sub(start).Seconds())
	logger.Infof("checkNilFile: nil file gotten from target storage system.")

	// TODO(evanSpendlove): Amend metric recording here to use api_READ or probe op read and verify metrics.
	start = time.Now()
	nilContents, err := ioutil.ReadAll(reader)
	if err != nil {
		target.LatencyMetrics.ProbeOpLatency[m.ProbeOps[m.READ_FILE]][m.ExitStatuses[m.FILE_READ_FAILURE]].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errorf("checkNilFile.read_file_failed: unable to read nil file: %w", err)
	}

	logger.Infof("checkNilFile: nil file read complete.")

	// Verify checksum
	// TODO(evanSpendlove): Re-evaluate how best to time verifying a file, possibly adding hashing time to other verification.
	hashTime := time.Now()
	hash := sha1.New()
	if _, err := io.Copy(hash, reader); err != nil {
		target.LatencyMetrics.ProbeOplatency[m.ProbeOps[m.VERIFY_FILE]][m.ExitStatuses[m.FILE_CORRUPTED]].Metric("latency").AddFloat64(time.Now().Sub(hashTime).Seconds())
		return nil, fmt.Errorf("checkNilFile.compute_checksum_failed: unable to compute checksum of file contents: %w", err)
	}
	if fmt.Sprintf("%x", hash) != nilFilename[len("Hermes_NIL_"):] {
		target.LatencyMetrics.ProbeOpLatency[m.ProbeOps[m.VERIFY_FILE]][m.ExitStatuses[m.FILE_CORRUPTED]].Metric("latency").AddFloat64(time.Now().Sub(hashTime).Seconds())
		return fmt.Errorf("checkNilFile.verify_checksum_failed: checksum of file did not match filename: %w", err)
	}
	hashTime = time.Now().Sub(hashTime)
	reader.Close()

	// TODO(evanSpendlove): Decode the contents of the NIL file.
	journal := &sjpb.Journal{}
	if err := proto.Unmarshal(nilContents, journal); err != nil {
		target.LatencyMetrics.ProbeOpLatency[m.ProbeOps[m.READ_FILE]][m.ExitStatuses[m.FILE_CORRUPTED]].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errrof("checkNilFile.unmarshal_failed: unable to unmarshal nil file proto: %w", err)
	}
	target.LatencyMetrics.ProbeOpLatency[m.ProbeOps[m.READ_FILE]][m.ExitStatuses[m.SUCCESS]].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
	logger.Infof("checkNilFile: unmarshal complete for nil file proto.")

	// TODO(evanSpendlove): Verify the contents against in memory and checksum.
	start = time.Now() // Start for verify op

	// TODO(evanSpendlove): Handle intent from NIL file.
	// Read intent, if not nil, check intent validity and update nil file in memory and defer nil file write.
	if journal.Intent != nil {
		// Check validity
		fileID := journal.Intent.Filename[len("Hermes_"):len("Hermes_ID")]
		switch journal.Intent.FileOperation {
		case sjpb.Intent_CREATE:
			// Check if the file with this filename has already been created
			// TODO: Should we check the file system for the file existing?
			if target.Journal.Filenames[fileID] != Journal.Intent.Filename {
				// TODO(evanSpendlove): Call createFile() passing the fileID
				logger.Info("File create not complete yet")
			} else {
				journal.Intent = nil // Remove intent as operation performed in full.
			}
		case sjpb.Intent_DELETE:
			if _, ok := target.Journal.Filenames[fileID]; !ok {
				// TODO(evanSpendlove): Call  DeleteFile() passing the fileID
				logger.Info("File create not complete yet")
			} else {
				journal.Intent = nil // Remove intent as operation performed in full.
			}
		}
	}

	if !reflect.DeepEqual(journal, target.Journal) {
		target.LatencyMetrics.ProbeOpLatency[m.ProbeOps[m.VERIFY_FILE]][m.ExitStatuses[m.FILE_CORRUPTED]].Metric("latency").AddFloat64(time.Now().Add(hashTime).Sub(start).Seconds())
		return fmt.Errorf("checkNilFile.file_inconsistent: nil file contents do not match expectations: %w", err)
	}

	// TODO(evanSpendlove): Add logging statements.
	logger.Infof("Object %v deleted in bucket %s.", file, bucket)
	return fileID, nil
}
