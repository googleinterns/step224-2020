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
	"testing"
	"fmt"
	"os"
	"context"

	myprobe "github.com/google/cloudprober/examples/extensions/myprober/myprobe"

	cp "github.com/googleinterns/step224-2020/cloudprober"
	// prober "github.com/google/cloudprober"
	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/metrics"
	// pb "github.com/google/cloudprober/prober/proto"
	"github.com/google/cloudprober/probes"
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
// Since cloudprober has to be aware of the probe type, we add testProbe to
// cloudprober as an EXTENSION probe type (done through the init() function
// below).
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

// generateTestExtensionProbeDef generates an inline probe definition of the FancyProbe probe extension.
// The FancyProbe extension is supplied as testdata in Cloudprober.
func generateTestExtensionProbeDef(name string) *probes_configpb.ProbeDef{
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

func generateTestRedisProbeDef(name string) *probes_configpb.ProbeDef{
	probeDef := &probes_configpb.ProbeDef{
		Name: proto.String(name),
		Type: probes_configpb.ProbeDef_EXTENSION.Enum(),
		Targets: &targetspb.TargetsDef{
			Type: &targetspb.TargetsDef_DummyTargets{},
		},
	}

	fmt.Println("Probe: ", myprobe.E_RedisProbe)
	op := myprobe.ProbeConf_SET
	proto.SetExtension(probeDef, myprobe.E_RedisProbe, &myprobe.ProbeConf{Op: op, Key: proto.String("testkey"), Value: proto.String("testval")})
	return probeDef
}

func TestRegisterAndAddProbe(t *testing.T) {
	// Create a probe and then register and add it to the prober
	probe := "testExtension1"
	RegisterAndAddProbe(generateTestExtensionProbeDef(probe), &testProbe{})

	// Check if the probe was added correctly
	resp := ListProbes()
	respProbes := resp.GetProbe()
	fmt.Println("Response: ", resp)

	if respProbes[0].GetName() != probe {
		t.Errorf("Expected probe %s to be running, got probe %s running", probe, respProbes[0].GetName())
	}

	probes.RegisterProbeType(200, func() probes.Probe {
		return &myprobe.Probe{}
	})

	probePath := "testExtensionProbe1.cfg"
	RegisterAndAddProbe(generateTestRedisProbeDef(probePath), &myprobe.Probe{})
	resp = ListProbes()
	fmt.Println("Response: ", resp)
}
func TestRegisterAndAddTwoProbes(t *testing.T) {
	x := 0
	x++
}

func TestRemoveProbe(t *testing.T) {
	x := 0
	x++
}

func TestListProbes(t *testing.T) {
	x := 0
	x++
}

