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
	m "github.com/googleinterns/step224-2020/hermes/probe/metrics"
)

const (
	begin            = 1    // ID of the first HermesFile
	end              = 50   // ID of the last HermesFile
	maxFileSizeBytes = 1000 // maximum allowed file size in bytes
)

type RandomHermesFile struct {
	ID   int64 // ID is a positive integer in the range [1,50]
	Size int   // File size in bytes
}

type randomHermesFileReader struct {
	size int // size in bytes
	i    int // currently reading this byte
	rand *rand.Rand
}

func (r *randomHermesFileReader) readDone() bool {
	return r.i >= r.size
}

func (r *randomHermesFileReader) bytesLeft() int {
	return (r.size - r.i)
}

func (r *randomHermesFileReader) bufferTooLong(buffer []byte) bool {
	return len(buffer) > r.bytesLeft()
}

func (r *randomHermesFileReader) Read(buf []byte) (n int, err error) {
	if r.readDone() {
		return 0, io.EOF
	}
	b := buf
	if r.bufferTooLong(b) {
		b = buf[:r.bytesLeft()]
	}
	n, err = r.rand.Read(b) //n is the length  of the buffer
	if err != nil {
		return n, err // in this case n=0 and we return 0 as nil can't be returned as a type int argument
	}
	r.i += n
	return n, err
}

// Warning: NewReader is not thread safe
func (f *RandomHermesFile) NewReader() *randomHermesFileReader {
	return &randomHermesFileReader{size: f.Size, rand: rand.New(rand.NewSource(f.ID))} // ID will serve as a Seed and i - index of the currently read byte  will be set to 0 automatically in the returned reader
}

func (f *RandomHermesFile) CheckSum() ([]byte, error) {
	r := f.NewReader()
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, fmt.Errorf("io.Copy: %v", err)
	}
	// returns checksum in hex notation
	return h.Sum(nil), nil
}

func (f *RandomHermesFile) FileName() (string, error) {
	if f.ID < begin || f.ID > end {
		return "", fmt.Errorf("The file ID provided %v wasn't in the required range [1,50]", f.ID)
	}
	if f.Size > maxFileSizeBytes {
		return "", fmt.Errorf("The file size provided %v bytes exceeded the limit %v bytes", f.Size, maxFileSizeBytes)
	}
	if f.Size <= 0 {
		return "", fmt.Errorf("The file size provided %v is not a positive number as required", f.Size)
	}
	checksum, err := f.CheckSum()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Hermes_%02d_%v", f.ID, fmt.Sprintf("%x", checksum)), nil
}

func CreateFile(ctx context.Context, target *probe.Target, fileID int32, fileSize int, client stiface.Client, logger *logger.Logger) error {
	f := RandomHermesFile{ID: int64(fileID), Size: fileSize}
	bucketName := target.Target.GetBucketName()
	fileName, err := f.FileName()
	if err != nil {
		return err
	}
	if _, ok := target.Journal.Filenames[fileID]; ok {
		var status m.ExitStatus
		status = m.UnknownFileFound
		return fmt.Errorf("CreateFile(ID: %d).%v could not create file as file with this ID already exists", fileID, status)
	}
	r := f.NewReader()
	start := time.Now()
	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		var status m.ExitStatus
		switch err {
		case storage.ErrBucketNotExist:
			status = m.BucketMissing
		default:
			status = m.ProbeFailed
		}
		target.LatencyMetrics.APICallLatency[m.APIDeleteFile][status].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errorf("CreateFile(id: %d).%v: could not create file %s: %w", fileID, status, fileName, err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	target.Journal.Filenames[fileID] = fileName
	if logger != nil {
		logger.Infof("Object %v added in bucket %s.", fileName, bucketName)
	}
	return nil

}