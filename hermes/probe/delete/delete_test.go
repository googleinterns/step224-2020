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
// Author: Evan Spendlove, GitHub: evanSpendlove.
//
// Delete_file_test tests the file deletion operation for MonitorProbe.

package delete

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/logger"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/googleinterns/step224-2020/hermes/probe"
	"github.com/googleinterns/step224-2020/hermes/probe/fakegcs"
	"google.golang.org/api/iterator"

	metricpb "github.com/google/cloudprober/metrics/proto"
	monitorpb "github.com/googleinterns/step224-2020/config/proto"
	m "github.com/googleinterns/step224-2020/hermes/probe/metrics"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

// genTestConfig generates a test HermesProbeDef proto config.
// Arguments:
//	- name: pass the name for this probe instance.
// Returns:
//	- hermesExtension: returns the HermesProbeDef extension.
func genTestConfig(name string) *monitorpb.HermesProbeDef {
	hermesExtension := &monitorpb.HermesProbeDef{
		ProbeName: proto.String(name),
		Targets: []*monitorpb.Target{
			genTargetPb(),
		},
		TargetSystem: monitorpb.HermesProbeDef_GCS.Enum(),
		IntervalSec:  proto.Int32(3600),
		TimeoutSec:   proto.Int32(60),
		ProbeLatencyDistribution: &metricpb.Dist{
			Buckets: &metricpb.Dist_ExplicitBuckets{
				ExplicitBuckets: "0.1,0.2,0.4,0.6,0.8,1.6,3.2,6.4,12.8,1000",
			},
		},
		ApiCallLatencyDistribution: &metricpb.Dist{
			Buckets: &metricpb.Dist_ExplicitBuckets{
				ExplicitBuckets: "0.000000002,0.000000004,0.000000008,0.000000016,0.00000032,0.000000064,0.000000128,100",
			},
		},
	}
	return hermesExtension
}

// getTargetPb generates a target proto config.
// Returns:
//	- target: target proto config
func genTargetPb() *monitorpb.Target {
	return &monitorpb.Target{
		Name:                   "hermes",
		TargetSystem:           monitorpb.Target_GOOGLE_CLOUD_STORAGE,
		TotalSpaceAllocatedMib: int64(100),
		BucketName:             bucketName,
	}
}

const (
	// TODO(#76): Change these to int from int32.
	firstID    = int32(1)
	lastID     = int32(50)
	contents   = "abc123"
	hash       = "6367c48dd193d56ea7b0baad25b19455e529f5ee"
	bucketName = "test_bucket_5"
)

// genTestTarget generates an initialised test Target struct.
func genTestTarget(cfg *monitorpb.HermesProbeDef, t *testing.T) *probe.Target {
	filenames := make(map[int32]string)

	for i := firstID; i <= lastID; i++ {
		var id string
		if i <= 9 {
			id = fmt.Sprintf("0%d", i)
		} else {
			id = fmt.Sprintf("%d", i)
		}
		filenames[i] = fmt.Sprintf("Hermes_%s_%s", id, hash)
	}

	metrics, err := m.NewMetrics(cfg, genTargetPb())
	if err != nil {
		t.Fatalf("could not initialise metrics using config and target provided: %v", err)
	}

	return &probe.Target{
		Target: cfg.GetTargets()[0],
		Journal: &journalpb.StateJournal{
			Intent:    &journalpb.Intent{},
			Filenames: filenames,
		},
		LatencyMetrics: metrics,
	}
}

// createTestFiles creates a test bucket and the required test files.
func createTestFiles(ctx context.Context, client stiface.Client, t *testing.T) {
	bucket := client.Bucket(bucketName)
	if err := bucket.Create(ctx, "", nil); err != nil {
		t.Fatalf("failed to create fake bucket, err: %v", err)
	}

	for i := firstID; i <= lastID; i++ {
		filename := fmt.Sprintf("Hermes_%02d_%s", i, hash)

		writer := client.Bucket(bucketName).Object(filename).NewWriter(ctx)
		n, err := writer.Write([]byte(contents))
		if err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
		if n != len([]byte(contents)) {
			t.Fatalf("short write: wrote %d bytes; wanted %d", n, len([]byte(contents)))
		}
		writer.Close()
	}
}

func TestDeleteRandomFile(t *testing.T) {
	testProbeName := "testDelete1"
	ctx := context.Background()

	client := fakegcs.NewClient()
	createTestFiles(ctx, client, t)

	target := genTestTarget(genTestConfig(testProbeName), t)

	logger, err := logger.NewCloudproberLog(testProbeName)
	if err != nil {
		t.Fatalf("failed to initialise logger: %v", err)
	}

	fileID, err := DeleteFile(ctx, PickFileToDelete(), target, client, logger)
	if err != nil {
		t.Errorf("deleteRandomFile(ID: %d) failed: expected error as %v, got %v", fileID, nil, err)
	}
	filename := target.Journal.Filenames[fileID]
	objects := client.Bucket(bucketName).Objects(ctx, &storage.Query{Prefix: filename})
	for {
		obj, err := objects.Next()
		if err == iterator.Done {
			break
		}
		if obj.Name == filename {
			t.Errorf("deleteRandomFile failed: expected object to be deleted, got object found.")
		}
		if err != nil {
			t.Fatalf("deleteRandomFile failed: %v", err)
		}
	}
}

// TODO(evanSpendlove): Add more tests that check that DeleteFile() throws the correct errors.
