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
// Test suite for file_gen.go

package hermes

import (
	"fmt"
	"testing"
)

// Hermes file test unit for generateFileName
type FileNameTestUnit struct {
	id       int
	checksum string
	want     string
}

var testUnitTable [5]FileNameTestUnit

// completeTestUnitTable completes the table of HermesFileTestUnit by iterating over its entries
func completeTestUnitTable() {
	for i := 0; i < 5; i++ {
		testUnitTable[i].id = i*10 + 1
		testUnitTable[i].checksum = "abba"
		testUnitTable[i].want = fmt.Sprintf("Hermes_%02d_%v", testUnitTable[i].id, testUnitTable[i].checksum)
	}
}

// TestChecksum tests the generateFileChecksum() method.
func TestChecksum(t *testing.T) {
	want := "68f3caf439065824dcf75651c202e9f7c28ebf07" //expected checksum result
	file := HermesFile{}
	file.contents = "jhfvjhdfjhfjjhjhdfvjvcvfjh"
	got := file.generateFileChecksum()
	if want != got {
		t.Errorf("generateFileChecksum() failed expected %v got %v", want, got)
	}
}

// TestFileName tests the generateFileName method.
func TestFileName(t *testing.T) {
	completeTestUnitTable()
	for i := 0; i < 5; i++ {
		file := HermesFile{}
		file.generateFileName(testUnitTable[i].id, testUnitTable[i].checksum)
		got := file.name
		if got != testUnitTable[i].want {
			t.Errorf("generateFileName(%v, \"abba\") failed expected %v got %v", fmt.Sprintf("%2d", testUnitTable[i].id), testUnitTable[i].want, got)
		}
	}
}
