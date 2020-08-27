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
	"sync"
	"fmt"

	proberpb "github.com/google/cloudprober/prober/proto"
	cpprobes "github.com/google/cloudprober/probes"
	configpb "github.com/google/cloudprober/probes/proto"
	"google.golang.org/grpc"
)

type CloudproberClient struct {
	conn  *grpc.ClientConn // Client connection for opening and closing the connection
	client    proberpb.CloudproberClient // Cloudprober gRPC client retained between RPC calls
	clientMux sync.Mutex                 // Mutex lock for making CloudproberClient safe for concurrent use
}

// InitClient() establishes an active client connection to the Cloudprober gRPC server.
// If there is already an active connection, it will do nothing and return nil.
// If there is not an active connection, it will make one and assign a Cloudprober client
// to the CloudproberClient.client field.
// Parameters:
// - rpcServer: Pass the rpcServer address and port as a string formatted as: "hostname:port", e.g. "localhost:1234"
// Returns:
// - error:
//	  - nil: No error occurred and the client was successfully initialised.
//	  - TODO(): Add documentation on GRPC dial() errors.
func (cpc *CloudproberClient) InitClient(rpcServer string) error {
	cpc.clientMux.Lock()
	defer cpc.clientMux.Unlock()

	if cpc.client == nil { // If there is not an active client connection, make one.
		conn, err := grpc.Dial(rpcServer, grpc.WithInsecure()) // Make a connection
		if err != nil {
			dialErr := fmt.Errorf("GRPC_DIAL_ERROR: %v", err)
			return dialErr
		}
		cpc.conn = conn
		cpc.client = proberpb.NewCloudproberClient(conn) // Create a new Cloudprober gRPC Client
	}

	return nil
}

// CloseConn() is used for closing the client connection with the gRPC server.
// Returns:
// - error:
//	   - nil: No error occurred and the connection was successfully closed.
//	   - Code 1, ErrClientConnClosing: This operation is illegal because the client connection is already closing.
func (cpc *CloudproberClient) CloseConn() error {
	err := cpc.conn.Close() // Close the connection
	return err
}

// addProbeFromConfig() adds a probe to Cloudprober via the gRPC client.
// Parameters:
// - probePb: - This probe config must be unmarshalled before being passed as an argument.
//			 - This probe type must be registered as an extension.
// Returns:
// - error:
func (cpc *CloudproberClient) addProbeFromConfig(probePb *configpb.ProbeDef) error {
	cpc.clientMux.Lock()
	defer cpc.clientMux.Unlock()

	_, err := cpc.client.AddProbe(context.Background(), &proberpb.AddProbeRequest{ProbeConfig: probePb}) // Adds the probe to Cloudprober
	return err
}

// RegisterAndAddProbe registers a probe type and adds the probe to Cloudprober to be run.
// Parameters:
// - probePb: - This probe config must be unmarshalled before being passed as an argument.
// Returns:
// - error:
//	   - ClientNotInitialised: the gRPC client of CloudproberClient is not registered.
//	     -> Solution: Call the InitClient(rpcServer string) method.
func (cpc *CloudproberClient)  RegisterAndAddProbe(extensionNumber int, probePb *configpb.ProbeDef, hermesProbeToAdd cpprobes.Probe) error {
	cpc.clientMux.Lock()
	defer cpc.clientMux.Unlock()

	if cpc.client == nil{
		err := fmt.Errorf("ClientNotInitialised: cannot register and add a probe when client is not initialised")
		return err
	}

	cpprobes.RegisterProbeType(extensionNumber, func() cpprobes.Probe { return hermesProbeToAdd }) // First, register the probe as an extension with Cloudprober.

	// Add the probe to Cloudprober
	// The probe will be scheduled and run by Cloudprober
	// Adding a probe will consume the probe type registration.
	// If you want to add multiple probes, you must register and add each one individually.
	// Only one extension type can be registered at any given time.
	// If more are registered, Cloudprober will throw an error.
	return cpc.addProbeFromConfig(probePb)
}

// RemoveProbe removes a probe from Cloudprober, given the probe name (located in the probe config).
// An error will be returned if there is no active probe with this probe name.
// Parameters:
// - probeName: This must be the name of an active probe in Cloudprober.
// Returns:
// - error:
//	   - Code: 3, InvalidArgument: probeName is an empty string
//	   - Code: 5,  NotFound: cannot find a probe matching this probe name
//	   - ClientNotInitialised: the gRPC client of CloudproberClient is not registered.
//	     -> Solution: Call the InitClient(rpcServer string) method.
func (cpc *CloudproberClient)  RemoveProbe(probeName string) error {
	cpc.clientMux.Lock()
	defer cpc.clientMux.Unlock()

	if cpc.client == nil{
		err := fmt.Errorf("ClientNotInitialised: cannot register and add a probe when client is not initialised")
		return err
	}

	probeToRemove := &probeName // Need to use a string pointer for RemoveProbeRequest{}

	_, err := cpc.client.RemoveProbe(context.Background(), &proberpb.RemoveProbeRequest{ProbeName: probeToRemove}) // Remove probe from Cloudprober
	return err
}

// ListProbes() returns a list of active probes from Cloudprober.
// Returns:
// - ListProbesResponse: This is the list of active probes in Cloudprober.
//			 If this is empty, there are no active probes.
// - Error:
//	   - ClientNotInitialised: the gRPC client of CloudproberClient is not registered.
//	     -> Solution: Call the InitClient(rpcServer string) method.
func (cpc *CloudproberClient)  ListProbes() (*proberpb.ListProbesResponse, error) {
	cpc.clientMux.Lock()
	defer cpc.clientMux.Unlock()

	if cpc.client == nil{
		err := fmt.Errorf("ClientNotInitialised: cannot register and add a probe when client is not initialised")
		return nil, err
	}

	return cpc.client.ListProbes(context.Background(), &proberpb.ListProbesRequest{}) // Submit ListProbes() rpc.
}
