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

const (
	fileSizeBytes = 100
)

func TestReadFile(t *testing.T) {
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
		t.Fatalf("metrics.NewMetrics(): %v", err)
	}

	ctx := context.Background()
	client := fakegcs.NewClient()
	bucketName := "test_bucket_probe0"
	fakeBucketHandle := client.Bucket(bucketName)
	if err := fakeBucketHandle.Create(ctx, bucketName, nil); err != nil {
		t.Error(err)
	}
	logger := fakegcs.NewLogger(ctx).Logger

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
		if err = create.CreateFile(ctx, target, tc.fileIDCreate, fileSizeBytes, client, logger); err != nil {
			t.Fatalf("CreateFile(fileID: %d) set up failed %v", tc.fileIDCreate, err)
		}
		err = ReadFile(ctx, target, tc.fileIDRead, fileSizeBytes, client, logger)
		if err != nil && !tc.wantErr {
			t.Errorf("ReadFile(fileID: %d) = %w, expected: nil", tc.fileIDRead, err)
		}
		if tc.wantErr && err == nil {
			t.Errorf("ReadFile(fileID: %d) = %w, expected an error", tc.fileIDRead, err)
		}

	}
}
