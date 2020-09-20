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
// Author: Evan Spendlove, GitHub: @evanSpendlove.
//
// Client_test tests the client.go file to ensure that the gRPC client for connecting
// to Cloudprober works as expected. This client is used for adding and removing probes.
//
// TODOs:
// TODO(evanSpendlove): Make these proper table-driven tests.
// TODO(evanSpendlove): Find a safe way to close down Cloudprober and restart it.
// TODO(evanSpendlove): Replace this test with a test using a Hermes probe definition

package client

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober"
	"github.com/google/cloudprober/examples/extensions/myprober/myprobe"

	probes_configpb "github.com/google/cloudprober/probes/proto"
	targetspb "github.com/google/cloudprober/targets/proto"
)

const (
	cloudproberConfig string = "grpc_port:9314"
	rpcServer         string = "localhost:9314" // Cloudprober gRPC server address
	testRunCount      int    = 20
)

func setupCloudproberAndClient(t *testing.T) *CloudproberClient {
	t.Helper()
	if err := cloudprober.InitFromConfig(cloudproberConfig); err != nil {
		t.Fatalf("Cloudprober could not be initialised: %v", err)
	}

	cloudprober.Start(context.Background())

	client, err := NewClient(rpcServer)
	if err != nil {
		t.Fatalf("Cloudprober gRPC Client could not be initialised: %v", err)
	}

	return client
}

func teardownCloudproberAndClient(t *testing.T, client *CloudproberClient) {
	// TODO(evanSpendlove): Find a safe way to close down Cloudprober and restart it.
	// Check if the probe was added correctly
	ctx := context.Background()
	t.Helper()
	probes, err := client.ListProbes(ctx)
	if err != nil {
		t.Fatalf("ListProbes() failed during teardown: %v", err)
	}

	for _, p := range probes {
		if err = client.RemoveProbe(ctx, p.GetName()); err != nil {
			t.Fatalf("RemoveProbe() failed: %v", err)
		}
	}

	client.CloseConn()
}

// generationRedisProbeDef generates an inline probe definition of the RedisProbe probe extension.
// The RedisProbe extension is supplied as an example extension in Cloudprober.
// TODO(evanSpendlove): Replace this test with a test using a Hermes probe definition
func generateRedisProbeDef(name string) *probes_configpb.ProbeDef {
	probeDef := &probes_configpb.ProbeDef{
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

func TestRegisterAndAddProbe(t *testing.T) {
	ctx := context.Background()
	client := setupCloudproberAndClient(t)
	defer teardownCloudproberAndClient(t, client)

	for i := 0; i < testRunCount; i++ {
		probeName := fmt.Sprintf("testExtension%d", i)
		if err := client.RegisterAndAddProbe(ctx, int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(probeName), &myprobe.Probe{}); err != nil {
			t.Errorf("Probe not correctly registered and added to Cloudprober: %v", err)
		}

		// Check if the probe was added correctly
		probes, err := client.ListProbes(ctx)
		if err != nil {
			t.Errorf("ListProbes() failed: %v", err)
		}

		if probes[i].GetName() != probeName {
			t.Errorf("Expected probe %s to be running, got probe %s running", probeName, probes[i].GetName())
		}
	}
}

func TestRemoveProbe(t *testing.T) {
	ctx := context.Background()
	client := setupCloudproberAndClient(t)
	defer teardownCloudproberAndClient(t, client)

	for i := 0; i < testRunCount; i++ {
		probeName := fmt.Sprintf("testExtension%d", i)
		if err := client.RegisterAndAddProbe(ctx, int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(probeName), &myprobe.Probe{}); err != nil {
			t.Errorf("Probe not correctly registered and added to Cloudprober: %v", err)
		}

		// Check if the probe was added correctly
		probes, err := client.ListProbes(ctx)
		if err != nil {
			t.Errorf("ListProbes() failed: %v", err)
		}

		if probes[0].GetName() != probeName {
			t.Errorf("Expected probe %s to be running, got probe %s running", probeName, probes[0].GetName())
		}

		if err = client.RemoveProbe(ctx, probeName); err != nil {
			t.Errorf("RemoveProbe() failed: %v", err)
		}

		// Check if the probe was removed correctly
		probes, err = client.ListProbes(ctx)
		if err != nil {
			t.Errorf("ListProbes() failed: %v", err)
		}

		if len(probes) != 0 {
			t.Errorf("Expected no probes to be running, got %d probes running; active probes running: %v", len(probes), probes)
		}
	}
}

func TestListProbes(t *testing.T) {
	ctx := context.Background()
	client := setupCloudproberAndClient(t)
	defer teardownCloudproberAndClient(t, client)

	for i := 0; i < testRunCount; i++ {
		probeName := fmt.Sprintf("testExtension%d", i)
		if err := client.RegisterAndAddProbe(ctx, int(myprobe.E_RedisProbe.Field), generateRedisProbeDef(probeName), &myprobe.Probe{}); err != nil {
			t.Errorf("Probe not correctly registered and added to Cloudprober: %v", err)
		}

		// Check if the probe was added correctly
		probes, err := client.ListProbes(ctx)
		if err != nil {
			t.Errorf("ListProbes() failed: %v", err)
		}

		if probes[0].GetName() != probeName {
			t.Errorf("Expected probe %s to be running, got probe %s running", probeName, probes[0].GetName())
		}

		if err = client.RemoveProbe(ctx, probeName); err != nil {
			t.Errorf("RemoveProbe() failed: %v", err)
		}

		// Check if the probe was removed correctly
		probes, err = client.ListProbes(ctx)
		if err != nil {
			t.Errorf("ListProbes() failed: %v", err)
		}

		if len(probes) != 0 {
			t.Errorf("Expected no probes to be running, got %d probes running; active probes running: %v", len(probes), probes)
		}
	}
}
