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
// Author: Alicja Kwiecinska GitHub: alicjakwie
//
// TODO (#72) change error types to be compatible with ProbeError

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
		fileID   int32
		fileSize int
		want   *randomFile
		wantErr  bool
	}{
		{51, 12, &randomFile{}, true},
		{0, 50, &randomFile{}, true},
		{3, 100, &randomFile{3, 100}, false},
		{12, 100, &randomFile{12, 100}, false},
		{3, 0, &randomFile{}, false},
		{3, 1001, &randomFile{}, false},
	}
	for _, test := range tests {
		got, err := newRandomFile(test.fileID, test.fileSize)
		if err != nil && !test.wantErr {
			t.Errorf("{%d, %d}.newRandomFile() failed and returned an unexpected error %v", test.fileID, test.fileSize, err)
		}
		if err == nil && test.wantErr {
			t.Errorf("{%d, %d}.newRandomFile() failed expected an error got nil", test.fileID, test.fileSize)
		}
		if got.id != test.wantRF.id || got.size != test.wantRF.size {
			t.Errorf("{%d, %d}.newRandomFile() failed expected {%d,%d}, got {%d,%d}", test.fileID, test.fileSize, test.wantRF.id, test.wantRF.size, got.id, got.size)
		}

	}

}

func TestFileName(t *testing.T) {
	tests := []struct {
		file     *randomFile
		wantName string
	}{
		{&randomFile{3, 100}, "Hermes_03"},
		{&randomFile{12, 100}, "Hermes_12"},
		{&randomFile{8, 20}, "Hermes_08"},
	}

	for _, test := range tests {
		got, err := test.file.fileName()
		if err != nil {
			t.Errorf("{%d, %d}.fileName() failed and returned an unexpected error %v", test.file.id, test.file.size, err)
		}
		if got[:FileNamePrefixSize] != test.wantName {
			t.Errorf("{%d, %d}.fileName() failed expected prefix %q got prefix %q", test.file.id, test.file.size, test.wantName, got[:FileNamePrefixSize])
		}
	}
}

func TestChecksum(t *testing.T) {
	file := randomFile{11, 100}
	checksum, err := file.checksum()
	if err != nil {
		t.Error(err)
	}
	file = randomFile{13, 100}
	otherChecksum, err := file.checksum()
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%x", checksum) == fmt.Sprintf("%x", otherChecksum) {
		t.Errorf("Checksum returned the same value for two different randomFiles {%d, %d} and {%d, %d}", file.id, file.size, file.id, file.size)

	}
	file = randomFile{11, 100}
	checksumAgain, err := file.checksum()
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%x", checksum) != fmt.Sprintf("%x", checksumAgain) {
		t.Errorf("Checksum returned different values for the same randomFiles {%d, %d}", file.id, file.size)
	}
}

func TestCreateFile(t *testing.T) {
	ctx := context.Background()
	bucketName := "test_bucket_probe0"
	client := fakegcs.NewClient()
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
		t.Fatalf("metrics.NewMetrics(): %v", err)
	}
	logger := fakegcs.NewLogger(ctx).Logger
	if err := CreateFile(ctx, target, fileID, fileSize, client, logger); err != nil {
		t.Error(err)
	}
}
