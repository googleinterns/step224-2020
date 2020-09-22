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
	"github.com/google/cloudprober/probes/options"
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"

	cpmetrics "github.com/google/cloudprober/metrics"
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

	fmt.Println("INIT COMPLETE")

	return nil
}

// Start runs the probe indefinitely, unless cancelled, at the configured interval.
// Probe metrics will be sent via the metricChan at the end of the probe run.
// This is a required method to implement the cloudprober.Probes.Probe interface.
// Arguments:
//	- ctx: context provided for cancelling probe.
//	- metricChan: unidirectional channel used for sending metrics to be surfaced.
func (mp *Probe) Start(ctx context.Context, metricChan chan *cpmetrics.EventMetrics) {
	fmt.Println("Interval: ", mp.interval())
	probeTicker := time.NewTicker(mp.interval())
	fmt.Println("START")

	for {
		select {
		case <-ctx.Done():
			probeTicker.Stop()
			return
		case <-probeTicker.C:
			fmt.Println("RUN PROBE")
			mp.runProbe(ctx, metricChan)
		}
	}
}

// reportMetrics sends the metrics recorded in the current probe run to Cloudprober.
// Arguments:
//	- run: metrics from a probe run on a target.
//	- metricChan: metric channel passed from Cloudprober.
func reportMetrics(run *metrics.Metrics, metricChan chan<- *cpmetrics.EventMetrics) {
	exitStatus := metrics.ProbeFailed
	dist, ok := run.ProbeOpLatency[metrics.TotalProbeRun][exitStatus].Metric("hermes_probe_latency_seconds").(*cpmetrics.Distribution)
	if !ok {
		fmt.Println("Type conversion failed for Value -> Distribution")
	}
	data := dist.Data()
	fmt.Printf("\n\nProbe Run, status: %v, count: %d\n\n", exitStatus, data.Count)
	for _, op := range run.ProbeOpLatency {
		for _, m := range op {
			m.Timestamp = time.Now()
			metricChan <- m
		}
	}
	dist, ok = run.ProbeOpLatency[metrics.TotalProbeRun][exitStatus].Metric("hermes_probe_latency_seconds").(*cpmetrics.Distribution)
	if !ok {
		fmt.Println("Type conversion failed for Value -> Distribution")
	}
	data = dist.Data()
	fmt.Printf("\n\nProbe Run, status: %v, count: %d\n\n", exitStatus, data.Count)

	for _, call := range run.APICallLatency {
		for _, m := range call {
			m.Timestamp = time.Now()
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
	fmt.Println("RUN PROBE INNER")
	fmt.Println("Targets length: ", len(mp.targets))
	for _, t := range mp.targets {
		wg.Add(1)
		go func(t *Target) {
			defer wg.Done()
			var err error
			fmt.Println("Pre-Metrics creation")
			if t.LatencyMetrics == nil {
				if t.LatencyMetrics, err = metrics.NewMetrics(mp.config, t.Target); err != nil {
					mp.logger.Errorf(err.Error())
					fmt.Printf("Err: %v", err.Error())
					return
				}
				fmt.Println("METRICS CREATED")
			}
			probeCtx, _ := context.WithDeadline(ctx, time.Now().Add(mp.interval()))
			start := time.Now()
			exitStatus, err := mp.runProbeForTarget(probeCtx, t)
			fmt.Printf("Exit status: %v, err: %v\n\n", exitStatus, err)
			dist, ok := t.LatencyMetrics.ProbeOpLatency[metrics.TotalProbeRun][exitStatus].Metric("hermes_probe_latency_seconds").(*cpmetrics.Distribution)
			if !ok {
				fmt.Println("Type conversion failed for Value -> Distribution")
			}
			data := dist.Data()
			fmt.Printf("Probe Run, status: %v, count: %d", exitStatus, data.Count)
			if err != nil {
				mp.logger.Errorf(err.Error())
				t.LatencyMetrics.ProbeOpLatency[metrics.TotalProbeRun][exitStatus].Metric("hermes_probe_latency_seconds").AddFloat64(time.Now().Sub(start).Seconds())
			} else {
				t.LatencyMetrics.ProbeOpLatency[metrics.TotalProbeRun][metrics.Success].Metric("hermes_probe_latency_seconds").AddFloat64(time.Now().Sub(start).Seconds())
			}
			dist, ok = t.LatencyMetrics.ProbeOpLatency[metrics.TotalProbeRun][exitStatus].Metric("hermes_probe_latency_seconds").(*cpmetrics.Distribution)
			if !ok {
				fmt.Println("Type conversion failed for Value -> Distribution")
			}
			data = dist.Data()
			fmt.Printf("Probe Run, status: %v, count: %d", exitStatus, data.Count)
			// fmt.Printf("Probe run output: %v\n\n", t.LatencyMetrics.ProbeOpLatency[metrics.TotalProbeRun][exitStatus])
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

	const fileSize = 10

	c, err := storage.NewClient(ctx)
	if err != nil {
		return metrics.ProbeFailed, err
	}
	client := stiface.AdaptClient(c)

	// SETUP HERE
	for id := 1; id < 50; id++ {
		if target.Journal.Filenames[int32(id)] == "" {
			// File does not exist
			if err := CreateFile(ctx, target, int32(id), fileSize, client, mp.logger); err != nil {
				return metrics.ProbeFailed, err
			}
		}
	}

	id, err := DeleteFile(ctx, PickFileToDelete(), target, client, mp.logger)
	fmt.Printf("\n\nDeleteFile: ID: %d, err: %v\n", id, err)
	if err != nil {
		return metrics.ProbeFailed, err
	}

	if err := CreateFile(ctx, target, int32(id), fileSize, client, mp.logger); err != nil {
		fmt.Printf("CreateFile: ID: %d, err: %v\n\n", id, err)
		return metrics.ProbeFailed, err
	}
	fmt.Printf("CreateFile: ID: %d, err: %v\n\n", id, err)

	return metrics.Success, nil
}
