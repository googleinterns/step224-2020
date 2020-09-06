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

import (
	"context"

	"github.com/google/cloudprober"
)

// FileOperation is an int used as part of the FileOperation enum within the intent log of Hermes' StateJournal.
type FileOperation int

// FileOperation enum is used for marking the operation intent within the StateJournal of Hermes.
// FileOperation has two possible values: CREATE, DELETE.
const (
	Create FileOperation = iota
	Delete
)

// String() allows the FileOperation constants to be conveniently converted to a print-friendly format.
// Returns:
// - string: Returns the print-friendly string version of the fileOperation enum.
func (fileOperation FileOperation) String() string {
	return [...]string{"Create", "Delete"}[fileOperation]
}

// Intent stores the next intended file operation of Hermes.
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
	intent    Intent         `json:"intent"`    // intent stores the next intended file operation and the name of the file that the operation is being performed on.
	filenames map[int]string `json:"filenames"` // filenames is a map of file IDs to filenames.
}

// Init initialises the map in the StateJournal so that entries can be added to it.
func (sj *StateJournal) Init() {
	sj.filenames = make(map[int]string) // initialise filenames to a map
}

// Hermes is the main Hermes prober that will startup Hermes and initiate monitoring targets.
type Hermes struct {
	Journal           StateJournal    // stateJournal stores the state of Hermes as a combination of next operation intent and a filenames map
	Ctx               context.Context // Context for starting Cloudprober
	CancelCloudprober func()          // CancelCloudprober is a cancel() function associated with the context passed to Cloudprober when initialised.
}

// InitialiseCloudproberFromConfig initialises Cloudprober from the config passed as an argument.
// Parameters:
// - config: config should be the contents of a Cloudprober config file. This is most likely: "grpc_port=9314"
//           -> the "grpc_port:" field is the only required field for the config.
// Returns:
// - error:
//	- logger.NewCloudproberLog() error: error initialising logging on GCE (Stackdriver)
//	- sysvars.Init():
//		- error getting local hostname: [error]:
//			-> error getting hostname from os.Hostname()
//		- other error
//			-> error initialising Cloud metadata
//	- config.ParseTemplate() error:
//		-> regex compilation issue of config or config could not be processed as a Go text template
//	- proto.UnmarshalText() error:
//		-> The config does not match the proto that it is being unmarshalled with.
//	- initDefaultServer() error:
//		- failed to parse default port from the env var: [serverEnvVar]=[parsedPort]
//		- error while creating listener for default HTTP server: [error]
//	- error while creating listener for default gRPC server: [error]
//	- tlsconfig.UpdateTLSConfig() error: an error occurred when updating the TLS config from the config passed.
func (h *Hermes) InitialiseCloudproberFromConfig(config string) error {
	return cloudprober.InitFromConfig(config)
}
