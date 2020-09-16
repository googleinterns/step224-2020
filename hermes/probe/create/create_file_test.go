// TODO (#70) add license header and author line
// TODO (#68) add doc strings

package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/fakegcs"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"

	metricpb "github.com/google/cloudprober/metrics/proto"
	probepb "github.com/googleinterns/step224-2020/config/proto"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

func TestFileName(t *testing.T) {
	tests := []struct {
		file     RandomFile
		wantName string
		wantErr  bool
	}{
		{RandomFile{51, 12}, "", true},
		{RandomFile{0, 50}, "", true},
		{RandomFile{3, 100}, "Hermes_03", false},
		{RandomFile{12, 100}, "Hermes_12", false},
		{RandomFile{3, 0}, "", true},
		{RandomFile{3, 1001}, "", true},
	}

	for _, test := range tests {
		got, err := test.file.FileName()
		if err == nil {
			if test.wantErr {
				t.Errorf("{%v, %v}.FileName() failed expected an error got nil", test.file.ID, test.file.Size)
			}
			if test.wantName != got[0:9] {
				t.Errorf("{%v, %v}.FileName() failed expected prefix %s, got %s", test.file.ID, test.file.Size, test.wantName, got[0:9])
			}
		} else {
			if got != "" {
				t.Errorf("{%v, %v}.FileName() failed expected empty string, got %s", test.file.ID, test.file.Size, got)
			}
			if !test.wantErr {
				t.Errorf("{%v, %v}.FileName() failed and gave unexpected error %s", test.file.ID, test.file.Size, err.Error())
			}
		}
	}
}

func TestChecksum(t *testing.T) {
	file := RandomFile{11, 100}
	otherFile := RandomFile{13, 1000}
	checksum, err := file.CheckSum()
	if err != nil {
		t.Error(err)
	}
	otherChecksum, err := otherFile.CheckSum()
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%x", checksum) == fmt.Sprintf("%x", otherChecksum) {
		t.Errorf("Checksum returned the same value for two different RandomFiles {%v, %v} and {%v, %v}", file.ID, file.Size, otherFile.ID, otherFile.Size)

	}
	checksumAgain, err := file.CheckSum()
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%x", checksum) != fmt.Sprintf("%x", checksumAgain) {
		t.Errorf("Checksum returned different values for the same RandomFiles {%v, %v}", file.ID, file.Size)
	}
}

func TestCreateFile(t *testing.T) {
	ctx := context.Background()
	bucketName := "test_bucket_probe0"
	client := fakegcs.NewClient()
	fbh := client.Bucket(bucketName)                         // fakeBucketHandle
	if err := fbh.Create(ctx, bucketName, nil); err != nil { // creates the bucket with name "test_bucket_probe0"
		t.Error(err)
	}
	fileID := int32(6)
	fileSize := 50
	target := &probe.Target{
		&probepb.Target{
			Name:                   "hermes",
			TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
			TotalSpaceAllocatedMib: int64(1000),
			BucketName:             "test_bucket_probe0",
		},
		&journalpb.StateJournal{
			Filenames: make(map[int32]string),
		},
		&metrics.Metrics{},
	}
	hp := &probepb.HermesProbeDef{
		ProbeName: proto.String("createfile_test"),
		Targets: []*probepb.Target{
			&probepb.Target{
				Name:                   "hermes",
				TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
				TotalSpaceAllocatedMib: int64(100),
				BucketName:             "test_bucket_probe0",
			},
		},
		TargetSystem: probepb.HermesProbeDef_GCS.Enum(),
		IntervalSec:  proto.Int32(3600),
		TimeoutSec:   proto.Int32(60),
		ProbeLatencyDistribution: &metricpb.Dist{
			Buckets: &metricpb.Dist_ExplicitBuckets{
				ExplicitBuckets: "0.1,0.2,0.4,0.6,0.8,1.6,3.2,6.4,12.8,1",
			},
		},
		ApiCallLatencyDistribution: &metricpb.Dist{
			Buckets: &metricpb.Dist_ExplicitBuckets{
				ExplicitBuckets: "0.1,0.2,0.4,0.6,0.8,1.6,3.2,6.4,12.8,1",
			},
		},
	}
	probeTarget := &probepb.Target{
		Name:                   "hermes",
		TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
		TotalSpaceAllocatedMib: int64(100),
		BucketName:             "test_bucket_probe0",
	}

	var err error
	if target.LatencyMetrics, err = metrics.NewMetrics(hp, probeTarget); err != nil {
		t.Error(err)
	}

	if err := CreateFile(ctx, target, fileID, fileSize, client, nil); err != nil {
		t.Error(err)
	}
}
