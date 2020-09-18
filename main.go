// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Evan Spendlove, GitHub: @evanSpendlove.
//
// Main program loop for Hermes. This initialises Cloudprober so that Hermes can
// interact with it through gRPCs.

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/google/cloudprober"
	"github.com/google/cloudprober/web"

	"github.com/golang/protobuf/proto"
	"github.com/googleinterns/step224-2020/client"
	"github.com/googleinterns/step224-2020/hermes/probe"

	metricpb "github.com/google/cloudprober/metrics/proto"
	probes_configpb "github.com/google/cloudprober/probes/proto"
	targetspb "github.com/google/cloudprober/targets/proto"
	monitorpb "github.com/googleinterns/step224-2020/config/proto"
)

var (
	rpcPort = flag.Int("rpc_port", 9314, "The port that the gRPC server of Cloudprober will run on.")
)

func main() {
	flag.Parse()

	if err := cloudprober.InitFromConfig(buildConfig()); err != nil {
		glog.Exitf("cloudprober could not be initialised from config: grpc_port: %d, err:%v", *rpcPort, err)
	}

	// Sets up web UI for cloudprober.
	web.Init()

	cloudprober.Start(context.Background())

	client, err := client.NewClient("localhost:9314")
	if err != nil {
		glog.Errorf("Cloudprober gRPC Client could not be initialised: %v", err)
	}

	if err := client.RegisterAndAddProbe(context.Background(), 200, genTestConfig("testProbe1"), &probe.Probe{}); err != nil {
		glog.Errorf("Probe not correctly registered and added to Cloudprober, expected error: %v, got: %v", nil, err)
	}

	// Wait forever
	select {}
}

// buildConfig() builds the configuration details for Cloudprober based on the flag contents.
// Returns:
// - string: Returns the configuration details for Cloudprober as a string.
func buildConfig() string {
	return fmt.Sprintf("grpc_port: %d", *rpcPort)
}

// genTestConfig generates a test HermesProbeDef proto config.
// Arguments:
//	- name: pass the name for this probe instance.
// Returns:
//	- hermesExtension: returns the HermesProbeDef extension.
func genTestConfig(name string) *probes_configpb.ProbeDef {
	probeDef := &probes_configpb.ProbeDef{
		Name: proto.String(name),
		Type: probes_configpb.ProbeDef_EXTENSION.Enum(),
		Targets: &targetspb.TargetsDef{
			Type: &targetspb.TargetsDef_DummyTargets{},
		},
	}

	hermesExtension := &monitorpb.HermesProbeDef{
		ProbeName: proto.String(name),
		Targets: []*monitorpb.Target{
			&monitorpb.Target{
				Name:                   "hermes",
				TargetSystem:           monitorpb.Target_GOOGLE_CLOUD_STORAGE,
				TotalSpaceAllocatedMib: int64(100),
				BucketName:             "test_bucket_5",
			},
		},
		TargetSystem: monitorpb.HermesProbeDef_GCS.Enum(),
		IntervalSec:  proto.Int32(5),
		TimeoutSec:   proto.Int32(60),
		ProbeLatencyDistribution: &metricpb.Dist{
			Buckets: &metricpb.Dist_ExplicitBuckets{
				ExplicitBuckets: "0.1,0.2,0.4,0.6,0.8,1.6,3.2,6.4,12.8,1000",
			},
		},
		ApiCallLatencyDistribution: &metricpb.Dist{
			Buckets: &metricpb.Dist_ExplicitBuckets{
				ExplicitBuckets: "0.000000002,0.000000004,0.000000008,0.000000016,0.00000032,0.000000064,0.000000128,100",
			},
		},
	}
	proto.SetExtension(probeDef, monitorpb.E_HermesProbeDef_HermesProbeDef, hermesExtension)
	return probeDef
}

const (
	bucketName = "test_bucket_5"
)

// getTargetPb generates a target proto config.
// Returns:
//	- target: target proto config
func genTargetPb() *monitorpb.Target {
	return &monitorpb.Target{
		Name:                   "hermes",
		TargetSystem:           monitorpb.Target_GOOGLE_CLOUD_STORAGE,
		TotalSpaceAllocatedMib: int64(100),
		BucketName:             bucketName,
	}
}
