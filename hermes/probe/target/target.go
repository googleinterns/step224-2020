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

// Package target implements the Target struct used to hold all information
// and state for a single probe run.
package target

import (
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
