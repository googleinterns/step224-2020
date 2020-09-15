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
	"fmt"

	"github.com/googleinterns/step224-2020/hermes/probe/metrics"
)

// ProbeError is an error that includes an exit status for used within Hermes.
type ProbeError struct {
	Status metrics.ExitStatus
	Err    error
}

// Wrap returns a new ProbeError containing the error passed
// with additional context.and the status.
// Arguments:
//	- status: pass the exit status associated with this error.
//	- msg: pass the additional context of the error as a string.
//	- innner: pass the inner error to be wrapped.
// Returns:
//	- ProbeError: returns a new ProbeError object containing the args passed.
func Wrap(status metrics.ExitStatus, msg string, inner error) ProbeError {
	return &ProbeError{
		Status: status,
		Err:    fmt.Errorf("%s: %w", msg, inner),
	}
}

// Error() returns the error string from the error.
// Returns:
//	- string: returns the error as a string.
func (e *ProbeError) Error() string {
	return fmt.Errorf("%v: %s: %w", e.Status, e.Msg, e.Inner)
}

// Is returns true if the argument matches this object.
// Each error that is wrapped is examined to find a match.
// Arguments:
//	- target: pass the error to be compared with.
// Returns:
//	- bool: true if a match, else false.
func (e *ProbeError) Is(target error) bool {
	t, ok := target.(*ProbeError)
	if !ok {
		return false
	}
	return (e.Status == t.Status) && (errors.Is(e.Err, t.Err))
}
