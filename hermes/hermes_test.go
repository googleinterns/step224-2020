// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// hermes_test.go tests that the Hermes struct and its associated functions work as expected.
// This includes initialising Cloudprober.

package hermes

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/cloudprober"
	"github.com/google/cloudprober/web"
)

const cfg = "grpc_port:9314"

// Test InitialiseCloudproberFromConfig tests that Hermes can initialise
// Cloudprober from a config and that it starts running properly.
func TestInitialiseCloudproberFromConfig(t *testing.T) {
	hermes := &Hermes{}

	if err := hermes.InitialiseCloudproberFromConfig(cfg); err != nil {
		t.Fatalf("Expected no error from hermes.InitialiseCloudproberFromConfig(), got %v", err)
	}

	hermes.Ctx, hermes.CancelCloudprober = context.WithCancel(context.Background())
	defer hermes.CancelCloudprober()

	// Sets up web UI for cloudprober.
	web.Init()

	// Start running Cloudprober instance from Hermes context
	cloudprober.Start(hermes.Ctx)

	setConfig := cloudprober.GetConfig() // This gets the current config that Cloudprober is running with.
	if cfg != setConfig {
		t.Error("Config mismatch")
	}
}

// TestFileOperationString tests that the String() method of FileOperation properly formats the enum as a string.
func TestFileOperationString(t *testing.T) {
	var fileOpTests = []struct {
		want  string
		input FileOperation
	}{
		{"Create", Create},
		{"Delete", Delete},
	}

	for _, op := range fileOpTests {
		got := fmt.Sprintf("%s", op.input)
		if got != op.want {
			t.Errorf("Want %s, got %s", op.want, got)
		}
	}
}
