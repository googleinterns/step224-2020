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
// Author: Alicja Kwiecinska (kwiecinskaa@google.com) github: alicjakwie
//
// File_gen provides file generation methods for Hermes.

// Package hermes is the general package for methods and function used within Hermes.
package hermes

import (
	"crypto/sha1"
	"fmt"
)

// HermesFile represents a file that Hermes will create on the storage system.
type HermesFile struct {
	name     string
	contents string
}

// generateFileName generates name of the form Hermes_ID_checksum.
// where ID is an integer & ID <= 49 & ID >= 0.
// Arguments:
//	- fileID: ID must be >= 0 and <= 49.
//	- fileChecksum: checksum must be a SHA1 checksum in hex notation.
func (f *HermesFile) generateFileName(fileID int, fileChecksum string) {
	f.name = fmt.Sprintf("Hermes_%02d_%v", fileID, fileChecksum)
}

// generateFileContents generates the contents of a file.
// TODO(alicjakwie): Use a pseudo-random number generator instead.
// Arguments:
//	- fileID: the ID of the file. Must be between 0 and 49.
func (f *HermesFile) generateFileContents(fileID int) {
	f.contents = "jhfvjhdfjhfjjhjhdfvjvcvfjh"
}

// generateFileChecksum generates the checksum of the file contents.
// Returns:
//	- checksum: Returns the checksum of the file contents in hex notation.
func (f *HermesFile) generateFileChecksum() string {
	fileContents := []byte(f.contents)
	hash := sha1.Sum(fileContents)
	// return checksum in hex notation
	return fmt.Sprintf("%x", hash)
}

// GenerateHermesFile generates a HermesFile object with a valid file name and contents.
// Arguments:
//	- id: id must be >= 0 and <= 49.
// Returns:
//	- HermesFile: returns a HermesFile object with a valid file name and contents.
func GenerateHermesFile(id int) HermesFile {
	file := HermesFile{}
	file.generateFileContents(id)
	file.generateFileName(id, file.generateFileChecksum())
	return file
}
