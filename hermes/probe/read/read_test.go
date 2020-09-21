package read

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/logger"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/create"
	"github.com/googleinterns/step224-2020/hermes/probe/fakegcs"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"

	metricpb "github.com/google/cloudprober/metrics/proto"
	probepb "github.com/googleinterns/step224-2020/config/proto"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

const (
	fileSizeBytes     = 100
	readTestProbeName = "read_test_probe"
)

func TestReadFile(t *testing.T) {
	target := &probe.Target{
		&probepb.Target{
			Name:                   "hermes",
			TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
			TotalSpaceAllocatedMib: int64(1),
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
				TotalSpaceAllocatedMib: int64(1),
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
		t.Fatalf("metrics.NewMetrics(): %v", err)
	}

	ctx := context.Background()
	client := fakegcs.NewClient()
	bucketName := "test_bucket_probe0"
	bucketHandle := client.Bucket(bucketName)
	if err := bucketHandle.Create(ctx, bucketName, nil); err != nil {
		t.Fatalf("error creating bucket %q: %v", bucketName, err)
	}
	logger, err := logger.NewCloudproberLog(readTestProbeName)
	if err != nil {
		t.Fatalf("failed to initialise logger: %v", err)
	}
	tests := []struct {
		fileIDCreate int32
		fileIDRead   int32
		wantErr      bool
	}{
		{3, 3, false},
		{4, 51, true},
		{10, 12, true},
		{12, 10, false},
		{7, 4, false},
		{6, 0, true},
	}
	for _, tc := range tests {
		if err := create.CreateFile(ctx, target, tc.fileIDCreate, fileSizeBytes, client, logger); err != nil {
			t.Fatalf("CreateFile(fileID: %d) set up failed %v", tc.fileIDCreate, err)
		}
		if err := ReadFile(ctx, target, tc.fileIDRead, fileSizeBytes, client, logger); (err != nil) != tc.wantErr {
			t.Errorf("ReadFile(fileID: %d) = %w, expected: %v", tc.fileIDRead, err, tc.wantErr)
		}
	}
}
