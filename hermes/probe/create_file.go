package probe

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/cloudprober/logger"
	"github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes/options"
	"github.com/google/cloudprober/targets/endpoint"
	"github.com/googleinterns/step224-2020/blob/evan/hermesMain/cmd"
)

type RandomFile struct {
	Seed int64
	Size int // File size in bytes
}

func (f *RandomFile) NewReader() *randomFileReader {
	return &randomFileReader{size: f.Size, rand: rand.New(rand.NewSource(f.Seed))}
}

type randomFileReader struct {
	size int
	// Current reading index.
	i    int
	rand *rand.Rand
}

func (r *randomFileReader) Read(b []byte) (n int, err error) {
	if r.i >= r.size {
		return 0, io.EOF
	}
	if len(b) > r.size-r.i {
		n, err = r.rand.Read(b[:r.size-r.i])
	} else {
		n, err = r.rand.Read(b)
	}
	r.i += n
	return
}

func (f *RandomFile) generateChecksum() (string, error) {
	r := f.NewReader()
	h := sha1.New()
	buff := make([]byte, 8)
	for {
		n, err := r.Read(buff)
		if err == io.EOF {
			break
		}
		if _, err = io.Copy(h, r); err != nil {
			return nil, fmt.Errorf("io.Copy: %v", err)
		}

	}

	// returns checksum in hex notation
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (f *RandomFile) generateFileName(fileID int) (string, error) {
	checksum, err := f.generateChecksum()
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("Hermes_%02d_%v", fileID, checksum), nil
}

func (p *MonitorProbe) CreateFile(ctx context.Context, target string, bucketName string, fileID int, fileSize int) error {
	if fileID < 1 || fileID > 50 {
		return fmt.Errorf("At %v the file ID provided wasn't in the required range [0,50]", time.Now())
	}
	f := RandomFile{Seed: fileID, Size: fileSize}
	fileName, err := f.generateFileName(fileID)
	if err != nil {
		return err
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(ctx, p.opts.Timeout)
	defer cancel()
	r := f.NewReader()

	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	buff := make([]byte, 8)
	for {
		n, err := r.Read(buff)
		if err == io.EOF {
			break
		}
		if _, err = io.Copy(wc, r); err != nil {
			return fmt.Errorf("io.Copy: %v", err)
		}

	}
	//io
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil

}