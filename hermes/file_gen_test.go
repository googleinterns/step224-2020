// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: Alicja Kwiecinska (kwiecinskaa@google.com) github: alicjakwie

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

var test_unit_table [5]FileNameTestUnit

// Function completes the table of HermesFileTestUnit by iterating over its entries
func completeTestUnitTable() {
	for i := 0; i < 5; i++ {
		test_unit_table[i].id = i*10 + 1
		test_unit_table[i].checksum = "abba"
		test_unit_table[i].want = fmt.Sprintf("Hermes_%02d_%v", test_unit_table[i].id, test_unit_table[i].checksum)
	}
}

// Test function for the generateFileChecksum function in file_gen.go
func TestChecksum(t *testing.T) {
	want := "68f3caf439065824dcf75651c202e9f7c28ebf07" //expected checksum result
	file := HermesFile{}
	file.contents = "jhfvjhdfjhfjjhjhdfvjvcvfjh"
	got := file.generateFileChecksum()
	if want != got {
		t.Errorf("generateFileChecksum() failed expected %v got %v", want, got)
	}
}

// Test function for the generateFileName function in file_gen.go
func TestFileName(t *testing.T) {
	completeTestUnitTable()
	for i := 0; i < 5; i++ {
		file := HermesFile{}
		file.generateFileName(test_unit_table[i].id, test_unit_table[i].checksum)
		got := file.name
		if got != test_unit_table[i].want {
			t.Errorf("generateFileName(%v, \"abba\") failed expected %v got %v", fmt.Sprintf("%2d", test_unit_table[i].id), test_unit_table[i].want, got)
		}
	}
}
