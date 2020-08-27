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
	"os"
	"testing"

	cp "github.com/googleinterns/step224-2020/cloudprober"
	"github.com/golang/protobuf/proto"
	probes_configpb "github.com/google/cloudprober/probes/proto"
	targetspb "github.com/google/cloudprober/targets/proto"
	"github.com/google/cloudprober/examples/extensions/myprober/myprobe"
)

var (
	cpCfg string = "grpc_port:9314"
)

//setupHermes initialises an instance of Cloudprober so that RPCs can be sent to it.
func setupHermes() {
	cp.InitialiseCloudproberFromConfig(cpCfg)
}

// TestMain is used to setup Cloudprober before tests are run.
// This is necessary for there to be a gRPC server to make RPCs to.
func TestMain(m *testing.M) {
	setupHermes()
	code := m.Run()
	os.Exit(code)
}

// generationRedisProbeDef generates an inline probe definition of the RedisProbe probe extension.
// The RedisProbe extension is supplied as an example extension in Cloudprober.
func generateRedisProbeDef(name string) *probes_configpb.ProbeDef{
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
	// Create a probe and then register and add it to the prober
	probeName1 := "testExtension1"
	RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(probeName1), &myprobe.Probe{})

	// Check if the probe was added correctly
	resp := ListProbes()
	respProbes := resp.GetProbe()

	if respProbes[0].GetName() != probeName1 {
		t.Errorf("Expected probe %s to be running, got probe %s running", probeName1, respProbes[0].GetName())
	}

	RemoveProbe(probeName1)

	// Try adding a second probe to see if the registration is still in place
	probeName2 := "testExtension2"
	RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(probeName2), &myprobe.Probe{})

	resp = ListProbes()
	respProbes = resp.GetProbe()

	if respProbes[0].GetName() != probeName2 {
		t.Errorf("Expected probe %s to be running, got probe %s running", probeName2, respProbes[0].GetName())
	}

	RemoveProbe(probeName2)
}

// TestRemoveProbe tests that the RemoveProbe() method removes a probe from Cloudprober
func TestRemoveProbe(t *testing.T) {
	// Create a probe and then register and add it to the prober
        probeName1 := "testExtension1"
	RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(probeName1), &myprobe.Probe{})

	// Remove probe
	RemoveProbe(probeName1)

	// Verify that the probe was removed
	resp := ListProbes()
	respProbes := resp.GetProbe()
	if len(respProbes) != 0 {
		t.Errorf("Expected no probes to be running, got %d probes running", len(respProbes))
		t.Error("Probes running: ", respProbes)
	}
}

// TestListProbes tests that the ListProbes() method returns all active probes in Cloudprober
func TestListProbes(t *testing.T) {
	// Create a probe and then register and add it to the prober
	probeName1 := "testExtension1"
	RegisterAndAddProbe(int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(probeName1), &myprobe.Probe{})

	// Check if the probe was added correctly
	resp := ListProbes()
	respProbes := resp.GetProbe()

	if respProbes[0].GetName() != probeName1 {
		t.Errorf("Expected probe %s to be running, got probe %s running", probeName1, respProbes[0].GetName())
	}

	// Remove probe
	RemoveProbe(probeName1)

	// Verify that the probe was removed
	resp = ListProbes()
	respProbes = resp.GetProbe()
	if len(respProbes) != 0 {
		t.Errorf("Expected no probes to be running, got %d probes running, probes running: %v", len(respProbes), respProbes)
	}
}
