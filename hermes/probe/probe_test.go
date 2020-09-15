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
// Authors:
// - Alicja Kwiecinska (kwiecinskaa@google.com) github: alicjakwie
// - Evan Spendlove, GitHub: evanSpendlove.
//
// Monitor_test tests the functions used in the monitor probe, as well
// as the probes integration with Cloudprober.

package probe

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/cloudprober/logger"
	"github.com/google/cloudprober/probes/options"

	metricpb "github.com/google/cloudprober/metrics/proto"
	probes_configpb "github.com/google/cloudprober/probes/proto"
	monitorpb "github.com/googleinterns/step224-2020/config/proto"
)

// GenTestConfig generates a test HermesProbeDef proto config for
// initialising and running a MonitorProbe.
// Arguments:
//	- name: pass the name for this probe instance.
// Returns:
//	- probeDef: a test probe config with a HermesProbeDef extension set.
//	- hermesExtension: returns the HermesProbeDef extension.
func GenTestConfig(name string) (*probes_configpb.ProbeDef, *monitorpb.HermesProbeDef) {
	probeDef := &probes_configpb.ProbeDef{
		Name: proto.String(name),
		Type: probes_configpb.ProbeDef_EXTENSION.Enum(),
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
		IntervalSec:  proto.Int32(3600),
		TimeoutSec:   proto.Int32(60),
		ProbeLatencyDistribution: &metricpb.Dist{
			Buckets: &metricpb.Dist_ExplicitBuckets{
				ExplicitBuckets: "0.1, 0.2, 0.4, 0.6 0.8, 1.6, 3.2, 6.4, 12.8, 1000",
			},
		},
		ApiCallLatencyDistribution: &metricpb.Dist{
			Buckets: &metricpb.Dist_ExplicitBuckets{
				ExplicitBuckets: "0.000000002, 0.000000004, 0.000000008, 0.000000016, 0.00000032, 0.000000064, 0.000000128, 100",
			},
		},
	}
	proto.SetExtension(probeDef, monitorpb.E_HermesProbeDef_HermesProbeDef, hermesExtension)
	return probeDef, hermesExtension
}

// GenOptsFromConfig generates the options for the MonitorProbe
// from the probe config passed.
// Arguments:
//	- t: testing object for throwing errors
//	- cfg: probe config.
// Returns:
//	- opts: options for a probe generated from the probe config.
func GenOptsFromConfig(t *testing.T, cfg *monitorpb.HermesProbeDef) *options.Options {
	opts := &options.Options{
		Interval: time.Duration(cfg.GetIntervalSec()) * time.Second,
		Timeout:  time.Duration(cfg.GetTimeoutSec()) * time.Second,
	}

	var err error
	if opts.Logger, err = logger.NewCloudproberLog(cfg.GetProbeName()); err != nil {
		t.Fatalf("error in initializing logger for the probe (%s): %v", cfg.GetProbeName(), err)
	}

	opts.ProbeConf = cfg

	return opts
}

func TestInit(t *testing.T) {
	wantName := "testProbe1"
	mp := &Probe{}
	_, wantConfig := GenTestConfig(wantName)
	wantOpts := GenOptsFromConfig(t, wantConfig)

	if err := mp.Init(wantName, wantOpts); err != nil {
		t.Errorf("error when calling Init() on MonitorProbe object: wanted %v, got %v for error return value", nil, err)
	}
}

// TODO(evanSpendlove): Add more tests for monitor.go methods.
