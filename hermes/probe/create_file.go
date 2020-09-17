// TODO (#70) add license header and author line
// TODO (#68) add doc strings
// TODO (#72) change error types and error messages

package probe

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

	m "github.com/googleinterns/step224-2020/hermes/probe/metrics"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

const (
	begin            = 1    // ID of the first HermesFile
	end              = 50   // ID of the last HermesFile
	maxFileSizeBytes = 1000 // maximum allowed file size in bytes
)

type RandomFile struct {
	ID   int64 // ID is a positive integer in the range [1,50]
	Size int   // File size in bytes
}

type randomFileReader struct {
	size int // size in bytes
	i    int // currently reading this byte
	rand *rand.Rand
}

func (r *randomFileReader) readDone() bool {
	return r.i >= r.size
}

func (r *randomFileReader) bytesLeft() int {
	return (r.size - r.i)
}

func (r *randomFileReader) bufferTooLong(buf []byte) bool {
	return len(buf) > r.bytesLeft()
}

func (r *randomFileReader) Read(buf []byte) (n int, err error) {
	if r.readDone() {
		return 0, io.EOF
	}
	b := buf
	if r.bufferTooLong(b) {
		b = buf[:r.bytesLeft()]
	}
	n, err = r.rand.Read(b) // n is the length  of the buffer
	if err != nil {
		return n, err // in this case n=0
	}
	r.i += n
	return n, err
}

// Warning: NewReader is not thread safe
func (f *RandomFile) NewReader() *randomFileReader {
	// ID will serve as a Seed and i - index of the currently read byte  will be set to 0 automatically in the returned reader
	return &randomFileReader{
		size: f.Size,
		rand: rand.New(rand.NewSource(f.ID)),
	}
}

func NewRandomFile(fileID int64, fileSize int) (*RandomFile, error) {
	if fileID < begin || fileID > end {
		return &RandomFile{}, fmt.Errorf("the file ID provided %v wasn't in the required range [%v, %v]", fileID, begin, end)
	}
	if fileSize > maxFileSizeBytes || fileSize <= 0 {
		return &RandomFile{}, fmt.Errorf("invalid argument: RandomFile.Size = %d; want 0 < size <= %d", fileSize, maxFileSizeBytes)
	}
	return &RandomFile{ID: int64(fileID), Size: fileSize}, nil
}

func (f *RandomFile) CheckSum() ([]byte, error) {
	r := f.NewReader()
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, fmt.Errorf("io.Copy: %v", err)
	}
	return h.Sum(nil), nil
}

func (f *RandomFile) FileName() (string, error) {
	checksum, err := f.CheckSum()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Hermes_%02d_%x", f.ID, checksum), nil
}

func CreateFile(ctx context.Context, target *probe.Target, fileID int32, fileSize int, client stiface.Client, logger *logger.Logger) error {
	f, err := NewRandomFile(int64(fileID), fileSize)
	if err != nil {
		return err
	}
	fileName, err := f.FileName()
	if err != nil {
		return err
	}
	target.Journal.Intent = &journalpb.Intent{FileOperation: journalpb.Intent_CREATE, Filename: fileName}
	if _, ok := target.Journal.Filenames[fileID]; ok {
		var status m.ExitStatus
		status = m.UnknownFileFound
		return fmt.Errorf("CreateFile(ID: %d).%v could not create file as file with this ID already exists", fileID, status)
	}
	r := f.NewReader()
	start := time.Now()
	bucketName := target.Target.GetBucketName()
	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		var status m.ExitStatus
		switch err {
		case storage.ErrBucketNotExist:
			status = m.BucketMissing
		default:
			status = m.ProbeFailed
		}
		target.LatencyMetrics.APICallLatency[m.APICreateFile][status].Metric("hermes_api_latency_s").AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errorf("CreateFile(id: %d).%v: could not create file %s: %w", fileID, status, fileName, err)
	}
	if err := wc.Close(); err != nil {
		target.LatencyMetrics.APICallLatency[m.APICreateFile][m.WriterCloseFailed].Metric("hermes_api_latency_s").AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errorf("Writer.Close: %v", err)
	}
	target.LatencyMetrics.APICallLatency[m.APICreateFile][m.Success].Metric("hermes_api_latency_s").AddFloat64(time.Now().Sub(start).Seconds())
	prefix := fileName[0:9]
	query := &storage.Query{Prefix: prefix}
	start = time.Now()
	objIter := client.Bucket(bucketName).Objects(ctx, query)
	end := time.Now()
	obj, err := objIter.Next()
	if err != nil {
		target.LatencyMetrics.APICallLatency[m.APIListFiles][m.FileMissing].Metric("hermes_api_latency_s").AddFloat64(end.Sub(start).Seconds())
		return fmt.Errorf("CreateFile check failed: %w", err)
	}
	if obj.Name != fileName {
		fmt.Errorf("CreateFile check failed expected file name present %v got %v", fileName, obj.Name)
	}
	target.LatencyMetrics.APICallLatency[m.APIListFiles][m.Success].Metric("hermes_api_latency_s").AddFloat64(end.Sub(start).Seconds())

	target.Journal.Filenames[fileID] = fileName
	if logger != nil {
		logger.Infof("Object %v added in bucket %s.", fileName, bucketName)
	}
	return nil
}
