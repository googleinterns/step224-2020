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
// Package probe implements the probe that Hermes uses to monitor
// a storage system.

package probe

import (
	"context"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"

	monitorpb "github.com/googleinterns/step224-2020/config/proto"
)

func TestDeleteRandomFile(t *testing.T) {
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		t.Fatalf("client conn failed: could not connect to storage system with client: %v", err)
	}

	client := stiface.AdaptClient(c) // Don't use this, use fakeClient.

	target := &monitorpb.Target{
		Name:                   "hermes",
		TargetSystem:           monitorpb.Target_GOOGLE_CLOUD_STORAGE,
		TotalSpaceAllocatedMib: int64(100),
		BucketName:             "test_bucket_5",
	}

	mp := &MonitorProbe{}

	fileID, err := deleteRandomFile(ctx, mp, target, &client)
	if err != nil {
		t.Errorf("deleteRandomFile(ID: %d) failed: expected error as %v, got %v", fileID, nil, err)
	}

	// w := client.Bucket("my-bucket").Object("my-object").NewWriter(ctx)
	// w.ObjectAttrs().ContentType = "text/plain"
}
