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
// Author: Alicja Kwiecinska, GitHub: alicjakwie
//
// Package read contains all of the logic necessary to verify the availability and consistency of the file contents and names in GCS.
//
// TODO(#76) change the type of fileID to int
// TODO(#79) unify  total space alocated Mib or MiB
package read

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/cloudprober/logger"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
	"google.golang.org/api/iterator"
)

const (
	// universal format of the names of files in the storage system Hermes_ID_checksum
	FileNameFormat          = "Hermes_%02d_%x"
	minFileID               = 1
	maxFileID               = 50
	maxFileSizeBytes        = 1000
	hermesAPILatencySeconds = "hermes_api_latency_seconds"
)

// ReadFile creates and stores a file with randomized contents in the target storage system.
// Before it creates and stores a file it logs an intent to do so in the target storage system's journal.
// It verifies that the creation and storage process was successful.
// Finally, it updates the filenames map in the target's journal and record the exit status in the logger.
// Arguments:
//          ctx: it carries deadlines and cancellation signals that might orinate from the main probe located in probe.go.
//          target: contains information about target storage system, carries an intent log in the form of a StateJournal and it used to export metrics.
//          fileID: the unique identifer of every file, it cannot be repeated. It needs to be in the range [minFileID, maxFileID]. FileID 0 is reserved for a special file called the NIL file.
//          client: is a storage client. It is used as an interface to interact with the target storage system.
//          logger: a cloudprober logger used to record the exit status of the ReadFile operation in a target bucket. The logger passed MUST be a valid logger.
// Returns:
//          error: an error string with detailed information about the status and fileID. Nil is returned when the operation is successful.
func ReadFile(ctx context.Context, target *probe.Target, fileID int32, fileSize int, client stiface.Client, logger *logger.Logger) error {
	if fileID < minFileID || fileID > maxFileID {
		return fmt.Errorf("invalid argument: fileID = %d; want %d <= fileID <= %d", fileID, minFileID, maxFileID)
	}
	bucketName := target.Target.GetBucketName()

	// Verify that the file is present in the State Journal
	fileName, ok := target.Journal.Filenames[fileID]
	if !ok {
		return fmt.Errorf("file with the ID: %d is missing from the State Journal", fileID)
	}

	// Verify that the file that we want to read is in fact present in the target system
	fileNamePrefix := fmt.Sprintf(FileNameFormat, fileID, "")
	query := &storage.Query{Prefix: fileNamePrefix}
	start := time.Now()
	objIter := client.Bucket(bucketName).Objects(ctx, query)
	var namesFound []string
	for {
		obj, err := objIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("existence check  for failed due to: %w", err)
		}
		namesFound = append(namesFound, obj.Name)
	}
	finish := time.Now()
	if len(namesFound) == 0 {
		target.LatencyMetrics.APICallLatency[metrics.APIListFiles][metrics.FileMissing].Metric(hermesAPILatencySeconds).AddFloat64(finish.Sub(start).Seconds())
		return fmt.Errorf("could not read file as the file with the provided ID %d does not exist in bucket %q", fileID, bucketName)
	}
	if len(namesFound) != 1 {
		return fmt.Errorf("check failed for ID %d expected exactly one file in bucket %q with prefix %q; found %d: %v", fileID, bucketName, fileNamePrefix, len(namesFound), namesFound)
	}
	if namesFound[0] != fileName {
		return fmt.Errorf("check failed  for ID %d expected file name present %q got %q", fileID, fileName, namesFound[0])
	}
	target.LatencyMetrics.APICallLatency[metrics.APIListFiles][metrics.Success].Metric(hermesAPILatencySeconds).AddFloat64(finish.Sub(start).Seconds())

	start = time.Now()
	reader, err := client.Bucket(bucketName).Object(fileName).NewReader(ctx)
	if err != nil {
		var status metrics.ExitStatus
		switch err {
		case storage.ErrObjectNotExist:
			status = metrics.FileMissing
		case storage.ErrBucketNotExist:
			status = metrics.BucketMissing
		default:
			status = metrics.ProbeFailed
		}
		target.LatencyMetrics.APICallLatency[metrics.APIGetFile][status].Metric(hermesAPILatencySeconds).AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errorf(".%q: could not read file %q: %w", status, fileName, err)
	}
	defer reader.Close()
	h := sha1.New()
	if _, err := io.Copy(h, reader); err != nil {
		return fmt.Errorf("checksum calculation failed io.Copy: %w", err)
	}
	gotChecksum := fmt.Sprintf("%x", h.Sum(nil))
	wantChecksum := fileName[len(fileNamePrefix):]
	if gotChecksum != wantChecksum {
		return fmt.Errorf("the calculated checksum: %q does not match the checksum in the file name: %q", gotChecksum, wantChecksum)
	}
	logger.Infof("verified consistency for object %q in bucket %q", fileName, bucketName)
	return nil
}
