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
// Client.go creates the client for connecting to, and making requests to, the
// gRPC server of Cloudprober.
//
// TODOs:
// TODO(#38) Add support for secure connection to RPC server with credentials, if supported by Cloudprober.

package cmd

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"

	"github.com/golang/glog"
	proberPb "github.com/google/cloudprober/prober/proto"
	"github.com/google/cloudprober/probes"
	configPb "github.com/google/cloudprober/probes/proto"
	"google.golang.org/grpc"
)

var (
	client proberPb.CloudproberClient // Cloudprober gRPC client retained between RPC calls
)

const (
	rpcServer string = "localhost:9314" // If the grpc_port in the Cloudprober config is changed, need to update this too.
)

// getClient() ensures there is an active client connection to the Cloudprober gRPC server.
// If there is not an active connection, it will make one and assign a Cloudprober client
// to the global client variable.Using a global client variable allows Hermes to
// retain the connection between RPCs.
func getClient() {
	if client == nil { // If there is not an active client connection, make one.
		conn, err := grpc.Dial(rpcServer, grpc.WithInsecure()) // Make a connection
		if err != nil {
			glog.Fatal(err)
		}
		client = proberPb.NewCloudproberClient(conn) // Create a new Cloudprober gRPC Client
	}
}

// addProbeFromConfig() adds a probe to Cloudprober via the gRPC client.
// Prerequisites - probePb:
// - The probe config passed must extend the Cloudprober ProbeDef found at /probes/proto in Cloudprober.
// - The probe config must also be unmarshalled before being passed as an argument.
// Prerequisite: The probe type must be registered as an extension.
func addProbeFromConfig(probePb *configPb.ProbeDef) {
	_, err := client.AddProbe(context.Background(), &proberPb.AddProbeRequest{ProbeConfig: probePb}) // Adds the probe to Cloudprober

	if err != nil {
		glog.Error("Probe could not be registered: ", probePb)
		glog.Exit(err)
	}
}

// RegisterAndAddProbe registers a probe type and adds the probe to Cloudprober to be run.
// Prequisites - probePb:
// - The probe config passed must extend the Cloudprober ProbeDef found at /probes/proto in Cloudprober.
// - The probe config must also be unmarshalled before being passed as an argument.
func RegisterAndAddProbe(extensionNumber int, probePb *configPb.ProbeDef, probe probes.Probe) {
	getClient() // Ensures there is an active client connection to Cloudprober gRPC server.
	probes.RegisterProbeType(extensionNumber, func() probes.Probe { return probe }) // First, register the probe as an extension with Cloudprober.

	fmt.Println("Probe: ", *probePb)
	fmt.Println("Probe extension: ", proto.RegisteredExtensions(probePb)[200])

	// Add the probe to Cloudprober
	// The probe will be scheduled and run by Cloudprober
	// Adding a probe will consume the probe type registration.
	// If you want to add multiple probes, you must register and add each one individually.
	// Only one extension type can be registered at any given time.
	// If more are registered, Cloudprober will throw an error.
	addProbeFromConfig(probePb)
}

// RemoveProbe removes a probe from Cloudprober, given the probe name (located in the probe config).
func RemoveProbe(probeName string) {
	getClient() // Ensures there is an active client connection to Cloudprober gRPC server.

	var probeToRemove *string // Need to use a string pointer for RemoveProbeRequest{}
	probeToRemove = &probeName

	_, err := client.RemoveProbe(context.Background(), &proberPb.RemoveProbeRequest{ProbeName: probeToRemove}) // Remove probe from Cloudprober
	if err != nil {
		glog.Exit(err)
	}
}

// ListProbes() returns a list of active probes from Cloudprober.
func ListProbes() *proberPb.ListProbesResponse {
	getClient() // Ensures there is an active client connection to Cloudprober gRPC server.

	response, err := client.ListProbes(context.Background(), &proberPb.ListProbesRequest{}) // Submit ListProbes() rpc.
	if err != nil {
		glog.Errorf("error while list probes: %v", err)
	}

	return response
}
