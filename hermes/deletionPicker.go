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

// Picks what file from 10 to 50 (files 0-9 are kept in the system for persistence) to delete
package hermes

import (
	"fmt"
	"math/rand"
	"time"
)

// PickFileToDelete picks which file to delete and returns a string: file name in the form "Hermes_ID".
func PickFileToDelete() string {
	rand.Seed(time.Now().UnixNano())
	file_id := rand.Intn(40) + 10;
	return fmt.Sprintf("Hermes_%02d", file_id)
}
