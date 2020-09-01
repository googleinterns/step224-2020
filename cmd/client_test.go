// Copyright 2020 Google LLC // // Licensed under the Apache License, Version 2.0 (the "License");
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
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober"
	"github.com/google/cloudprober/examples/extensions/myprober/myprobe"
	probes_configpb "github.com/google/cloudprober/probes/proto"
	targetspb "github.com/google/cloudprober/targets/proto"
	"github.com/googleinterns/step224-2020/hermes"
)

const (
	cpCfg     string = "grpc_port:9314" // Cloudprober config
	rpcServer string = "localhost:9314" // Cloudprober gRPC server address
)

// TODO: Make these proper table-driven tests.
var (
	testProbes = []string{
		"testExtension0",
		"testExtension1"}
)

func setupCloudproberAndClient() (*hermes.Hermes, *CloudproberClient) {
	testHermes := &hermes.Hermes{}
	err := testHermes.InitialiseCloudproberFromConfig(cpCfg)
	if err != nil {
		glog.Fatalf("Cloudprober could not be initialised, err: %v", err)
	}

	cloudprober.Start(testHermes.Ctx)

	client, errClient := NewClient(rpcServer)
	if errClient != nil {
		glog.Fatalf("Cloudprober gRPC Client could not be initialised: %v", errClient)
	}

	return testHermes, client
}

func teardownCloudproberAndClient(testHermes *hermes.Hermes, client *CloudproberClient) {
	// TODO: Find a safe way to close down Cloudprober and restart it.
	for i := 0; i < len(testProbes); i++ {
		err := client.RemoveProbe(testProbes[i])
		if err != nil {
			glog.Errorf("RemoveProbe() failed, expected error: %v, got %v", nil, err)
		}

	}

	client.CloseConn()

	// Need to wait for the RemoveProbe() RPCs to be processed on the server side
	time.Sleep(1 * time.Second)
}

// TODO: Replace this test with a test using a Hermes probe definition
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

	// Add RedisProbe extension to probeDef proto
	op := myprobe.ProbeConf_Op.Enum(myprobe.ProbeConf_SET)
	proto.SetExtension(probeDef, myprobe.E_RedisProbe, &myprobe.ProbeConf{Op: op, Key: proto.String("testkey"), Value: proto.String("testval")})
	return probeDef
}

// TestRegisterAndAddProbe tests that an extension probe type can be registered and added to Cloudprober without error.
func TestRegisterAndAddProbe(t *testing.T) {
	testHermes, client := setupCloudproberAndClient()
	defer teardownCloudproberAndClient(testHermes, client)

	for i := 0; i < len(testProbes); i++ {
		err := client.RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), testHermes.Ctx, generateRedisProbeDef(testProbes[i]), &myprobe.Probe{})
		if err != nil {
			t.Errorf("Probe not correctly registered and added to Cloudprober, expected error: %v, got: %v", nil, err)
		}

		// Check if the probe was added correctly
		listProbesResp, err := client.ListProbes()
		if err != nil {
			t.Errorf("ListProbes() failed, expected error: %v, got: %v", nil, err)
		}

		probesList := listProbesResp.GetProbe()
		if probesList[i].GetName() != testProbes[i] {
			t.Errorf("Expected probe %s to be running, got probe %s running", testProbes[i], probesList[i].GetName())
		}
	}
}

// TestRemoveProbe tests that the RemoveProbe() method removes a probe from Cloudprober
func TestRemoveProbe(t *testing.T) {
	testHermes, client := setupCloudproberAndClient()
	defer teardownCloudproberAndClient(testHermes, client)

	for i := 0; i < len(testProbes); i++ {
		err := client.RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), testHermes.Ctx, generateRedisProbeDef(testProbes[i]), &myprobe.Probe{})
		if err != nil {
			t.Errorf("Probe not correctly registered and added to Cloudprober, expected error: %v, got: %v", nil, err)
		}

		// Check if the probe was added correctly
		listProbesResp, err := client.ListProbes()
		if err != nil {
			t.Errorf("ListProbes() failed, expected error: %v, got: %v", nil, err)
		}

		probesList := listProbesResp.GetProbe()
		if probesList[0].GetName() != testProbes[i] {
			t.Errorf("Expected probe %s to be running, got probe %s running", testProbes[i], probesList[0].GetName())
		}

		err = client.RemoveProbe(testProbes[i])
		if err != nil {
			t.Errorf("RemoveProbe() failed, expected error: %v, got %v", nil, err)
		}

		// Check if the probe was removed correctly
		listProbesResp, err = client.ListProbes()
		if err != nil {
			t.Errorf("ListProbes() failed, expected error: %v, got: %v", nil, err)
		}

		probesList = listProbesResp.GetProbe()

		if len(probesList) != 0 {
			t.Errorf("Expected no probes to be running, got %d probes running; active probes running: %v", len(probesList), probesList)
		}
	}
}

// TestListProbes tests that the ListProbes() method returns all active probes in Cloudprober
func TestListProbes(t *testing.T) {
	testHermes, client := setupCloudproberAndClient()
	defer teardownCloudproberAndClient(testHermes, client)

	for i := 0; i < len(testProbes); i++ {
		err := client.RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), testHermes.Ctx, generateRedisProbeDef(testProbes[i]), &myprobe.Probe{})
		if err != nil {
			t.Errorf("Probe not correctly registered and added to Cloudprober, expected error: %v, got: %v", nil, err)
		}

		// Check if the probe was added correctly
		listProbesResp, err := client.ListProbes()
		if err != nil {
			t.Errorf("ListProbes() failed, expected error: %v, got: %v", nil, err)
		}

		probesList := listProbesResp.GetProbe()
		if probesList[0].GetName() != testProbes[i] {
			t.Errorf("Expected probe %s to be running, got probe %s running", testProbes[i], probesList[0].GetName())
		}

		err = client.RemoveProbe(testProbes[i])
		if err != nil {
			t.Errorf("RemoveProbe() failed, expected error: %v, got %v", nil, err)
		}

		// Check if the probe was removed correctly
		listProbesResp, err = client.ListProbes()
		if err != nil {
			t.Errorf("ListProbes() failed, expected error: %v, got: %v", nil, err)
		}

		probesList = listProbesResp.GetProbe()

		if len(probesList) != 0 {
			t.Errorf("Expected no probes to be running, got %d probes running; active probes running: %v", len(probesList), probesList)
		}
	}
}
