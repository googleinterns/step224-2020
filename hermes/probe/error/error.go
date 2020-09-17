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
// Error implements the ProbeError used in Hermes for recording the
// exit status of functions.

// Package error defines the ProbeError used in Hermes for recording the
// exit status of functions.
package error

import (
	"errors"
	"fmt"

	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
)

// ProbeError is an error that includes an exit status for used within Hermes.
type ProbeError struct {
	// TODO(evanSpendlove): Refactor metrics.ExitStatus to be named something similar to metrics.APICallStatus.
	Status metrics.ExitStatus
	Err    error
}

// New returns a new ProbeError containing the error and status passed.
// Arguments:
//	- status: pass the exit status associated with this error.
//	- err: pass the error to be embedded.
// Returns:
//	- ProbeError: returns a new ProbeError object containing the args passed.
func New(status metrics.ExitStatus, err error) *ProbeError {
	return &ProbeError{
		Status: status,
		Err:    err,
	}
}

// Error returns the error string from the error.
// Returns:
//	- string: returns the error as a string.
func (e *ProbeError) Error() string {
	return fmt.Sprintf("%v: %v", e.Status, e.Err)
}
