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
// hermes.go defines the structures necessary for the main hermes object to
// monitor a storage system.

package hermes

// FileOperation is used as part of the FileOperation enum within the intent log of Hermes' StateJournal.
type FileOperation int

// FileOperation is used for marking the operation intent within the StateJournal of Hermes.
// FileOperation has two possible values: CREATE, DELETE.
const (
	Create FileOperation = iota
	Delete
)

// String() allows the FileOperation constants to be conveniently converted to a print-friendly format.
// Returns:
// - string: Returns the print-friendly string version of the fileOperation enum.
func (fileOperation FileOperation) String() string {
	var FileOperationName = map[FileOperation]string{
		Create: "Create",
		Delete: "Delete",
	}

	return FileOperationName[fileOperation]
}

// Intent stores the next intended file operation.
// This is used as part of the StateJournal of Hermes.
type Intent struct {
	operation FileOperation `json:"fileoperation"` // Stores the file operation intent, either CREATE or DELETE
	filename  string        `json:"filename"`      // Stores the filename that the operation is being performed on.
}

// StateJournal stores the state of Hermes in two parts: the next operation intent and a map of filenames.
// The intent stores the operation to be performed and the name of the file that the operation is being performed on.
// The filenames map is a map of file IDs to filenames.
// If an entry does not exist for a given ID, then the file does not exist (i.e. has been deleted)
type StateJournal struct {
	// intent stores the next intended file operation and the name of the file that the operation is being performed on.
	intent Intent `json:"intent"`
	// filenames is a map of file IDs to filenames.
	filenames map[int]string `json:"filenames"`
}

func (sj *StateJournal) Init() {
	sj.filenames = make(map[int]string)
}

// Hermes is the main Hermes prober that will startup Hermes and initiate monitoring targets.
type Hermes struct {
	// Journal stores the state of Hermes as a combination of next operation intent and a filenames map
	Journal StateJournal
}
