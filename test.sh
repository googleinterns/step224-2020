#!/bin/bash
# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Author: Evan Spendlove, GitHub: @evanSpendlove.
#
# This script runs all of tests in the project and compiles all of the proto
# files.

# Compile all protos
cd $HOME/go/src

protoc -I=. --go_out=. github.com/googleinterns/step224-2020/config/proto/*.proto

# Run go fmt on all .go files to format them
go fmt github.com/googleinterns/step224-2020/...

# Run all go tests
go test github.com/googleinterns/step224-2020/...
