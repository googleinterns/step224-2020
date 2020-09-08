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
// TODO(evanSpendlove): Setup labels and metrics within probe
// TODO(evanSpendlove): Add handling of additional labels
// TODO(evanSpendlove): Add Start() method to Probe.
// TODO(evanSpendlove): Add initProbeMetrics() method to Probe.
// TODO(evanSpendlove): Add runProbeForTarget() method to Probe.
// TODO(evanSpendlove): Add runProbe() method to Probe.

// Package gcs implements the probes required to monitor a GCS instance.
package gcs

import (
	"fmt"

	"github.com/google/cloudprober/logger"
	"github.com/google/cloudprober/metrics"
	"github.com/google/cloudprober/probes/options"

	deletepb "github.com/googleinterns/step224-2020/gcs/deleteFile/proto"
)

// Probe holds aggregate information about all probe runs, per-target.
// It also holds the config and options used to initialise the probe.
type Probe struct {
	name    string
	config  *deletepb.DeleteFileProbe
	targets []string
	opts    *options.Options

	// probeOpLatency is used to record the latency of a probe operation
	// per target.
	probeOpLatency map[string]*metrics.EventMetrics
	// apiCallLatency is used to record latency values with distinct labels
	// per API call per target.
	// Usage: (map[target]map[apiCall][]metric).
	apiCallLatency map[string]map[string][]*metric.EventMetrics
	logger         *logger.Logger

	bucket string
	// CancelFunc cancels this probe and its goroutines.
	// If this probe needs to be cancelled immediately, call CancelFunc().
	CancelFunc context.CancelFunc
}

// Init initializes the probe with the given parameters.
// This is a required method to implement the Probe interface specified
// by Cloudprober.
// Arguments:
//	- name: probe name
//	- opts: probe configuration options built within Cloudprober
// Return:
//	- error: if an error is returned, the config is not a valid delete_file config.
func (p *Probe) Init(name string, opts *options.Options) error {
	conf, ok := opts.ProbeConf.(*deletepb.DeleteFileProbe) // Casting ProbeConf interface to DeleteFileProbe type
	if !ok {
		return fmt.Errorf("invalid argument: opts.ProbeConf is not of type *deletepb.DeleteFileProbe.")
	}
	p.config = conf
	p.name = name
	p.opts = opts
	p.logger = opts.Logger
	p.probeOpLatency = make(map[string]*metrics.EventMetrics)
	p.apiCallLatency = make(map[string]map[string]*metrics.EventMetrics)
	p.bucket = p.config.GetBucketName()

	return nil
}

// TODO(evanSpendlove): Add Start() method to Probe.
// TODO(evanSpendlove): Add initProbeMetrics() method to Probe.
// TODO(evanSpendlove): Add runProbeForTarget() method to Probe.
// TODO(evanSpendlove): Add runProbe() method to Probe.
