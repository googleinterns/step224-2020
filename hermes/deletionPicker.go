// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Alicja Kwiecinska (kwiecinskaa@google.com) github: alicjakwie
//
// DeletionPicker picks what file from 10 to 50 (files 0-9 are kept in the system for persistence) to delete
package hermes

import (
	"fmt"
	"math/rand"
	"time"
)

// PickFileToDelete picks which file to delete and returns a string: file name in the form "Hermes_ID".
func PickFileToDelete() string {
	rand.Seed(time.Now().UnixNano())
	const beg = 10; // we can delete files staring from the file Hermes_10
	const number_of_deletable_files = 41 // there are 41 files to delete from [Hermes_10,Hermes_50]
	file_id := rand.Intn(number_of_deletable_files) + beg // rand.Intn will return a natural number in the range [0, number_of_deletable_files) so file_id will be in the range [beg, number_of_deletable_files) 
	return fmt.Sprintf("Hermes_%02d", file_id)
}
