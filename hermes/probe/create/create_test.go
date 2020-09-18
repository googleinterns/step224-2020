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
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/fakegcs"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"

	metricpb "github.com/google/cloudprober/metrics/proto"
	probepb "github.com/googleinterns/step224-2020/config/proto"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

var fileNamePrefixLength = len(fmt.Sprintf(FileNameFormat, 0, ""))

func TestNewRandomFile(t *testing.T) {
	tests := []struct {
		fileID   int32
		fileSize int
		want     *randomFile
		wantErr  bool
	}{
		{51, 12, nil, true},
		{0, 50, nil, true},
		{3, 100, &randomFile{3, 100}, false},
		{12, 100, &randomFile{12, 100}, false},
		{3, 0, nil, true},
		{3, 1001, nil, true},
	}
	for _, tc := range tests {
		got, err := newRandomFile(tc.fileID, tc.fileSize)
		if tc.want == nil && got != nil {
			t.Errorf("{%d, %d}.newRandomFile = {%d, %d} expected nil", tc.fileID, tc.fileSize, got.id, got.sizeBytes)
		}
		if got == nil && tc.want != nil {
			t.Errorf("{%d, %d}.newRandomFile = nil expected {%d, %d}", tc.fileID, tc.fileSize, tc.want.id, tc.want.sizeBytes)
		}
		if err != nil && !tc.wantErr {
			t.Errorf("{%d, %d}.newRandomFile() failed and returned an unexpected error %w", tc.fileID, tc.fileSize, err)
		}
		if err == nil && tc.wantErr {
			t.Errorf("{%d, %d}.newRandomFile() failed expected an error got nil", tc.fileID, tc.fileSize)
		}
		if got != nil && (got.id != tc.want.id || got.sizeBytes != tc.want.sizeBytes) {
			t.Errorf("{%d, %d}.newRandomFile() =  {%d,%d}, expected {%d,%d}", tc.fileID, tc.fileSize, got.id, got.sizeBytes, tc.want.id, tc.want.sizeBytes)
		}

	}

}

func TestFileName(t *testing.T) {
	tests := []struct {
		file *randomFile
		want string
	}{
		{&randomFile{3, 100}, "Hermes_03_"},
		{&randomFile{12, 100}, "Hermes_12_"},
		{&randomFile{8, 20}, "Hermes_08_"},
	}

	for _, tc := range tests {
		got, err := tc.file.fileName()
		if err != nil {
			t.Errorf("{%d, %d}.fileName() failed and returned an unexpected error %w", tc.file.id, tc.file.sizeBytes, err)
		}
		if !strings.HasPrefix(got, tc.want) {
			t.Errorf("{%d, %d}.fileName() =  %qchecksum expected %qchecksum", tc.file.id, tc.file.sizeBytes, got[:fileNamePrefixLength], tc.want)
		}
	}
}

func TestChecksum(t *testing.T) {
	file := randomFile{11, 100}
	checksum, err := file.checksum()
	if err != nil {
		t.Error(err)
	}
	otherFile := randomFile{13, 100}
	otherChecksum, err := otherFile.checksum()
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%x", checksum) == fmt.Sprintf("%x", otherChecksum) {
		t.Errorf("{%d, %d}.checksum = {%d,%d}.checksum expected {%d, %d}.checksum != {%d, %d}.checksum ", file.id, file.sizeBytes, otherFile.id, otherFile.sizeBytes, file.id, file.sizeBytes, otherFile.id, otherFile.sizeBytes)

	}
	file = randomFile{11, 100}
	checksumAgain, err := file.checksum()
	if err != nil {
		t.Error(err)
	}
	if fmt.Sprintf("%x", checksum) != fmt.Sprintf("%x", checksumAgain) {
		t.Errorf("{%d, %d}.checksum != {%d, %d}.checksum expected {%d, %d}.checksum = {%d, %d}.checksum", file.id, file.sizeBytes, file.id, file.sizeBytes, file.id, file.sizeBytes, file.id, file.sizeBytes)
	}
}

func TestCreateFile(t *testing.T) {
	ctx := context.Background()
	client := fakegcs.NewClient()
	bucketName := "test_bucket_probe0"
	fakeBucketHandle := client.Bucket(bucketName)
	if err := fakeBucketHandle.Create(ctx, bucketName, nil); err != nil {
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
