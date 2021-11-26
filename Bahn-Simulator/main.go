// SPDX-License-Identifier: Apache-2.0
// Copyright 2021 Marcus Soll
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime/pprof"
	"sync"
)

func main() {
	inputPath := flag.String("input", "input.txt", "path to input file")
	outputPath := flag.String("output", "output.txt", "path to input file")
	profile := flag.String("pprof", "", "if set to a path, a pprof profile will be written")
	verbose := flag.Bool("verbose", false, "verbose output")
	flag.Parse()

	if *profile != "" {
		f, err := os.Create(*profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Panicln(err)
		}
		defer pprof.StopCPUProfile()
	}

	delay, successful := runSimulation(*inputPath, *outputPath, *verbose)
	if !successful {
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("Printing score")
	}
	fmt.Println(delay.String())
}

func runSimulation(input, output string, verbose bool) (*big.Int, bool) {
	world, err := ParseInput(input)
	if err != nil {
		fmt.Println("Can not read input file:", err)
		return big.NewInt(-1), false
	}

	if verbose {
		fmt.Println("Read output plans")
	}

	err = ParsePlan(world, output)
	if err != nil {
		fmt.Println("Can not read input file:", err)
		return big.NewInt(-1), false
	}

	if verbose {
		fmt.Println("Validating word begin")
	}
	errs := world.ValidateStart()
	if errs != nil {
		fmt.Println("initial validation failed:")
		for i := range errs {
			fmt.Println(errs[i].Error())
		}
		return big.NewInt(-1), false
	}

	// Run simulation
	for world.CurrentTime.Cmp(&world.MaxTime) != +1 {
		world.CurrentTime.Add(&world.CurrentTime, big.NewInt(1))
		if verbose {
			fmt.Println("Timestep", world.CurrentTime.String())
		}
		var errFound bool

		e := make(chan error, 1)
		var wg sync.WaitGroup

		// Trains
		for k := range world.Trains {
			wg.Add(1)
			go world.Trains[k].Update(world, e, &wg)
		}

		go func() {
			wg.Wait()
			close(e)
		}()

		for err := range e {
			errFound = true
			fmt.Println("trains -", world.CurrentTime.String(), "-", err.Error())
		}

		if errFound {
			return big.NewInt(-1), false
		}

		// Passengers
		e = make(chan error, 1)

		for k := range world.Passengers {
			wg.Add(1)
			go world.Passengers[k].Update(world, e, &wg)
		}

		go func() {
			wg.Wait()
			close(e)
		}()

		for err := range e {
			errFound = true
			fmt.Println("passengers -", world.CurrentTime.String(), "-", err.Error())
		}

		if errFound {
			return big.NewInt(-1), false
		}

		// Validate
		if verbose {
			fmt.Println("Validate", world.CurrentTime.String())
		}

		errs = world.Validate()
		if errs != nil {
			fmt.Println("validation", world.CurrentTime.String(), "failed:")
			for i := range errs {
				fmt.Println(errs[i].Error())
			}
			return big.NewInt(-1), false
		}
	}

	// Check result and calculate delay
	if verbose {
		fmt.Println("Calculating score")
	}

	delay := big.NewInt(0)
	valid := true

	for k := range world.Passengers {
		d := world.Passengers[k].Delay()
		if d.Cmp(InvalidDelay) == 0 {
			valid = false
			fmt.Println("passenger", k, "does not reach goal")
		}
		delay.Add(delay, d)
	}

	if !valid {
		return big.NewInt(-1), false
	}

	return delay, true
}
