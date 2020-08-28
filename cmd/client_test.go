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
// Client_test tests the client.go file to ensure that the gRPC client for connecting
// to Cloudprober works as expected. This client is used for adding and removing probes.

package cmd

import (
	"github.com/golang/glog"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/examples/extensions/myprober/myprobe"
	probes_configpb "github.com/google/cloudprober/probes/proto"
	targetspb "github.com/google/cloudprober/targets/proto"
	"github.com/google/cloudprober"
	"github.com/google/cloudprober/web"

	"github.com/googleinterns/step224-2020/hermes"
)

const (
	cpCfg     string = "grpc_port:9314" // Cloudprober config
	rpcServer string = "localhost:9314" // Cloudprober gRPC server address

	testProbe0 string = "testExtension0"
	testProbe1 string = "testExtension1"
)

//setupHermes initialises an instance of Cloudprober so that RPCs can be sent to it.
func setupHermes(hermes *hermes.Hermes) {
	err := hermes.InitialiseCloudproberFromConfig(cpCfg)
	if err != nil {
		glog.Fatalf("Cloudprober could not be initialised, err: %v", err)
	}
}

// TestMain is used to setup Cloudprober before tests are run.
// This is necessary for there to be a gRPC server to make RPCs to.
func TestMain(m *testing.M) {
	test_hermes := &hermes.Hermes{}

	setupHermes(test_hermes)              // Set up Hermes and initialise Cloudprober
	defer hermes.CancelCloudprober() // Cancel Cloudprober after tests are completed

	// Sets up web UI for cloudprober.
	web.Init()

	// Start running Cloudprober instance from Hermes context
	cloudprober.Start(hermes.Ctx)

	code := m.Run()
	os.Exit(code)
}

// generationRedisProbeDef generates an inline probe definition of the RedisProbe probe extension.
// The RedisProbe extension is supplied as an example extension in Cloudprober.
func generateRedisProbeDef(name string) *probes_configpb.ProbeDef {
	probeDef := &probes_configpb.ProbeDef{ // Create probe def
		Name: proto.String(name),
		Type: probes_configpb.ProbeDef_EXTENSION.Enum(),
		Targets: &targetspb.TargetsDef{
			Type: &targetspb.TargetsDef_DummyTargets{},
		},
	}

	// Add RedisProbe extension to probeDef
	op := myprobe.ProbeConf_Op.Enum(myprobe.ProbeConf_SET)
	proto.SetExtension(probeDef, myprobe.E_RedisProbe, &myprobe.ProbeConf{Op: op, Key: proto.String("testkey"), Value: proto.String("testval")})
	return probeDef
}

// TestRegisterAndAddProbe tests that an extension probe type can be registered and added to Cloudprober without error.
func TestRegisterAndAddProbe(t *testing.T) {
	cpc := &CloudproberClient{}
	cpc.InitClient(rpcServer)

	// Create a probe and then register and add it to the prober
	err := cpc.RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(testProbe0), &myprobe.Probe{})
	if err != nil {
		t.Errorf("Probe not correctly registered and added to Cloudprober, error: %v", err)
	}

	// Check if the probe was added correctly
	resp, errResp := cpc.ListProbes()
	if errResp != nil {
		t.Errorf("ListProbes() failed, error: %v", errResp)
	}

	respProbes := resp.GetProbe()

	if respProbes[0].GetName() != testProbe0 {
		t.Errorf("Expected probe %s to be running, got probe %s running", testProbe0, respProbes[0].GetName())
	}

	// Try adding a second probe to see if the registration is still in place
	err = cpc.RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(testProbe1), &myprobe.Probe{})

	// Check if the probe was added correctly
	resp, errResp = cpc.ListProbes()
	if errResp != nil {
		t.Errorf("ListProbes() failed, error: %v", errResp)
	}

	respProbes = resp.GetProbe()

	if respProbes[1].GetName() != testProbe1 {
		t.Errorf("Expected probe %s to be running, got probe %s running", testProbe0, respProbes[1].GetName())
	}

	err = cpc.RemoveProbe(testProbe0)
	if err != nil {
		t.Errorf("RemoveProbe() failed, error: %v", err)
	}

	err = cpc.RemoveProbe(testProbe1)
	if err != nil {
		t.Errorf("RemoveProbe() failed, error: %v", err)
	}

	err = cpc.CloseConn()
	if err != nil {
		t.Errorf("CloseConn() failed, error: %v", err)
	}
}

// TestRemoveProbe tests that the RemoveProbe() method removes a probe from Cloudprober
func TestRemoveProbe(t *testing.T) {
	cpc := &CloudproberClient{}
	cpc.InitClient(rpcServer)

	// Create a probe and then register and add it to the prober
	err := cpc.RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(testProbe0), &myprobe.Probe{})
	if err != nil {
		t.Errorf("Probe not correctly registered and added to Cloudprober, error: %v", err)
	}

	// Remove probe
	err = cpc.RemoveProbe(testProbe0)
	if err != nil {
		t.Errorf("RemoveProbe() failed, error: %v", err)
	}

	// Check if the probe was removed correctly
	resp, errResp := cpc.ListProbes()
	if errResp != nil {
		t.Errorf("ListProbes() failed, error: %v", errResp)
	}

	respProbes := resp.GetProbe()

	if len(respProbes) != 0 {
		t.Errorf("Expected no probes to be running, got %d probes running; active probes running: %v", len(respProbes), respProbes)
	}

	// Close connection to gRPC server
	err = cpc.CloseConn()
	if err != nil {
		t.Errorf("CloseConn() failed, error: %v", err)
	}
}

// TestListProbes tests that the ListProbes() method returns all active probes in Cloudprober
func TestListProbes(t *testing.T) {
	// Create a probe and then register and add it to the prober
	cpc := &CloudproberClient{}
	cpc.InitClient(rpcServer)

	// Create a probe and then register and add it to the prober
	err := cpc.RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(testProbe0), &myprobe.Probe{})
	if err != nil {
		t.Errorf("Probe not correctly registered and added to Cloudprober, error: %v", err)
	}

	// Check if the probe was added correctly
	resp, errResp := cpc.ListProbes()
	if errResp != nil {
		t.Errorf("ListProbes() failed, error: %v", errResp)
	}

	respProbes := resp.GetProbe()

	if respProbes[0].GetName() != testProbe0 {
		t.Errorf("Expected probe %s to be running, got probe %s running", testProbe0, respProbes[0].GetName())
	}

	// Remove probe
	err = cpc.RemoveProbe(testProbe0)
	if err != nil {
		t.Errorf("RemoveProbe() failed, error: %v", err)
	}

	// Check if the probe was removed correctly
	resp, errResp = cpc.ListProbes()
	if errResp != nil {
		t.Errorf("ListProbes() failed, error: %v", errResp)
	}

	respProbes = resp.GetProbe()

	if len(respProbes) != 0 {
		t.Errorf("Expected no probes to be running, got %d probes running; active probes running: %v", len(respProbes), respProbes)
	}

	err = cpc.CloseConn()
	if err != nil {
		t.Errorf("CloseConn() failed, error: %v", err)
	}
}
