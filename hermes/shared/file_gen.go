// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// File generation for Hermes
package shared

import (
	"crypto/sha1"
	"fmt"
)

type HermesFile struct {
	name     string
	contents string
}

// method of HermesFile generates HermesFile.name of the form Hermes_id_checksum where id is an integer & id <= 50
func (file *HermesFile) generateFileName(file_id int, file_checksum string) {
	file.name = fmt.Sprintf("Hermes_%02d_%v", file_id, file_checksum)
}

// method of HermesFile generates HermesFile.contents now a string without any significance in the future a pseudo random byte generator will be used
func (file *HermesFile) generateFileContents(file_id int) {
	file.contents = "jhfvjhdfjhfjjhjhdfvjvcvfjh"
}

// method of HermesFile generates the checksum of the file contents
func (file HermesFile) generateFileChecksum() string {
	file_contents := []byte(file.contents)
	hash := sha1.Sum(file_contents)
	// return checksum in hex notation
	return fmt.Sprintf("%x", hash)
}

// method of HermesFile generates the file takes id as a parameter
func GenerateHermesFile(id int) HermesFile {
	file := HermesFile{}
	file.generateFileContents(id)
	checksum := file.generateFileChecksum()
	file.generateFileName(id, checksum)
	return file
}
