// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Evan Spendlove (@evanSpendlove)
//
// Client_test tests the client.go file to ensure that the gRPC client for connecting
// to Cloudprober works as expected. This client is used for adding and removing probes.

package cmd

import (
	"testing"
	"github.com/google/cloudprober/probes"
	"github.com/googleinterns/step224-2020/cmd"
	"github.com/googleinterns/step224-2020/hermes"
)

func setupHermes() {

}

func teardownHermes() {

}

func TestMain(m *testing.M) {
	setupHermes()
	code := m.Run()
	teardownHermes()
	os.Exit(code)
}

func TestRegisterAndAddProbeFromFilepath(t *testing.T) {

}

func TestRegisterAndAddProbe(t *testing.T) {

}

func TestRemoveProbe(t *testing.T) {

}

func TestListProbes(t *testing.T) {

}

