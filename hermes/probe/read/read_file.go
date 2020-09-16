// TODO (#70) add license header and author line
// TODO (#69) add doc strings

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
	m "github.com/googleinterns/step224-2020/hermes/probe/metrics"
)

func checkSumFromName(fileName string) (string, error) {
	if len(fileName) < 11 {
		return "", fmt.Errorf("The filename stored in the journal is shorter than 11 characters, hence does not match Hermes_id_checksum")
	}
	return fileName[10:len(fileName)], nil
}

func ReadFile(ctx context.Context, target *probe.Target, fileID int32, client stiface.Client, logger *logger.Logger) error {
	bucketName := target.Target.GetBucketName()
	fileName, ok := target.Journal.Filenames[fileID]
	if !ok {
		var status m.ExitStatus
		status = m.FileMissing
		return fmt.Errorf("ReadFile(ID: %d).%v could not read file as the file with the provided ID does not exist", fileID, status)
	}
	prefix := fileName[0:9]
	query := &storage.Query{Prefix: prefix}
	start := time.Now()
	objIter := client.Bucket(bucketName).Objects(ctx, query)
	end := time.Now()
	obj, err := objIter.Next()
	if err != nil {
		target.LatencyMetrics.APICallLatency[m.APIListFiles][m.FileMissing].Metric("hermes_api_latency_s").AddFloat64(end.Sub(start).Seconds())
		return fmt.Errorf("CreateFile check failed: %w", err)
	}
	if obj.Name != fileName {
		fmt.Errorf("The name %v of the file with id %v in the storage system doesn't match the name %v stored in the Journal", obj.Name, fileID, fileID)
	}
	target.LatencyMetrics.APICallLatency[m.APIListFiles][m.Success].Metric("hermes_api_latency_s").AddFloat64(end.Sub(start).Seconds())
	wantChecksum, err := checkSumFromName(fileName)
	if err != nil {
		return err
	}
	start = time.Now()
	reader, err := client.Bucket(bucketName).Object(fileName).NewReader(ctx)
	if err != nil {
		var status m.ExitStatus
		switch err {
		case storage.ErrObjectNotExist:
			status = m.FileMissing
		case storage.ErrBucketNotExist:
			status = m.BucketMissing
		default:
			status = m.ProbeFailed
		}
		target.LatencyMetrics.APICallLatency[m.APIGetFile][status].Metric("hermes_api_latency_s").AddFloat64(time.Now().Sub(start).Seconds())
		return fmt.Errorf("CreateFile(id: %d).%v: could not create file %s: %w", fileID, status, fileName, err)
	}
	h := sha1.New()
	if _, err := io.Copy(h, reader); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	gotChecksum := fmt.Sprintf("%x", h.Sum(nil))
	if gotChecksum != wantChecksum {
		return fmt.Errorf("The calculated checksum does not match the checksum stored in the file name")
	}
	if err := reader.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}
