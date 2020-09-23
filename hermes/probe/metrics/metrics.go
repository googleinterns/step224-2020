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
// Author: Evan Spendlove, GitHub: evanSpendlove.
//
// Metrics implements the latency distribution metrics for Hermes.

// Package metrics implements the metrics used for recording information
// when monitoring a storage system.
package metrics

import (
	"fmt"
	"time"

	"github.com/google/cloudprober/metrics"

	probepb "github.com/googleinterns/step224-2020/config/proto"
)

// ProbeOperation represents a possible probe operation metric label.
type ProbeOperation int

const (
	// TotalProbeRun is the metric label for the overall probe run.
	TotalProbeRun ProbeOperation = iota
	// CheckNil is the metric label for the check nil file existence and consistency operation.
	CheckNil
	// ReadFile is the metric label for the read file operation.
	ReadFile
	// VerifyFileContents is the metric label for the verify file contents operation.
	VerifyFileContents
	// DeleteFile is the metric label for the delete file operation.
	DeleteFile
	// CreateFile is the metric label for the create file operation.
	CreateFile
)

// APICall represents a possible API call metric label.
type APICall int

const (
	// APIListFiles is the metric label for the list files API call.
	APIListFiles APICall = iota
	// APICreateFile is the metric label for the create file API call.
	APICreateFile
	// APIDeleteFile is the metric label for the delete file API call.
	APIDeleteFile
	// APIGetFile is the metric label for the get file API call.
	APIGetFile
)

// ExitStatus represents a possible exit status metric label.
type ExitStatus int

const (
	// Success indicates the operation/API call was successful.
	Success ExitStatus = iota
	// OpTimeout indicates that the operation/API call timed out.
	OpTimeout
	// ProbeFailed indicates that the probe failed due to an error.
	ProbeFailed
	// APICallFailed indicates that the API call failed due to an error.
	APICallFailed
	// FileMissing indicates that the target file could not be found.
	FileMissing
	// BucketMissing indicates that the target bucket could not be found.
	BucketMissing
	// FileCorrupted indicates that the target file was corrupted and the contents could not be read.
	FileCorrupted
	// FileReadFailure indicates that the target file could not be read.
	FileReadFailure
	// FileMetadataMismatch indicates that the target file did not match the expectation from its metadata
	// TODO(evanSpendlove): Review if we need this? When would we use it?
	FileMetadataMismatch
	// UnknownFileFound indicates that an unknown file, not created by Hermes, was found in the target bucket.
	UnknownFileFound
	// AllFilesMissing indicates that all of the Hermes files were missing.
	AllFilesMissing
)

var (
	// ProbeOpName maps ProbeOperation constants to their metric label string equivalent.
	ProbeOpName = map[ProbeOperation]string{
		TotalProbeRun:      "total_probe_run",
		CheckNil:           "check_nil",
		ReadFile:           "read_file",
		VerifyFileContents: "verify_file_contents",
		DeleteFile:         "delete_file",
		CreateFile:         "create_file",
	}
	// APICallName maps ApiCall constants to their metric label string equivalent.
	APICallName = map[APICall]string{
		APIListFiles:  "list_files",
		APICreateFile: "create_file",
		APIDeleteFile: "delete_file",
		APIGetFile:    "get_file",
	}
	// ExitStatusName maps ExitStatus constants to their metric label string equivalent.
	ExitStatusName = map[ExitStatus]string{
		Success:              "success",
		OpTimeout:            "op_timeout",
		ProbeFailed:          "probe_failed",
		APICallFailed:        "api_call_failed",
		FileMissing:          "file_missing",
		BucketMissing:        "bucket_missing",
		FileCorrupted:        "file_corrupted",
		FileReadFailure:      "file_read_failure",
		FileMetadataMismatch: "file_metadata_mismatch",
		UnknownFileFound:     "unknown_file_found",
		AllFilesMissing:      "all_files_missing",
	}
)

// Metrics stores the cumulative metrics for probe runs for a target.
type Metrics struct {
	// probeOpLatency is used to record latency values with distinct labels
	// per exit status per probe operation per target.
	// Recommended usage: probeOpLatency[ProbeOperation][ExitStatus].Metric("latency").AddFloat64(<val>)
	ProbeOpLatency map[ProbeOperation]map[ExitStatus]*metrics.EventMetrics
	// apiCallLatency is used to record latency values with distinct labels
	// per exit status per API call per target.
	// Recommended usage: apiCallLatency[ApiCall][ExitStatus].Metric("latency").AddFloat64(<val>)
	APICallLatency map[APICall]map[ExitStatus]*metrics.EventMetrics
}

// NewMetrics creates a new *Metrics object and initialises the fields inside it.
// Arguments:
//	- conf: pass a HermesProbeDef config
//	- target: pass the target for which metrics are to be collected.
// Returns:
//	- m: returns an initialised *Metrics object.
//	- err: returns an error if a latency distribution cannot be created from the config proto.
func NewMetrics(conf *probepb.HermesProbeDef, target *probepb.Target) (*Metrics, error) {
	m := &Metrics{
		ProbeOpLatency: make(map[ProbeOperation]map[ExitStatus]*metrics.EventMetrics, len(ProbeOpName)),
		APICallLatency: make(map[APICall]map[ExitStatus]*metrics.EventMetrics, len(APICallName)),
	}

	probeOpLatDist, err := metrics.NewDistributionFromProto(conf.GetProbeLatencyDistribution())
	if err != nil {
		return nil, fmt.Errorf("invalid argument: error creating probe latency distribution from the specification (%v): %w", conf.GetProbeLatencyDistribution(), err)
	}

	for op := range ProbeOpName {
		m.ProbeOpLatency[op] = make(map[ExitStatus]*metrics.EventMetrics, len(ExitStatusName))
		for e := range ExitStatusName {
			m.ProbeOpLatency[op][e] = metrics.NewEventMetrics(time.Now()).
				AddMetric("hermes_probe_latency_seconds", probeOpLatDist.Clone()).
				AddLabel("storage_system", target.GetTargetSystem().String()).
				AddLabel("target", fmt.Sprintf("%s:%s", target.GetName(), target.GetBucketName())).
				AddLabel("probe_operation_type", ProbeOpName[op]).
				AddLabel("exit_status", ExitStatusName[e])
		}
	}

	apiCallLatDist, err := metrics.NewDistributionFromProto(conf.GetApiCallLatencyDistribution())
	if err != nil {
		return nil, fmt.Errorf("invalid argument: error creating probe latency distribution from the specification (%v): %v", conf.GetApiCallLatencyDistribution(), err)
	}

	for call := range APICallName {
		m.APICallLatency[call] = make(map[ExitStatus]*metrics.EventMetrics, len(ExitStatusName))
		for e := range ExitStatusName {
			m.APICallLatency[call][e] = metrics.NewEventMetrics(time.Now()).
				AddMetric("hermes_api_latency_seconds", apiCallLatDist.Clone()).
				AddLabel("storage_system", target.GetTargetSystem().String()).
				AddLabel("target", fmt.Sprintf("%s:%s", target.GetName(), target.GetBucketName())).
				AddLabel("probe_operation_type", APICallName[call]).
				AddLabel("exit_status", ExitStatusName[e])
		}
	}

	return m, nil
}
