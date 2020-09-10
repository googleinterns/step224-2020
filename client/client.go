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
// Author: Evan Spendlove, GitHub: @evanSpendlove.
//
// Client.go creates the client for connecting to, and making requests to, the
// gRPC server of Cloudprober.
//
// TODOs:
// TODO(#38): Add support for secure connection to RPC server with credentials, if supported by Cloudprober.
// TODO(evanSpendlove): Include details of these errors (addProbeFromConfig method) in the error message.

// Client contains the gRPC client implementation used to send RPCs to Cloudprober.
package client

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	proberpb "github.com/google/cloudprober/prober/proto"
	cpprobes "github.com/google/cloudprober/probes"
	configpb "github.com/google/cloudprober/probes/proto"
)

type CloudproberClient struct {
	conn      *grpc.ClientConn           // Client connection for opening and closing the connection
	clientMux sync.Mutex                 // Mutex lock for making CloudproberClient safe for concurrent use
	client    proberpb.CloudproberClient // Cloudprober gRPC client retained between RPC calls
}

// NewClient creates, initialises and returns a new CloudproberClient.
// This is the correct way to create a new instance of CloudproberClient.
// Parameters:
// - rpcServer: Pass the rpcServer address and port as a string formatted as: "hostname:port", e.g. "localhost:1234"
// Returns:
// - *CloudproberClient: returns a CloudproberClient with a client connection to the Cloudprober gRPC server.
// - error:
//	  - nil: No error occurred and the client was successfully initialised.
//	  - gRPC Error: See https://godoc.org/google.golang.org/grpc/codes for error codes.
//		-> See https://github.com/grpc/grpc-go/blob/d25c71b54334380ff1febd25d88064b36de44b3c/clientconn.go#L123
func NewClient(rpcServer string) (*CloudproberClient, error) {
	conn, err := grpc.Dial(rpcServer, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("GRPC_DIAL_ERROR: %v", err)
	}
	return &CloudproberClient{conn: conn, client: proberpb.NewCloudproberClient(conn)}, nil
}

// CloseConn is used for closing the client connection with the gRPC server.
func (c *CloudproberClient) CloseConn() {
	c.clientMux.Lock()
	defer c.clientMux.Unlock()

	if err := c.conn.Close(); err != nil {
		glog.Errorf("Cloudprober gRPC client could not close connection, err: %v", err)
	}
}

// RegisterAndAddProbe registers a probe type and adds the probe to Cloudprober to be run.
// Parameters:
// - ctx: Context used for cancelling RPCs.
// - extensionNumber: The proto extension number from the proto of this probe type.
// - probePb: This probe config must be unmarshalled before being passed as an argument.
// - hermesProbeToAdd - An empty probe object of the desired probe type.
// Returns:
// TODO(evanSpendlove): Include details of these errors in the error message.
// - error:
//         - Code 3, InvalidArgument: probe config cannot be nil
//         - error parsing regexp [...]: the machine name for this probe to run on does not compile as a regexp.
//         - Code 6, AlreadyExists: probe [name] is already defined - the probe has already been added
//         - Code 2, Unknown:
//		- only one probe extension is allowed per probe, got %d extensions
//			-> More than one probe extension type has been registered using cpprobes.RegisterProbeType()
//		- no probe extension in probe config
//			-> The probe proto does not have an extension registered.
//		- no probes registered for the extension: [extension code]
//			-> The proto has a probe extension, but there is no probe extension registered matching that extension code.
//		- proto.GetExtension() error: extension missing, invalid, incomplete or proto cannot be extended.
//			-> The proto.SetExtension() method must be called and supplied with valid data.
//		- unregistered user defined probe: [probe name]
//			-> The probe type has not been registered as a user defined probe.
//		- unknown probe type: [probe type]
//			-> The probe type does not match any expected types (including user-defined and extension).
//		- probe.Init() error: an error occurred when the probe, created from the config passsed, was initialised.
//			-> There was a problem in the Init() function of the probe type passed.
//		- options.BuildProbeOptions() error: an error occurred when building the options for this probe from the config supplied.
//			-> The options could not be built from the config supplied.
func (c *CloudproberClient) RegisterAndAddProbe(ctx context.Context, extensionNumber int, probePb *configpb.ProbeDef, hermesProbeToAdd cpprobes.Probe) error {
	c.clientMux.Lock()
	defer c.clientMux.Unlock()

	cpprobes.RegisterProbeType(extensionNumber, func() cpprobes.Probe { return hermesProbeToAdd })

	// Only one extension type per proto config can be registered at any given time.
	// If more are registered, Cloudprober will throw an error.
	_, err := c.client.AddProbe(ctx, &proberpb.AddProbeRequest{ProbeConfig: probePb})
	return err
}

// RemoveProbe removes a probe from Cloudprober, given the probe name (located in the probe config).
// An error will be returned if there is no active probe with this probe name.
// Parameters:
// - ctx: Context used for cancelling RPCs.
// - probeName: This must be the name of an active probe in Cloudprober.
// Returns:
// - error:
//	   - Code: 3, InvalidArgument: probeName is an empty string
//	   - Code: 5,  NotFound: cannot find a probe matching this probe name
func (c *CloudproberClient) RemoveProbe(ctx context.Context, probeName string) error {
	c.clientMux.Lock()
	defer c.clientMux.Unlock()

	_, err := c.client.RemoveProbe(ctx, &proberpb.RemoveProbeRequest{ProbeName: &probeName})
	return err
}

// ListProbes() returns a (stable) sorted array of active probes from Cloudprober.
// Parameters:
// - ctx: Context used for cancelling RPCs.
// Returns:
// - probesList: This is the array of active probes in Cloudprober.
//			 If this is empty, there are no active probes.
// - Error:
//	   - See Cloudprober ListProbes() RPC for details on an error.
func (c *CloudproberClient) ListProbes(ctx context.Context) ([]*proberpb.Probe, error) {
	listProbesResp, err := c.client.ListProbes(ctx, &proberpb.ListProbesRequest{})
	if err != nil {
		return nil, err
	}

	probesList := listProbesResp.GetProbe()

	// Sorts the probes by their extension number as they are in the format:
	// "testExtension" + <number>, e.g. "testExtension0".
	sort.SliceStable(probesList[:], func(i, j int) bool {
		probeNum0, _ := strconv.Atoi(probesList[i].GetName()[len("testExtension"):])
		probeNum1, _ := strconv.Atoi(probesList[j].GetName()[len("testExtension"):])
		return probeNum0 < probeNum1
	})

	return probesList, nil
}
