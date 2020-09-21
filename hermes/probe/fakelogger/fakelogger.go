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
// Authors: Alicja Kwiecinska GitHub: alicjakwie
//
// package fakelogger contains all of the logic necessary to create a fake instance of cloudprober's logger.
package fakelogger

import (
	"context"
	"fmt"

	"github.com/google/cloudprober/logger"
)

type FakeLogger struct {
	Logger *logger.Logger
}

// NewLogger returns a new Logger object
// Argument:
//      ctx: it carries deadlines and cancellation signals that originate from read_test.go and create_test.go
// Returns:
//		a fake instance of cloudprober's logger
func NewLogger(ctx context.Context) *FakeLogger {
	l, err := logger.New(ctx, "FakeLogger")
	if err != nil {
		fmt.Errorf("NewLogger failed to create a new instance of cloudprober logger.Loggger: %v", err)
	}
	return &FakeLogger{
		Logger: l,
	}
}
