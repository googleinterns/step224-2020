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
// Parameters:
// - config: config should be the contents of a Cloudprober config file. This is most likely: "grpc_port=9314"
//           -> the "grpc_port:" field is the only required field for the config.
// Returns:
// - context: 
// - cancel func():
// - error: 
func InitialiseCloudproberFromConfig(config string) (context.Context, func(), error) {

	err := cloudprober.InitFromConfig(config)
	if err != nil {
		glog.Errorf("failed to initialise cloudprober, err: %v", err)
		return nil, nil, err
	}

	// web.Init sets up web UI for cloudprober.
	web.Init()

	ctx, cancelCloudprober := context.WithCancel(context.Background()) // Create new context with a cancel() function

	cloudprober.Start(ctx) // Start running Cloudprober instance

	return ctx, cancelCloudprober, nil
}
