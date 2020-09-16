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
	"google.golang.org/api/iterator"

	probe "github.com/googleinterns/step224-2020/hermes/probe"
	m "github.com/googleinterns/step224-2020/hermes/probe/metrics"
	pb "github.com/googleinterns/step224-2020/hermes/proto"
)

// DeleteRandomFile deletes a random file in the target storage system bucket.
// It then checks that the file has been deleted by trying to get the object.
// Arguments:
//	- ctx: pass the context so this probe can be cancelled if needed.
//	- target: pass the target run information.
//	- client: pass an initialised storage client for this target system.
//	- logger: pass the logger associated with the probe calling this function.
//		-> logger can log to Google Cloud if configured to do so.
// Returns:
//	- fileID: returns the ID of the file delete OR a missing file to be created if one is found.
//	- err:
//		Status:
//		- StateJournalInconsistent: the file to be deleted does not exist in Hermes' StateJournal.
//		- FileMissing: the file to be deleted could not be found in the target bucket.
//		- BucketMissing: the target bucket on this target system was not found.
//		- ProbeFailed: there was an error during one of the API calls and the probe failed.
func DeleteRandomFile(ctx context.Context, target *probe.Target, client stiface.Client, logger *logger.Logger) (int32, error) {
	fileID := pickFileToDelete()
	return DeleteFile(ctx, fileID, target, client, logger)
}

// DeleteFile deletes the file, corresponding to the ID passed, in the target storage system bucket.
// It then checks that the file has been deleted by trying to get the object.
// Arguments:
//	- ctx: pass the context so this probe can be cancelled if needed.
//	- config: pass the HermesProbeDef config for the probe calling this function.
//	- target: pass the target run information.
//	- client: pass an initialised storage client for this target system.
//	- logger: pass the logger associated with the probe calling this function.
// Returns:
//	- fileID: returns the ID of the file delete OR a missing file to be created if one is found.
//	- err:
//		Status:
//		- StateJournalInconsistent: the file to be deleted does not exist in Hermes' StateJournal.
//		- FileMissing: the file to be deleted could not be found in the target bucket.
//		- BucketMissing: the target bucket on this target system was not found.
//		- ProbeFailed: there was an error during one of the API calls and the probe failed.
func DeleteFile(ctx context.Context, fileID int32, target *probe.Target, client stiface.Client, logger *logger.Logger) (int32, error) {
	bucket := target.Target.GetBucketName()

	// TODO(evanSpendlove): Add custom error object to return value and modify all returns.
	filename, ok := target.Journal.Filenames[fileID]
	if !ok {
		return fileID, fmt.Errorf("delete(%q, %q) failed; status StateJournalInconsistent: expected fileID %d to exist", bucket, filename, fileID)
	}

	target.Journal.Intent = &pb.Intent{
		FileOperation: pb.Intent_CREATE,
		Filename:      filename,
	}

	file := client.Bucket(bucket).Object(filename)

	start := time.Now()
	if err := file.Delete(ctx); err != nil {
		var status m.ExitStatus
		switch err {
		case storage.ErrObjectNotExist:
			status = m.FileMissing
		case storage.ErrBucketNotExist:
			status = m.BucketMissing
		default:
			status = m.ProbeFailed
		}

		target.LatencyMetrics.APICallLatency[m.APIDeleteFile][status].Metric("latency_s").AddFloat64(time.Now().Sub(start).Seconds())
		return fileID, fmt.Errorf("delete(%q, %q) failed; status %v: %w", bucket, filename, status, err)
	}
	target.LatencyMetrics.APICallLatency[m.APIDeleteFile][m.Success].Metric("latency_s").AddFloat64(time.Now().Sub(start).Seconds())

	query := &storage.Query{Prefix: filename}
	start = time.Now()
	objects := client.Bucket(bucket).Objects(ctx, query)
	target.LatencyMetrics.APICallLatency[m.APIListFiles][m.Success].Metric("latency_s").AddFloat64(time.Now().Sub(start).Seconds())
	for {
		obj, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if obj.Name == filename {
			status := m.ProbeFailed
			return fileID, fmt.Errorf("delete(%q, %q) failed; status %v: object %v still listed after delete", bucket, filename, status, obj.Name)
		}
		if err != nil {
			status := m.BucketMissing
			return fileID, fmt.Errorf("delete(%q, %q) failed; status %v: %w", bucket, filename, status, err)
		}
	}

	// Update in-memory NIL file after delete operation.
	delete(target.Journal.Filenames, int32(fileID))

	logger.Infof("Object %v deleted in bucket %s.", file, bucket)
	return fileID, nil
}

// pickFileToDelete picks which file to delete and returns the integer ID of this file.
// Returns:
//	- ID: returns the ID of the file to be deleted.
func pickFileToDelete() int32 {
	const (
		beg                    = 10 // we can delete files staring from the file Hermes_10
		numberOfDeletableFiles = 41 // there are 41 files to delete from [Hermes_10,Hermes_50]
	)
	rand.Seed(time.Now().UnixNano())
	return int32(rand.Intn(numberOfDeletableFiles) + beg)
}
