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

func TestNewRandomFile(t *testing.T) {
	tests := []struct {
		fileID   int64
		fileSize int
		wantErr  bool
		wantRF   *RandomFile
	}{
		{51, 12, true, &RandomFile{}},
		{0, 50, true, &RandomFile{}},
		{3, 100, false, &RandomFile{3, 100}},
		{12, 100, false, &RandomFile{12, 100}},
		{3, 0, true, &RandomFile{}},
		{3, 1001, true, &RandomFile{}},
	}
	for _, test := range tests {
		got, err := NewRandomFile(test.fileID, test.fileSize)
		if err != nil && !test.wantErr {
			t.Errorf("{%v, %v}.NewRandomFile() failed and returned an unexpected error %s", test.fileID, test.fileSize, err.Error())
		}
		if err == nil && test.wantErr {
			t.Errorf("{%v, %v}.NewRandomFile() failed expected an error got nil", test.fileID, test.fileSize)
		}
		if got.ID != test.wantRF.ID || got.Size != test.wantRF.Size {
			t.Errorf("{%v, %v}.NewRandomFile() failed expected {%v,%v}, got {%v,%v}", test.fileID, test.fileSize, test.wantRF.ID, test.wantRF.Size, got.ID, got.Size)
		}

	}

}

func TestFileName(t *testing.T) {
	tests := []struct {
		file     *RandomFile
		wantName string
	}{
		{&RandomFile{3, 100}, "Hermes_03"},
		{&RandomFile{12, 100}, "Hermes_12"},
		{&RandomFile{9, 20}, "Hermes_09"},
	}

	for _, test := range tests {
		got, err := test.file.FileName()
		if err != nil {
			t.Errorf("{%v, %v}.NewRandomFile() failed and returned an unexpected error %s", test.file.ID, test.file.Size, err.Error())
		}
		if got[0:9] != test.wantName {
			t.Errorf("{%v, %v}.NewRandomFile() failed expected prefix %s got prefix %s", test.file.ID, test.file.Size, test.wantName, got[0:9])
		}
	}
}

func TestChecksum(t *testing.T) {
	file := RandomFile{11, 100}
	checksum, err := file.CheckSum()
	if err != nil {
		t.Error(err)
	}
	file = RandomFile{13, 100}
	otherChecksum, err := file.CheckSum()
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%x", checksum) == fmt.Sprintf("%x", otherChecksum) {
		t.Errorf("Checksum returned the same value for two different RandomFiles {%v, %v} and {%v, %v}", file.ID, file.Size, file.ID, file.Size)

	}
	file = RandomFile{11, 100}
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
		t.Fatalf("Metric set up failed due to %s", err.Error())
	}

	if err := CreateFile(ctx, target, fileID, fileSize, client, nil); err != nil {
		t.Error(err)
	}
}
