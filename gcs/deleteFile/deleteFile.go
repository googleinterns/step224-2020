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
// DeleteFile is a probe that deletes a file from GCS.
// It tests the file deletion aspect of GCS and ensures that the file
// is deleted after the deletion request is made.
//
// TODOs:
// TODO: Define probe initialisation and run structure.
// TODO: Setup labels and metrics within probe
// TODO: Add handling of additional labels

package gcs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/glog"
	"github.com/google/cloudprober/logger"
	"github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes/options"
	"github.com/google/cloudprober/targets/endpoint"
	"google.golang.org/api/iterator"

	delete_pb "github.com/googleinterns/step224-2020/gcs/deleteFile/proto"
)

// Probe holds aggregate information about all probe runs, per-target.
type Probe struct {
	name    string
	config  *delete_pb.DeleteFileProbe
	targets []string
	opts    *options.Options

	probeOpLatency map[string]*metrics.EventMetrics           // Probe Operation Latency per target
	apiCallLatency map[string]map[string]*metric.EventMetrics // Latency per API call per target (map[target]map[apiCall]metric).
	logger         *logger.Logger

	// Add other required fields here - bucket
	// Array of Cancel funcs so that we can cancel externally, or just the context passed.
}

// Init initializes the probe with the given params.
func (p *Probe) Init(name string, opts *options.Options) error {
	c, ok := opts.ProbeConf.(*delete_pb.DeleteFileProbe)
	if !ok {
		return fmt.Errorf("not a gcs.deleteFile config")
	}
	p.config = c
	p.name = name
	p.opts = opts
	p.logger = opts.Logger
	p.probeOpLatency = make(map[string]*metrics.EventMetrics)
	p.apiCallLatency = make(map[string]map[string]*metrics.EventMetrics)

	return nil
}

// Start starts and runs the probe indefinitely.
func (p *Probe) Start(ctx context.Context, dataChan chan *metrics.EventMetrics) {
	probeTicker := time.NewTicker(p.opts.Interval)

	for {
		select {
		case <-ctx.Done():
			probeTicker.Stop()
			return
		case <-probeTicker.C:
			// On probe tick, write data to the channel and run probe.
			for _, metric := range p.probeOpLatency {
				dataChan <- metric
			}

			for _, apiCall := range p.apiCallLatency {
				for _, metric := range apiCall {
					dataChan <- metric
				}
			}

			// TODO: CHANGE THIS BIT TO OUR PROBE IMPLEMENTATION
			p.targets = endpoint.NamesFromEndpoints(p.opts.Targets.ListEndpoints())
			probeCtx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(p.opts.Timeout))
			p.initProbeMetrics()
			p.runProbe(probeCtx)
			cancelFunc()
		}
	}
}

// TODO: Update this section with the required metrics.
// TODO: Investigate how we go about adding additional distribution metrics.
// initProbeMetrics initializes missing probe metrics.
func (p *Probe) initProbeMetrics() {
	for _, target := range p.targets {
		if p.probeOpLatency[target] == nil {

			var probeLatDist *metrics.Distribution
			probeLatDistPb := p.config.HermesProbeDef.GetProbeLatencyDistribution()
			if probeLatDist, err = metrics.NewDistributionFromProto(probeLatDistPb); err != nil {
				p.logger.Errorf("error creating probe latency distribution from the specification (%v): %v", probeLatDistPb, err)
			}

			p.probeOpLatency[target] = metrics.NewEventMetrics(time.Now()).
				AddMetric("hermes_probe_operation_latency_s", probeLatDist).
				AddLabel("storage_system", p.config.GetTarmgetSystem().String()).
				AddLabel("probe_operation_type", p.config.GetProbeType().String()).
				AddLabel("target", target)
		}

		if p.apiCallLatency[target] == nil {
			p.apiCallLatency[target] = make(map[string]*metric.EventMetrics)

			var apiCallLatDist *metrics.Distribution
			apiCallLatDistPb := p.config.HermesProbeDef.GetApiCallLatencyDistribution()
			if probeLatDist, err = metrics.NewDistributionFromProto(apiCallLatDistPb); err != nil {
				p.logger.Errorf("error creating probe latency distribution from the specification (%v): %v", apiCallLatDistPb, err)
			}

			// List, Delete
			apiCallList := []string{"list", "delete"}
			for apiCall := range apiCallList {
				p.apiCallLatDist[target][apiCall] = metrics.NewEventMetrics(time.Now()).
					AddMetric("hermes_api_call_latency_ms", apiCallLatDist.Clone()).
					AddLabel("storage_system", p.config.GetTargetSystem().String()).
					AddLabel("api_call_type", apiCall).
					AddLabel("target", target)
			}
		}
	}
}

// TODO: Improve this by passing the storage client variable/making it accessible to reduce num of conns.
// TODO: Investigate getting credentials for multiple GCS instances.
// runProbeForTarget runs probe for a single target.
func (p *Probe) runProbeForTarget(ctx context.Context, target string) error {

	// Algorithm
	// List objects, pick one, delete it.
	// Check for errors (does GCS throw them?) -- YES.
	bucket := p.config.GetBucketName()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// First, list objects
	var names []string

	start := time.Now()
	it := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: ""})
	p.apiCallLatDist[target]["list"].AddFloat64(time.Now().Sub(start)) // What unit???

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			p.logger.Errorf("%v", err)
		}
		names = append(names, attrs.Name)
	}

	// Then, pick one to delete.
	object := names[0]

	// Then delete it
	o := client.Bucket(bucket).Object(object)
	start = time.Now()
	if err := o.Delete(ctx); err != nil {
		p.apiCallLatDist[target]["delete"].AddFloat64(time.Now().Sub(start))     // TODO(evanSpendlove) Add unit
		p.apiCallLatDist[target]["delete"].AddLabel("exit_status", err.String()) // TODO(evanSpendlove) Add a function that turns this into a proper label value
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	p.apiCallLatDist[target]["delete"].AddFloat64(time.Now().Sub(start))
	p.apiCallLatdist[target]["delete"].AddLabel("exit_status", "success")

	p.logger.Infof("Object %v deleted in bucket %s.", o, bucket)
	return nil
}

// runProbe runs probe for all targets and update EventMetrics.
func (p *Probe) runProbe(ctx context.Context) {
	p.targets = endpoint.NamesFromEndpoints(p.opts.Targets.ListEndpoints())

	var wg sync.WaitGroup
	for _, target := range p.targets {
		wg.Add(1)

		go func(target string, em *metrics.EventMetrics) {
			defer wg.Done()
			em.Metric("total").AddInt64(1)
			start := time.Now()
			err := p.runProbeForTarget(ctx, target) // run probe just for a single target
			if err != nil {
				p.logger.Errorf(err.Error())
				return
			}
			em.Metric("success").AddInt64(1)
			em.Metric("latency").AddFloat64(time.Now().Sub(start).Seconds() / p.opts.LatencyUnit.Seconds())
		}(target, p.results[target])

	}

	wg.Wait()
}
