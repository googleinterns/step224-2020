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
// Delete_file implements the probe operation for deleting a file in a
// storage system.
// Package probe implements the probe that Hermes uses to monitor
// a storage system.

package probe

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/cloudprober/logger"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/googleinterns/step224-2020/hermes"
	"google.golang.org/api/iterator"

	probepb "github.com/googleinterns/step224-2020/config/proto"
	m "github.com/googleinterns/step224-2020/hermes/probe/metrics"
)

// deleteRandomFile deletes a random file in the target storage system bucket.
// It then checks that the file has been deleted by trying to get the object.
// Arguments:
// TODO(evanSpendlove): Update comment here
//	- ctx: pass the context so this probe can be cancelled if needed.
//	- config: pass the HermesProbeDef config for the probe calling this function.
//	- target: pass the target run information.
//	- client: pass an initialised storage client for this target system.
//	- logger: pass the logger associated with the probe calling this function.
// Returns:
//	- fileID: returns the ID of the file delete OR a missing file to be created if one is found.
//	- err:
//		- [...].hermes_file_missing: the file to be deleted does not exist in Hermes' StateJournal.
//		- [...].file_missing: the file to be deleted could not be found in the target bucket.
//		- [...].bucket_missing: the target bucket on this target system was not found.
//		- [...].probe_failed: there was an error during one of the API calls and the probe failed.
//		- [...].deleted_file_found: the file was deleted but it is still found in the target bucket.
//		- [...].list_bucket_failed: the listBucket operation failed when checking if the target file was deleted.
func DeleteRandomFile(ctx context.Context, config *probepb.HermesProbeDef, target *Target, client *stiface.Client, logger *logger.Logger) (int, error) {
	bucket := target.Target.GetBucketName()

	fileID := hermes.PickFileToDelete()

	filename, ok := target.Journal.Filenames[int32(fileID)]
	if !ok {
		return fileID, fmt.Errorf("deleteRandomFile(id: %d).hermes_file_missing: could not delete file %s, file not found", fileID, filename)
	}

	file := (*client).Bucket(bucket).Object(filename)

	start := time.Now()
	if err := file.Delete(ctx); err != nil {
		var status string
		if err == storage.ErrObjectNotExist {
			status = m.ExitStatuses[m.FILE_MISSING]
		} else if err == storage.ErrBucketNotExist {
			status = m.ExitStatuses[m.BUCKET_MISSING]
		} else {
			status = m.ExitStatuses[m.PROBE_FAILED]
		}
		target.LatencyMetrics.ApiCallLatency[m.ApiCalls[m.API_DELETE_FILE]][status].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
		return fileID, fmt.Errorf("deleteRandomFile(id: %d).%s: could not delete file %s: %w", fileID, status, filename, err)
	}
	target.LatencyMetrics.ApiCallLatency[m.ApiCalls[m.API_DELETE_FILE]][m.ExitStatuses[m.SUCCESS]].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())

	query := &storage.Query{Prefix: filename}
	start = time.Now()
	it := (*client).Bucket(bucket).Objects(ctx, query)
	target.LatencyMetrics.ApiCallLatency[m.ApiCalls[m.API_LIST_FILES]][m.ExitStatuses[m.SUCCESS]].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
	for {
		obj, err := it.Next()
		if obj.Name == filename {
			return fileID, fmt.Errorf("deleteRandomFile(id %d).deleted_file_found: object %v in bucket %q still listed after delete", fileID, *obj, bucket)
		}
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fileID, fmt.Errorf("deleteRandomFile(id: %d).list_bucket_failed: unable to list bucket %q: %w", fileID, bucket, err)
		}
	}

	// Update in-memory NIL file after delete operation.
	delete(target.Journal.Filenames, int32(fileID))

	logger.Infof("Object %v deleted in bucket %s.", file, bucket)
	return fileID, nil
}
