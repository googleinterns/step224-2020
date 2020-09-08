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
// Main program loop for Hermes. This initialises Cloudprober so that Hermes can
// interact with it through gRPCs.

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/google/cloudprober"
	"github.com/google/cloudprober/web"
)

var (
	rpcPort = flag.Int("rpc_port", 9314, "The port that the gRPC server of Cloudprober will run on.")
)

func main() {
	flag.Parse()

	if err := cloudprober.InitFromConfig(buildConfig()); err != nil {
		glog.Exitf("cloudprober could not be initialised from config: grpc_port: %d, err:%v", *rpcPort, err)
	}

	// Sets up web UI for cloudprober.
	web.Init()

	cloudprober.Start(context.Background())

	// Wait forever
	select {}
}

// buildConfig() builds the configuration details for Cloudprober based on the flag contents.
// Returns:
// - string: Returns the configuration details for Cloudprober as a string.
func buildConfig() string {
	return fmt.Sprintf("grpc_port: %d", *rpcPort)
}
