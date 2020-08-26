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
// cloudprober.go initialises and runs a cloudprober instance using a custom configuration file.

package cloudprober

import (
	"context"

	"github.com/golang/glog"
	"github.com/google/cloudprober"
	"github.com/google/cloudprober/web"
)

// InitialiseCloudproberFromConfig initialises Cloudprober from the config passed as an argument.
// It then sets up the web UI for Cloudprober, and starts running Cloudprober.
// The config passed as a string should be the contents of a Cloudprober config file, not a file path or filename.
// For most instances, this is likely just "grpc_port=9314" as this is the only required value in the config.
func InitialiseCloudproberFromConfig(config string) {

	err := cloudprober.InitFromConfig(config)
	if err != nil {
		glog.Exitf("Error initializing cloudprober. Err: %v", err)
	}

	// web.Init sets up web UI for cloudprober.
	web.Init()

	cloudprober.Start(context.Background()) // Start running Cloudprober instance
}
