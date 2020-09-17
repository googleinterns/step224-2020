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
// Monitor implements a probe that monitors a storage system using the
// Hermes algorithm.

// Package probe implements the probe that Hermes uses to monitor
// a storage system.
package probe

import (
	"context"
	"fmt"
	"sync"
	"time"
	"cloud.google.com/go/storage"
	"github.com/google/cloudprober/logger"
	cpmetrics "github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes/options"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"

	probepb "github.com/googleinterns/step224-2020/config/proto"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

// Target holds all of the required information and state for a given target run.
type Target struct {
	// Target stores the proto config for the target to be probed.
	Target *probepb.Target

	// Journal stores the state of MonitorProbe as a combination of a next operation intent enum and a filenames map.
	Journal *journalpb.StateJournal

	// LatencyMetrics stores the api call and probe operation latency for a given target run.
	// Metrics are stored with additional labels to record operation type and exit status.
	LatencyMetrics *metrics.Metrics
}

// Probe holds aggregate information about all probe runs, per-target.
// It also holds the config and options used to initialise the probe.
type Probe struct {
	name    string
	config  *probepb.HermesProbeDef
	targets []*Target
	opts    *options.Options
	logger  *logger.Logger
}

// interval returns the probing interval in seconds as a time.Duration.
// Returns:
//	- time.Duration: returns the probing interval
func (mp *Probe) interval() time.Duration {
	return time.Duration(mp.config.GetIntervalSec()) * time.Second
}

// timeout returns the probe timeout in seconds as a time.Duration.
// Returns:
//	- time.Duration: returns the probe timeout
func (mp *Probe) timeout() time.Duration {
	return time.Duration(mp.config.GetTimeoutSec()) * time.Second
}

// Init initializes the probe with the given parameters.
// This is a required method to implement the cloudprober.Probes.Probe interface.
func (mp *Probe) Init(name string, opts *options.Options) error {
	// Cast from Cloudprober defined interface to HermesProbeDef.
	// This allows for accessing the variables and methods of a HermesProbeDef object.
	conf, ok := opts.ProbeConf.(*probepb.HermesProbeDef)
	if !ok {
		return fmt.Errorf("invalid argument: opts.ProbeConf is not of type *probepb.HermesProbeDef")
	}
	mp.name = name
	mp.config = conf

	for _, t := range mp.config.GetTargets() {
		mp.targets = append(mp.targets, &Target{
			Target: t,
			Journal: &journalpb.StateJournal{
				Intent:    &journalpb.Intent{},
				Filenames: make(map[int32]string),
			},
		})
	}
	mp.opts = opts
	mp.logger = opts.Logger

	return nil
}

// Start runs the probe indefinitely, unless cancelled, at the configured interval.
// Probe metrics will be sent via the metricChan at the end of the probe run.
// This is a required method to implement the cloudprober.Probes.Probe interface.
// Arguments:
//	- ctx: context provided for cancelling probe.
//	- metricChan: unidirectional channel used for sending metrics to be surfaced.
func (mp *Probe) Start(ctx context.Context, metricChan chan<- *cpmetrics.EventMetrics) {
	probeTicker := time.NewTicker(mp.interval())

	for {
		select {
		case <-ctx.Done():
			probeTicker.Stop()
			return
		case <-probeTicker.C:
			mp.runProbe(ctx, metricChan)
		}
	}
}

// reportMetrics sends the metrics recorded in the current probe run to Cloudprober.
// Arguments:
//	- run: metrics from a probe run on a target.
//	- metricChan: metric channel passed from Cloudprober.
func reportMetrics(run *metrics.Metrics, metricChan chan<- *cpmetrics.EventMetrics) {
	for _, op := range run.ProbeOpLatency {
		for _, m := range op {
			metricChan <- m
		}
	}

	for _, call := range run.APICallLatency {
		for _, m := range call {
			metricChan <- m
		}
	}
}

// runProbe runs the probe against each target, collects metrics on probe run
// and surface metrics to Cloudprober.
// Arguments:
//	- ctx: pass context to allow for complete cancellation of the probe.
//	- metricChan: pass the metrics channel for surfacing metrics to Cloudprober.
func (mp *Probe) runProbe(ctx context.Context, metricChan chan<- *cpmetrics.EventMetrics) {
	var wg sync.WaitGroup
	for _, t := range mp.targets {
		wg.Add(1)
		go func(t *Target) {
			defer wg.Done()
			var err error
			if t.LatencyMetrics, err = metrics.NewMetrics(mp.config, t.Target); err != nil {
				mp.logger.Errorf(err.Error())
				return
			}
			probeCtx, _ := context.WithDeadline(ctx, time.Now().Add(mp.interval()))
			start := time.Now()
			exitStatus, err := mp.runProbeForTarget(probeCtx, t)
			if err != nil {
				mp.logger.Errorf(err.Error())
				t.LatencyMetrics.ProbeOpLatency[metrics.TotalProbeRun][exitStatus].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
				return
			}
			t.LatencyMetrics.ProbeOpLatency[metrics.TotalProbeRun][metrics.Success].Metric("latency").AddFloat64(time.Now().Sub(start).Seconds())
			reportMetrics(t.LatencyMetrics, metricChan)
		}(t)
	}
	wg.Wait()
}

// runProbeForTarget runs the Hermes probing algorithm on a single target.
// Arguments:
//	- ctx: pass context to allow for cancellation of the probe.
//	- target: the target to be probed
// Returns:
//	- status: returns the exit status of the probe run.
//	- error: returns an error if one occurred during the probe run.
func (mp *Probe) runProbeForTarget(ctx context.Context, target *Target) (metrics.ExitStatus, error) {
	// TODO(evanSpendlove): Add implementation of runProbeForTarget, i.e. Hermes probing algorithm.
	
	return metrics.Success, nil
}

func runCreate(){
	
ctx := context.Background()
  client :=storage.NewClient()
  fileID := int32(6)
  fileSize := 50
  target := Target{
       &probepb.Target{
           Name:                   "hermes",
           TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
           TotalSpaceAllocatedMib: int64(1000),
           BucketName:             "test_bucket_probe0",
       },
       &journalpb.StateJournal{
           Filenames: make(map[int32]string),
       },
       &metrics.Metrics{},
   }
   hp := &probepb.HermesProbeDef{
       ProbeName: proto.String("createfile_test"),
       Targets: []*probepb.Target{
           &probepb.Target{
               Name:                   "hermes",
               TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
               TotalSpaceAllocatedMib: int64(100),
               BucketName:             "test_bucket_probe0",
           },
       },
       TargetSystem: probepb.HermesProbeDef_GCS.Enum(),
       IntervalSec:  proto.Int32(3600),
       TimeoutSec:   proto.Int32(60),
       ProbeLatencyDistribution: &metricpb.Dist{
           Buckets: &metricpb.Dist_ExplicitBuckets{
               ExplicitBuckets: "0.1,0.2,0.4,0.6,0.8,1.6,3.2,6.4,12.8,1",
           },
       },
       ApiCallLatencyDistribution: &metricpb.Dist{
           Buckets: &metricpb.Dist_ExplicitBuckets{
               ExplicitBuckets: "0.1,0.2,0.4,0.6,0.8,1.6,3.2,6.4,12.8,1",
           },
       },
   }
   probeTarget := &probepb.Target{
       Name:                   "hermes",
       TargetSystem:           probepb.Target_GOOGLE_CLOUD_STORAGE,
       TotalSpaceAllocatedMib: int64(100),
       BucketName:             "test_bucket_probe0",
   }
 
   var err error
   if target.LatencyMetrics, err = metrics.NewMetrics(hp, probeTarget); err != nil {
       t.Fatalf("Metric set up failed due to %s", err.Error())
   }
 
   if err := CreateFile(ctx, target, fileID, fileSize, client, nil); err != nil {
       t.Error(err)
   }	
	
	
	
}
