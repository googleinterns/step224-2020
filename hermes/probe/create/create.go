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
// package create contains all of the logic necessary to create files with randomized contents in GCS.
//
// TODO(#76) change the type of fileID to int
package create

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
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
	FileNameFormat          = "Hermes_%02d_%x"
	minFileID               = 1
	maxFileID               = 50
	maxFileSizeBytes        = 1000
	hermesAPILatencySeconds = "hermes_api_latency_s"
)

type randomFile struct {
	id        int32
	sizeBytes int
}

type randomFileReader struct {
	sizeBytes int
	// currently reading this byte
	i    int
	rand *rand.Rand
}

// Read implements the standard Read interface, it returns a random stream of bytes generated by a pseudo-random generator.
// Argument:
//      buf: a byte slice serves as an output buffer
// Returns:
//		n:  the number of bytes read
//		err: error it should be nil when not done reading and io.EOF once the whole file contents have been read
func (r *randomFileReader) Read(buf []byte) (n int, err error) {
	// check whether Read is done
	if r.i >= r.sizeBytes {
		return 0, io.EOF
	}
	b := buf
	if len(buf) > r.sizeBytes-r.i {
		// if the length of buffer is greater than the number of bytes left to read make the sizeBytes of the buffer match the number of bytes left to read
		b = buf[:r.sizeBytes-r.i]
	}
	// n is now the length  of the buffer
	n, err = r.rand.Read(b)
	if err != nil {
		// in this case n = 0
		return n, err
	}
	r.i += n
	return n, err
}

func (f *randomFile) newReader() *randomFileReader {
	//  id will serve as a Seed and i - index of the currently read byte  will be set to 0 automatically in the returned reader
	return &randomFileReader{
		sizeBytes: f.sizeBytes,
		rand:      rand.New(rand.NewSource(int64(f.id))),
	}
}

func newRandomFile(id int32, sizeBytes int) (*randomFile, error) {
	if id < minFileID || id > maxFileID {
		return nil, fmt.Errorf("invalid argument: id = %d; want %d <= id <= %d", id, minFileID, maxFileID)
	}
	if sizeBytes > maxFileSizeBytes || sizeBytes <= 0 {
		return nil, fmt.Errorf("invalid argument: sizeBytes = %d; want 0 < sizeBytes <= %d", sizeBytes, maxFileSizeBytes)
	}
	return &randomFile{id: id, sizeBytes: sizeBytes}, nil
}

func (f *randomFile) checksum() ([]byte, error) {
	r := f.newReader()
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, fmt.Errorf("io.Copy: %w", err)
	}
	return h.Sum(nil), nil
}

func (f *randomFile) fileName() (string, error) {
	checksum, err := f.checksum()
	if err != nil {
		return "", fmt.Errorf("{%d, %d}.checksum = nil,  %w", f.id, f.sizeBytes, err)
	}
	return fmt.Sprintf(FileNameFormat, f.id, checksum), nil
}

// CreateFile creates and stores a file with randomized contents in the target storage system.
// Before it creates and stores a file it logs an intent to do so in the target storage system's journal.
// It verifies that the creation and storage process was successful.
// Finally, it updates the filenames map in the target's journal and record the exit status in the logger.
// Arguments:
//          ctx: it carries deadlines and cancellation signals that might orinate from the main probe located in probe.go.
//          target: contains information about target storage system, carries an intent log in the form of a StateJournal and it used to export metrics.
//          fileID: the unique identifer of every randomFile, it cannot be repeated. It needs to be in the range [minFileID, maxFileID]. FileID 0 is reserved for a special file called the NIL file.
//          client: is a storage client. It is used as an interface to interact with the target storage system.
//          logger: a cloudprober logger used to record the exit status of the CreateFile operation in a target bucket. The logger passed MUST be a valid logger.
// Returns:
//          error: an error string with detailed information about the status and fileID. Nil is returned when the operation is successful.
func CreateFile(ctx context.Context, target *probe.Target, fileID int32, fileSize int, client stiface.Client, logger *logger.Logger) error {
	f, err := newRandomFile(fileID, fileSize)
	if err != nil {
		return err
	}
	fileName, err := f.fileName()
	if err != nil {
		return err
	}
	target.Journal.Intent = &pb.Intent{FileOperation: pb.Intent_CREATE, Filename: fileName}
	var status metrics.ExitStatus
	if _, ok := target.Journal.Filenames[fileID]; ok {
		status = metrics.UnknownFileFound
		return fmt.Errorf("CreateFile(fileID: %d).%q could not create file as file with this ID already exists", fileID, status)
	}
	r := f.newReader()
	start := time.Now()
	bucketName := target.Target.GetBucketName()
	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		switch err {
		case storage.ErrBucketNotExist:
			status = metrics.BucketMissing
		default:
			status = metrics.ProbeFailed
		}
		target.LatencyMetrics.APICallLatency[metrics.APICreateFile][status].Metric(hermesAPILatencySeconds).AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errorf("CreateFile(id: %d).%q: could not create file %q: %w", fileID, status, fileName, err)
	}
	if err := wc.Close(); err != nil {
		status = metrics.WriterCloseFailed
		target.LatencyMetrics.APICallLatency[metrics.APICreateFile][status].Metric(hermesAPILatencySeconds).AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errorf("Writer.Close: %w with status %q", err, status)
	}
	status = metrics.Success
	target.LatencyMetrics.APICallLatency[metrics.APICreateFile][status].Metric(hermesAPILatencySeconds).AddFloat64(time.Now().Sub(start).Seconds())

	// Verify that the file that has just been created is in fact present in the target system
	fileNamePrefix := fmt.Sprintf(FileNameFormat, fileID, "")
	query := &storage.Query{Prefix: fileNamePrefix}
	start = time.Now()
	objIter := client.Bucket(bucketName).Objects(ctx, query)
	var namesFound []string
	for {
		obj, err := objIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("CreateFile check failed due to: %w", err)
		}
		namesFound = append(namesFound, obj.Name)
	}
	finish := time.Now()
	if len(namesFound) == 0 {
		target.LatencyMetrics.APICallLatency[metrics.APIListFiles][metrics.FileMissing].Metric(hermesAPILatencySeconds).AddFloat64(finish.Sub(start).Seconds())
		return fmt.Errorf("CreateFile check failed no files with prefix %q found", fileNamePrefix)
	}
	if len(namesFound) != 1 {
		return fmt.Errorf("expected exactly one file in bucket %q with prefix %q; found %d: %v", bucketName, fileNamePrefix, len(namesFound), namesFound)
	}
	if namesFound[0] != fileName {
		return fmt.Errorf("CreateFile check failed expected file name present %q got %q", fileName, namesFound[0])
	}
	target.LatencyMetrics.APICallLatency[metrics.APIListFiles][metrics.Success].Metric(hermesAPILatencySeconds).AddFloat64(finish.Sub(start).Seconds())

	target.Journal.Filenames[fileID] = fileName
	logger.Infof("Object %q added in bucket %q.", fileName, bucketName)
	return nil
}