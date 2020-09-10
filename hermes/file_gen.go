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

// File generation for Hermes
package hermes

import (
	"crypto/sha1"
	"fmt"
	"time"
)

type HermesFile struct {
	name     string
	contents string
}

// generateFileName generates HermesFile.name of the form Hermes_id_checksum where id is an integer & id <= 50
func generateFileName(file_id int, file_checksum string) string {
	return fmt.Sprintf("Hermes_%02d_%v", file_id, file_checksum)
}

func generateFileContents(file_id int) string {
	return "jhfvjhdfjhfjjhjhdfvjvcvfjh"
}

// method of HermesFile generates the checksum of the file contents
func (file HermesFile) generateFileChecksum() string {
	file_contents := []byte(file.contents)
	hash := sha1.Sum(file_contents)
	// return checksum in hex notation
	return fmt.Sprintf("%x", hash)
}

// method of HermesFile generates the file takes id as a parameter
func NewHermesFile(id int) (*HermesFile, error) {
	if id < 0 || id > 50 {
		return nil, fmt.Errorf("At %v The file id provided wasn't in the required range [0,50]", time.Now())
	}
	file := HermesFile{}
	file.contents = generateFileContents(id)
	checksum := file.generateFileChecksum()
	file.name = generateFileName(id, checksum)
	return &file, nil
}