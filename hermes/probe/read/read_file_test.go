package read

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/create"
	"github.com/googleinterns/step224-2020/hermes/probe/fakegcs"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"

	metricpb "github.com/google/cloudprober/metrics/proto"
	probepb "github.com/googleinterns/step224-2020/config/proto"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

func TestReadFile(t *testing.T) {
	ctx := context.Background()
	bucketName := "test_bucket_probe0"
	client := fakegcs.NewClient()
	fbh := client.Bucket(bucketName)                         // fakeBucketHandle
	if err := fbh.Create(ctx, bucketName, nil); err != nil { // creates the bucket with name "test_bucket_probe0"
		t.Error(err)
	}
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
	fileID := int32(6)
	fileSize := 50
	if err := create.CreateFile(ctx, target, fileID, fileSize, client, nil); err != nil {
		t.Error(err)
	}
	if err := ReadFile(ctx, target, fileID, client, nil); err != nil {
		t.Error(err)
	}
}
