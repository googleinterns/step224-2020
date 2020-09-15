package create

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	metricpb "github.com/google/cloudprober/metrics/proto"
	probepb "github.com/googleinterns/step224-2020/config/proto"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
	"github.com/googleinterns/step224-2020/hermes/probe/mock"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

func TestFileName(t *testing.T) {
	tests := []struct {
		file     RandomHermesFile
		wantName string
		wantErr  string
	}{
		{RandomHermesFile{51, 12}, "", "The file ID provided 51 wasn't in the required range [1,50]"},
		{RandomHermesFile{0, 50}, "", "The file ID provided 0 wasn't in the required range [1,50]"},
		{RandomHermesFile{3, 100}, "Hermes_03", ""},
		{RandomHermesFile{12, 100}, "Hermes_12", ""},
		{RandomHermesFile{3, 0}, "", "The file size provided 0 is not a positive number as required"},
		{RandomHermesFile{3, 1001}, "", "The file size provided 1001 bytes exceeded the limit 1000 bytes"},
	}

	for _, test := range tests {
		got, err := test.file.FileName()
		if err == nil {
			if test.wantErr != "" {
				t.Errorf("{%v, %v}.FileName() gave wrong error expected %s, got nil", test.file.ID, test.file.Size, test.wantErr)
			}
			if test.wantName != got[0:9] {
				t.Errorf("{%v, %v}.FileName() failed expected prefix %s, got %s", test.file.ID, test.file.Size, test.wantName, got[0:9])
			}
		} else {
			if got != "" {
				t.Errorf("{%v, %v}.FileName() failed expected empty string, got %s", test.file.ID, test.file.Size, got)
			}
			if err.Error() != test.wantErr {
				t.Errorf("{%v, %v}.FileName() gave wrong error expected %s, got %s", test.file.ID, test.file.Size, test.wantErr, err.Error())
			}
		}
	}
}

func TestChecksum(t *testing.T) {
	f := RandomHermesFile{11, 100}
	fTwo := RandomHermesFile{13, 1000}
	fCopy := RandomHermesFile{11, 100}
	checksum, err := f.CheckSum()
	if err != nil {
		t.Error(err)
	}
	checksumTwo, err := fTwo.CheckSum()
	if err != nil {
		t.Error(err)
	}

	checksumCopy, err := fCopy.CheckSum()
	if err != nil {
		t.Error(err)
	}

	if fmt.Sprintf("%x", checksum) != fmt.Sprintf("%x", checksumCopy) { //comparing checksums converted to strings as a slice can only be compared to nil
		t.Errorf("Checksum returned different values for the same RandomHermesFiles {%v, %v}", f.ID, f.Size)
	}

	if fmt.Sprintf("%x", checksum) == fmt.Sprintf("%x", checksumTwo) {
		t.Errorf("Checksum returned the same value for two different RandomHermesFiles {%v, %v} and {%v, %v}", f.ID, f.Size, fTwo.ID, fTwo.Size)
	}

}

func TestCreateFile(t *testing.T) {
	ctx := context.Background()
	bucketName := "test_bucket_probe0"
	client := mock.NewFakeClient()
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

	err := fmt.Errorf("") //err cannot be defined in the if statement below using ':=' as then the set up fails with the error:  expected identifier on left side of :=
	if target.LatencyMetrics, err = metrics.NewMetrics(hp, probeTarget); err != nil {
		t.Error(err)
	}

	if err := CreateFile(ctx, target, fileID, fileSize, client, nil); err != nil {
		t.Error(err)
	}
}
