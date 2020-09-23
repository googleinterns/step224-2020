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
// TODO(#76) change the type of fileID to int
// TODO(#79) unify  total space alocated Mib or MiB

package read

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/logger"
	"github.com/googleinterns/step224-2020/hermes/probe/create"
	"github.com/googleinterns/step224-2020/hermes/probe/fakegcs"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
	"github.com/googleinterns/step224-2020/hermes/probe/target"

	metricpb "github.com/google/cloudprober/metrics/proto"
	probepb "github.com/googleinterns/step224-2020/config/proto"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

const (
	fileSizeBytes     = 100
	readTestProbeName = "read_test_probe"
)

func TestReadFile(t *testing.T) {
	target := &target.Target{
		&probepb.Target{
			Name:                   "hermes",
			TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
			TotalSpaceAllocatedMib: 1,
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
		TotalSpaceAllocatedMib: 1,
		BucketName:             "test_bucket_probe0",
	}

	var err error
	if target.LatencyMetrics, err = metrics.NewMetrics(hp, probeTarget); err != nil {
		t.Fatalf("metrics.NewMetrics(): %v", err)
	}

	ctx := context.Background()
	client := fakegcs.NewClient()
	bucket := "test_bucket_probe0"
	bucketHandle := client.Bucket(bucket)
	if err := bucketHandle.Create(ctx, bucket, nil); err != nil {
		t.Fatalf("error creating bucket %q: %v", bucket, err)
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
		{2, 51, true},
		{8, 7, true},
		{10, 10, false},
		{4, 4, false},
		{6, 0, true},
	}
	for _, tc := range tests {
		if err := create.CreateFile(ctx, target, tc.fileIDCreate, fileSizeBytes, client, logger); err != nil {
			t.Fatalf("CreateFile(fileID: %d) set up failed %v", tc.fileIDCreate, err)
		}
		if err := ReadFile(ctx, target, tc.fileIDRead, fileSizeBytes, client, logger); (err != nil) != tc.wantErr {
			t.Errorf("ReadFile(fileID: %d) = %w, want: %v", tc.fileIDRead, err, tc.wantErr)
		}
	}
}
