// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Marcus Soll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestMain(t *testing.T) {
	dirs, err := os.ReadDir("test/")
	if err != nil {
		fmt.Println("can not read test dir:", err.Error())
		t.FailNow()
	}
	for i := range dirs {
		if dirs[i].IsDir() {
			input := path.Join("test", dirs[i].Name(), "input.txt")
			output := path.Join("test", dirs[i].Name(), "output.txt")
			_, successful := runSimulation(input, output, false)
			if !successful {
				fmt.Println("Test", dirs[i].Name(), "failed")
				t.Fail()
			}
		}
	}
}
