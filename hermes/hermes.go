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
// hermes.go defines the structures necessary for the main hermes object to
// monitor a storage system.

package hermes

import (
	pb "github.com/googleinterns/step224-2020/hermes/proto"
)

// Hermes is the main prober struct. It contains the monitoring state and is used when monitoring target systems.
type Hermes struct {
	// Journal stores the state of Hermes as a combination of a next operation intent enum and a filenames map
	Journal *pb.StateJournal
}
