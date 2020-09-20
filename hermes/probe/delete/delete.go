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

// Package delete implements the file deletion operation with a storage system.
package delete

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/cloudprober/logger"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
	"google.golang.org/api/iterator"

	pb "github.com/googleinterns/step224-2020/hermes/proto"
)

const (
	apiLatency        = "hermes_api_latency_seconds" // TODO(evanSpendlove) add this constant to metrics.go
	minFileIDToDelete = 11                           // we can delete files starting from the file Hermes_11
	maxFileIDToDelete = 50                           // there are 40 files to delete from [Hermes_11,Hermes_50]
)

// DeleteFile deletes the file, corresponding to the ID passed, in the target storage system bucket.
// It then checks that the file has been deleted by trying to get the object.
// Arguments:
//	- ctx: context allows this probe can be cancelled if needed.
// - fileID: ID of the file to be deleted. Must be within inclusive range: 11-50.
//	- config: HermesProbeDef config for the probe calling this function.
//	- target: target run information stored in struct from probe/probe.go
//	- client: initialised storage client for this target system.
//	- logger: logger associated with the probe calling this function.
// Returns:
//	- fileID: returns the ID of the deleted file.
//	- err:
//		Status:
//		- StateJournalInconsistent: the file to be deleted does not exist in Hermes' StateJournal.
//		- FileMissing: the file to be deleted could not be found in the target bucket.
//		- BucketMissing: the target bucket on this target system was not found.
//		- ProbeFailed: there was an error during one of the API calls and the probe failed.
func DeleteFile(ctx context.Context, fileID int32, target *probe.Target, client stiface.Client, logger *logger.Logger) (int32, error) {
	bucket := target.Target.GetBucketName()

	if fileID < minFileIDToDelete || fileID > maxFileIDToDelete {
		return fileID, fmt.Errorf("DeleteFile(%q, %q) failed; status %v: expected fileID %d to be within valid inclusive range: 11-50", bucket, fileID, metrics.InvalidArgument, fileID)
	}

	// TODO(evanSpendlove): Add custom error object to return value and modify all returns.
	filename, ok := target.Journal.Filenames[fileID]
	if !ok {
		return fileID, fmt.Errorf("StateJournalInconsistent: Journal.Filenames has no entry with file ID = %d", fileID)
	}

	target.Journal.Intent = &pb.Intent{
		FileOperation: pb.Intent_CREATE,
		Filename:      filename,
	}

	file := client.Bucket(bucket).Object(filename)

	// TODO(#77): Refactor timing into using function from metrics.go
	start := time.Now()
	if err := file.Delete(ctx); err != nil {
		var status metrics.ExitStatus
		switch err {
		case storage.ErrObjectNotExist:
			status = metrics.FileMissing
		case storage.ErrBucketNotExist:
			status = metrics.BucketMissing
		default:
			status = metrics.ProbeFailed
		}

		target.LatencyMetrics.APICallLatency[metrics.APIDeleteFile][status].Metric(apiLatency).AddFloat64(time.Now().Sub(start).Seconds())
		return fileID, fmt.Errorf("DeleteFile(%q, %q) failed; status %v: %w", bucket, filename, status, err)
	}
	target.LatencyMetrics.APICallLatency[metrics.APIDeleteFile][metrics.Success].Metric(apiLatency).AddFloat64(time.Now().Sub(start).Seconds())

	query := &storage.Query{Prefix: filename}
	// TODO(#77): Refactor timing into using function from metrics.go
	start = time.Now()
	objects := client.Bucket(bucket).Objects(ctx, query)
	for {
		obj, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if obj.Name == filename {
			status := metrics.ProbeFailed
			target.LatencyMetrics.APICallLatency[metrics.APIListFiles][status].Metric(apiLatency).AddFloat64(time.Now().Sub(start).Seconds())
			return fileID, fmt.Errorf("DeleteFile(%q, %q) failed; status %v: object %v still listed after delete", bucket, filename, status, obj.Name)
		}
		if err != nil {
			status := metrics.BucketMissing
			target.LatencyMetrics.APICallLatency[metrics.APIListFiles][status].Metric(apiLatency).AddFloat64(time.Now().Sub(start).Seconds())
			return fileID, fmt.Errorf("DeleteFile(%q, %q) failed; status %v: %w", bucket, filename, status, err)
		}
	}
	target.LatencyMetrics.APICallLatency[metrics.APIListFiles][metrics.Success].Metric(apiLatency).AddFloat64(time.Now().Sub(start).Seconds())

	// Update in-memory NIL file after delete operation.
	delete(target.Journal.Filenames, fileID)

	logger.Infof("Object %v deleted in bucket %s.", file, bucket)
	return fileID, nil
}

// PickFileToDelete picks which file to delete and returns the integer ID of this file.
// Returns:
//	- ID: returns the ID of the file to be deleted.
func PickFileToDelete() int32 {
	rand.Seed(time.Now().UnixNano())
	return int32(rand.Intn(maxFileIDToDelete-minFileIDToDelete+1) + minFileIDToDelete)
}
