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
	"sync"

	proberpb "github.com/google/cloudprober/prober/proto"
	cpprobes "github.com/google/cloudprober/probes"
	configpb "github.com/google/cloudprober/probes/proto"
	"google.golang.org/grpc"
)

type CloudproberClient struct {
	conn      *grpc.ClientConn           // Client connection for opening and closing the connection
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
//	  - gRPC Error: See https://godoc.org/google.golang.org/grpc/codes for error codes.
//		-> See https://github.com/grpc/grpc-go/blob/d25c71b54334380ff1febd25d88064b36de44b3c/clientconn.go#L123
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
		return nil
	}

	return nil
}

// CloseConn() is used for closing the client connection with the gRPC server.
// Returns:
// - error:
//	   - nil: No error occurred and the connection was successfully closed.
//	   - Code 1, ErrClientConnClosing: This operation is illegal because the client connection is already closing.
func (cpc *CloudproberClient) CloseConn() error {
	cpc.clientMux.Lock()
	cpc.clientMux.Unlock()

	err := cpc.conn.Close() // Close the connection
	return err
}

// addProbeFromConfig() adds a probe to Cloudprober via the gRPC client.
// Parameters:
// - probePb: - This probe config must be unmarshalled before being passed as an argument.
//			 - This probe type must be registered as an extension.
// Returns:
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
func (cpc *CloudproberClient) addProbeFromConfig(probePb *configpb.ProbeDef) error {
	// No mutex locking handled here as this is a private method called by
	// a public method which handles locking.
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
func (cpc *CloudproberClient) RegisterAndAddProbe(extensionNumber int, probePb *configpb.ProbeDef, hermesProbeToAdd cpprobes.Probe) error {
	cpc.clientMux.Lock()
	defer cpc.clientMux.Unlock()

	if cpc.client == nil {
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
func (cpc *CloudproberClient) RemoveProbe(probeName string) error {
	cpc.clientMux.Lock()
	defer cpc.clientMux.Unlock()

	if cpc.client == nil {
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
func (cpc *CloudproberClient) ListProbes() (*proberpb.ListProbesResponse, error) {
	cpc.clientMux.Lock()
	defer cpc.clientMux.Unlock()

	if cpc.client == nil {
		err := fmt.Errorf("ClientNotInitialised: cannot register and add a probe when client is not initialised")
		return nil, err
	}

	return cpc.client.ListProbes(context.Background(), &proberpb.ListProbesRequest{}) // Submit ListProbes() rpc.
}
