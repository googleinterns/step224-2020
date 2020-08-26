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
	"context"
	"os"
	"testing"

	cp "github.com/googleinterns/step224-2020/cloudprober"
	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes/options"
	probes_configpb "github.com/google/cloudprober/probes/proto"
	testdatapb "github.com/google/cloudprober/probes/testdata"
	targetspb "github.com/google/cloudprober/targets/proto"
)

var (
	cpCfg string = "grpc_port:9314"
)

// testProbe implements the probes.Probe interface, while providing
// facilities to examine the probe status for the purpose of testing.
type testProbe struct {
	intialized      bool
	runningStatusCh chan bool
}

func (p *testProbe) Init(name string, opts *options.Options) error {
	p.intialized = true
	p.runningStatusCh = make(chan bool)
	return nil
}

func (p *testProbe) Start(ctx context.Context, dataChan chan *metrics.EventMetrics) {
	p.runningStatusCh <- true

	// If context is done (used to stop a running probe before removing it),
	// change probe state to not-running.
	<-ctx.Done()
	p.runningStatusCh <- false
	close(p.runningStatusCh)
}

func setupHermes() {
	cp.InitialiseCloudproberFromConfig(cpCfg)
}

func TestMain(m *testing.M) {
	setupHermes()
	code := m.Run()
	os.Exit(code)
}

// generateExtensionProbeDef generates an inline probe definition of the FancyProbe probe extension.
// The FancyProbe extension is supplied as testdata in Cloudprober.
func generateExtensionProbeDef(name string) *probes_configpb.ProbeDef {
	probeDef := &probes_configpb.ProbeDef{
		Name: proto.String(name),
		Type: probes_configpb.ProbeDef_EXTENSION.Enum(),
		Targets: &targetspb.TargetsDef{
			Type: &targetspb.TargetsDef_DummyTargets{},
		},
	}

	proto.SetExtension(probeDef, testdatapb.E_FancyProbe, &testdatapb.FancyProbe{Name: proto.String("fancy")})
	return probeDef
}

// TestRegisterAndAddProbe tests that an extension probe type can be registered and added to Cloudprober without error.
func TestRegisterAndAddProbe(t *testing.T) {
	// Create a probe and then register and add it to the prober
	probeName1 := "testExtension1"
	RegisterAndAddProbe(int(testdatapb.E_FancyProbe.Field), generateExtensionProbeDef(probeName1), &testProbe{})

	// Check if the probe was added correctly
	resp := ListProbes()
	respProbes := resp.GetProbe()

	if respProbes[0].GetName() != probeName1 {
		t.Errorf("Expected probe %s to be running, got probe %s running", probeName1, respProbes[0].GetName())
	}

	RemoveProbe(probeName1)

	// Try adding a second probe to see if the registration is still in place
	probeName2 := "testExtension2"
	RegisterAndAddProbe(int(testdatapb.E_FancyProbe.Field), generateExtensionProbeDef(probeName2), &testProbe{})

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
        RegisterAndAddProbe(int(testdatapb.E_FancyProbe.Field), generateExtensionProbeDef(probeName1), &testProbe{})

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
	RegisterAndAddProbe(int(testdatapb.E_FancyProbe.Field), generateExtensionProbeDef(probeName1), &testProbe{})

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
		t.Errorf("Expected no probes to be running, got %d probes running", len(respProbes))
		t.Error("Probes running: ", respProbes)
	}
}
