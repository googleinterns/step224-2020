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

func TestChecksum(t *testing.T) {
	want := "68f3caf439065824dcf75651c202e9f7c28ebf07" //expected checksum result
	file := HermesFile{}
	file.contents = "jhfvjhdfjhfjjhjhdfvjvcvfjh"
	got := file.generateFileChecksum()
	if want != got {
		t.Errorf("generateFileChecksum() failed expected %v got %v", want, got)
	}
}

type FileNameTestUnit struct {
	id       int
	checksum string
	want     string
}

var test_unit_table [5]FileNameTestUnit

func TestFileName(t *testing.T) {
	test_unit_table[0].id = 1
	test_unit_table[0].checksum = "abba"
	test_unit_table[0].want = fmt.Sprintf("Hermes_%02d_%v", test_unit_table[0].id, test_unit_table[0].checksum)
	test_unit_table[1].id = 11
	test_unit_table[1].checksum = "abba"
	test_unit_table[1].want = fmt.Sprintf("Hermes_%02d_%v", test_unit_table[1].id, test_unit_table[1].checksum)
	test_unit_table[2].id = 21
	test_unit_table[2].checksum = "abba"
	test_unit_table[2].want = fmt.Sprintf("Hermes_%02d_%v", test_unit_table[2].id, test_unit_table[2].checksum)
	test_unit_table[3].id = 31
	test_unit_table[3].checksum = "abba"
	test_unit_table[3].want = fmt.Sprintf("Hermes_%02d_%v", test_unit_table[3].id, test_unit_table[3].checksum)
	test_unit_table[4].id = 41
	test_unit_table[4].checksum = "abba"
	test_unit_table[4].want = fmt.Sprintf("Hermes_%02d_%v", test_unit_table[4].id, test_unit_table[4].checksum)
	for i := 0; i < 5; i++ {
		got := generateFileName(test_unit_table[i].id, test_unit_table[i].checksum)
		if got != test_unit_table[i].want {
			t.Errorf("generateFileName(%v, \"abba\") failed expected %v got %v", fmt.Sprintf("%2d", test_unit_table[i].id), test_unit_table[i].want, got)
		}
	}
}
