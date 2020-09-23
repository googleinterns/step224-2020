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
// probe implements a probe that monitors a storage system using the
// Hermes algorithm.

// Package probe implements the probe that Hermes uses to monitor
// a storage system.
package probe

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/cloudprober/logger"
	cpmetrics "github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes/options"
	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
	"github.com/googleinterns/step224-2020/hermes/probe/target"

	probepb "github.com/googleinterns/step224-2020/config/proto"
	journalpb "github.com/googleinterns/step224-2020/hermes/proto"
)

const (
	// TODO(evanSpendlove): Refactor to use constant from metrics.go
	probeLatency = "hermes_probe_latency_seconds"
)

// Probe holds aggregate information about all probe runs, per-target.
// It also holds the config and options used to initialise the probe.
type Probe struct {
	name    string
	config  *probepb.HermesProbeDef
	targets []*target.Target
	opts    *options.Options
	logger  *logger.Logger
}

// interval returns the probing interval in seconds as a time.Duration.
// Returns:
//	- time.Duration: returns the probing interval
func (p *Probe) interval() time.Duration {
	return time.Duration(p.config.GetIntervalSec()) * time.Second
}

// timeout returns the probe timeout in seconds as a time.Duration.
// Returns:
//	- time.Duration: returns the probe timeout
func (p *Probe) timeout() time.Duration {
	return time.Duration(p.config.GetTimeoutSec()) * time.Second
}

// Init initializes the probe with the given parameters.
// This is a required method to iplement the cloudprober.Probes.Probe interface.
func (p *Probe) Init(name string, opts *options.Options) error {
	// Cast from Cloudprober defined interface to HermesProbeDef.
	// This allows for accessing the variables and methods of a HermesProbeDef object.
	conf, ok := opts.ProbeConf.(*probepb.HermesProbeDef)
	if !ok {
		return fmt.Errorf("invalid argument: opts.ProbeConf is not of type *probepb.HermesProbeDef")
	}
	p.name = name
	p.config = conf

	for _, t := range p.config.GetTargets() {
		p.targets = append(p.targets, &target.Target{
			Target: t,
			Journal: &journalpb.StateJournal{
				Intent:    &journalpb.Intent{},
				Filenames: make(map[int32]string),
			},
		})
	}
	p.opts = opts
	p.logger = opts.Logger

	return nil
}

// Start runs the probe indefinitely, unless cancelled, at the configured interval.
// Probe metrics will be sent via the metricChan at the end of the probe run.
// This is a required method to iplement the cloudprober.Probes.Probe interface.
// Arguments:
//	- ctx: context provided for cancelling probe.
//	- metricChan: bidirectional channel used for sending metrics to be surfaced.
//		- Must be bidirectional to satisfy cloudprober.Probes.Probe interface.
func (p *Probe) Start(ctx context.Context, metricChan chan *cpmetrics.EventMetrics) {
	probeTicker := time.NewTicker(p.interval())

	for {
		select {
		case <-ctx.Done():
			probeTicker.Stop()
			return
		case <-probeTicker.C:
			p.runProbe(ctx, metricChan)
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
			m.Timestamp = time.Now()
			metricChan <- m
		}
	}

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
func (p *Probe) runProbe(ctx context.Context, metricChan chan<- *cpmetrics.EventMetrics) {
	var wg sync.WaitGroup
	for _, t := range p.targets {
		wg.Add(1)
		t := t
		go func() {
			defer wg.Done()
			if t.LatencyMetrics == nil {
				lm, err := metrics.NewMetrics(p.config, t.Target)
				if err != nil {
					p.logger.Errorf("NewMetrics(%v) failed: %v", t.Target, err)
					return
				}
				t.LatencyMetrics = lm
			}

			probeCtx, _ := context.WithDeadline(ctx, time.Now().Add(p.interval()))
			// TODO(evanSpendlove): Refactor to use closure func from metrics.go in metrics PR.
			start := time.Now()
			status, err := p.runProbeForTarget(probeCtx, t)
			if err != nil {
				p.logger.Errorf(err.Error())
			}

			t.LatencyMetrics.ProbeOpLatency[metrics.TotalProbeRun][status].Metric(probeLatency).AddFloat64(time.Now().Sub(start).Seconds())
			reportMetrics(t.LatencyMetrics, metricChan)
		}()
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
func (p *Probe) runProbeForTarget(ctx context.Context, target *target.Target) (metrics.ExitStatus, error) {
	// TODO(evanSpendlove): Add implementation of runProbeForTarget, i.e. Hermes probing algorithm.
	return metrics.Success, nil
}
