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
// Author: Evan Spendlove (@evanSpendlove)
//
// Main program loop for Hermes. This initialises Cloudprober so that Hermes can
// interact with it through gRPCs.

package main

import (
	"flag"
	"strconv"

	cp "github.com/googleinterns/step224-2020/cloudprober"
)

var (
	rpc_port = flag.Int("rpc_port", 9314, "The port that the gRPC server of Cloudprober will run on")
)

func main() {
	flag.Parse()

	err := cp.InitialiseCloudproberFromConfig("grpc_port: " + strconv.Itoa(*rpc_port))

	if err != nil {
		glog.Exitf("cloudprober could not be initialised from config: grpc_port: %d, err:%v", *rpc_port, err)
	}

	// Wait forever
	select {}
}
