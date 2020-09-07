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

package gcsDelete

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/cloudprober/logger"
	"github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes/options"
	"github.com/google/cloudprober/targets/endpoint"

	"cloud.google.com/go/storage"
)

// Probe holds aggregate information about all probe runs, per-target.
type Probe struct {
	name    string
	config  *ProbeConf
	targets []string
	opts    *options.Options

	res map[string]*metrics.EventMetrics // Results by target
	l   *logger.Logger
}

// Init initializes the probe with the given params.
func (p *Probe) Init(name string, opts *options.Options) error {
	c, ok := opts.ProbeConf.(*ProbeConf)
	if !ok {
		return fmt.Errorf("not a gcs.deleteFile config")
	}
	p.config = c
	p.name = name
	p.opts = opts
	p.l = opts.Logger
	p.res = make(map[string]*metrics.EventMetrics)

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
			for _, metric := range p.res {
				dataChan <- metric
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
		if p.res[target] != nil {
			continue
		}
		var latVal metrics.Value
		if p.opts.LatencyDist != nil {
			latVal = p.opts.LatencyDist.Clone()
		} else {
			latVal = metrics.NewFloat(0)
		}
		p.res[target] = metrics.NewEventMetrics(time.Now()).
			AddMetric("total", metrics.NewInt(0)).
			AddMetric("success", metrics.NewInt(0)).
			AddMetric("latency", latVal).
			AddLabel("ptype", "redis").
			AddLabel("probe", p.name).
			AddLabel("dst", target)
	}
}

// TODO: Improve this by passing the storage client variable/making it accessible to reduce num of conns.
// TODO: Investigate getting credentials for multiple GCS instances.
// runProbeForTarget runs probe for a single target.
func (p *Probe) runProbeForTarget(ctx context.Context, target string) error {

	// Algorithm
	// List objects, pick one, delete it.
	// Check for errors (does GCS throw them?) -- YES.
	bucket := p.c.GetBucket()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// First, list objects
	var names []string
	iterator := bkt.Objects(ctx, &storage.Query{Prefix: ""})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			p.logger.Error(err)
		}
		names = append(names, attrs.Name)
	}

	// Then, pick one to delete.

	o := client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}

	glog.Infof("Object %v deleted in bucket %s.", o, bucket)
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
				p.l.Errorf(err.Error())
				return
			}
			em.Metric("success").AddInt64(1)
			em.Metric("latency").AddFloat64(time.Now().Sub(start).Seconds() / p.opts.LatencyUnit.Seconds())
		}(target, p.res[target])

	}

	wg.Wait()
}
