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
// Author: Evan Spendlove (@evanSpendlove)
//
// cloudprober_test.go tests that Cloudprober can be initialised from a config
// and that it starts running properly.

package cloudprober

import (
	"testing"

	"github.com/google/cloudprober"
)

// Test InitialiseCloudproberFromConfig tests that Hermes can initialise
// Cloudprober from a config and that it starts running properly.
func TestInitialiseCloudproberFromConfig(t *testing.T) {
	cfg := "grpc_port:9314"
	_, cancel, err := InitialiseCloudproberFromConfig(cfg) // This initialises Cloudprober and starts it running.
	if err != nil {
		t.Errorf("Expected no error for InitialiseCloudproberFromConfig(), got %v", err)
	}

	setConfig := cloudprober.GetConfig() // This gets the current config that Cloudprober is running with.
	if cfg != setConfig {
		t.Error("Config mismatch")
	}

	cancel() // Close down Cloudprober instance
}
