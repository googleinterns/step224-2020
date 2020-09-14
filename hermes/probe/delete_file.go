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
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/googleinterns/step224-2020/hermes"
	"google.golang.org/api/iterator"

	probepb "github.com/googleinterns/step224-2020/config/proto"
)

// deleteRandomFile deletes a random file in the target storage system bucket.
// It then checks that the file has been deleted by trying to get the object.
// Arguments:
//	- ctx: pass the context so this probe can be cancelled if needed.
//	- mp: pass the monitorProbe object calling this function.
//	- target: pass the target to be probed.
//	- client: pass an initialised storage client for this target system.
// Returns:
//	- fileID: returns the ID of the file delete OR a missing file to be created if one is found.
//	- err:
//		- [...].hermes_file_missing: the file to be deleted does not exist in Hermes' StateJournal.
//		- [...].file_missing: the file to be deleted could not be found in the target bucket.
//		- [...].bucket_missing: the target bucket on this target system was not found.
//		- [...].probe_failed: there was an error during one of the API calls and the probe failed.
//		- [...].deleted_file_found: the file was deleted but it is still found in the target bucket.
//		- [...].list_bucket_failed: the listBucket operation failed when checking if the target file was deleted.
func deleteRandomFile(ctx context.Context, mp *MonitorProbe, target *probepb.Target, client *stiface.Client) (int, error) {
	bucket := target.GetBucketName()

	fileID := hermes.PickFileToDelete()

	filename, ok := mp.Hermes.Journal.Filenames[int32(fileID)]
	if !ok {
		return fileID, fmt.Errorf("error, deleteRandomFile(id: %d).hermes_file_missing: could not delete file %s, file not found", fileID, filename)
	}

	file := (*client).Bucket(bucket).Object(filename)

	start := time.Now()
	if err := file.Delete(ctx); err != nil {
		var status string
		if err == storage.ErrObjectNotExist {
			status = "file_missing"
		} else if err == storage.ErrBucketNotExist {
			status = "bucket_missing"
		} else {
			status = "probe_failed"
		}
		mp.LatencyMetrics[target.GetName()].ApiCallLatency["delete"][status].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
		return fileID, fmt.Errorf("error, deleteRandomFile(id: %d).%s: could not delete file %s: %v", fileID, status, filename, err)
	}
	mp.LatencyMetrics[target.GetName()].ApiCallLatency["delete"]["success"].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())

	query := &storage.Query{Prefix: filename}
	start = time.Now()
	it := (*client).Bucket(bucket).Objects(ctx, query)
	mp.LatencyMetrics[target.GetName()].ApiCallLatency["list"]["success"].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
	for {
		obj, err := it.Next()
		if obj.Name == filename {
			return fileID, fmt.Errorf("error, deleteRandomFile(id %d).deleted_file_found: object %v in bucket %q still listed after delete", fileID, *obj, bucket)
		}
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fileID, fmt.Errorf("error, deleteRandomFile(id: %d).list_bucket_failed: unable to list bucket %q: %v", fileID, bucket, err)
		}
	}

	mp.Logger.Infof("Object %v deleted in bucket %s.", file, bucket)
	return fileID, nil
}