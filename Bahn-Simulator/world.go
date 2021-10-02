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
	"math/big"
	"sync"
)

type World struct {
	Lines       map[string]*Line
	Stations    map[string]*Station
	Trains      map[string]*Train
	Passengers  map[string]*Passenger
	CurrentTime big.Int
	MaxTime     big.Int
}

type Line struct {
	ID              string
	L               sync.Mutex
	End             []string
	Length          big.Rat
	MaxCapacity     big.Int
	CurrentCapacity big.Int
}

type Station struct {
	ID           string
	Capacity     big.Int
	CurrenTrains big.Int
	L            sync.Mutex
}

func (l *Line) IsValidStart(w *World) error {
	if len(l.End) != 2 {
		return fmt.Errorf("line (%s): ends do not fit (must: 2, is: %d)", l.ID, len(l.End))
	}

	_, ok := w.Stations[l.End[0]]
	if !ok {
		return fmt.Errorf("line (%s): unknown station '%s'", l.ID, l.End[0])
	}

	_, ok = w.Stations[l.End[1]]
	if !ok {
		return fmt.Errorf("line (%s): unknown station '%s'", l.ID, l.End[1])
	}

	if l.Length.Cmp(big.NewRat(0, 1)) != +1 {
		return fmt.Errorf("line (%s): length '%s' must be larger than 0", l.ID, l.Length.String())
	}

	if l.MaxCapacity.Cmp(big.NewInt(0)) != +1 {
		return fmt.Errorf("line (%s): maximum capacity '%s' must be larger than 0", l.ID, l.Length.String())
	}
	return l.IsValid(w)
}

func (l *Line) IsValid(w *World) error {
	l.L.Lock()
	defer l.L.Unlock()

	if l.MaxCapacity.Cmp(&l.CurrentCapacity) == -1 {
		return fmt.Errorf("line (%s): too many trains (capacity: %s, current: %s)", l.ID, l.MaxCapacity.String(), l.CurrentCapacity.String())
	}

	return nil
}

func (s *Station) IsValidStart(w *World) error {
	if s.Capacity.Cmp(big.NewInt(0)) != +1 {
		return fmt.Errorf("station (%s): maximum capacity '%s' must be larger than 0", s.ID, s.Capacity.String())
	}
	return s.IsValid(w)
}

func (s *Station) IsValid(w *World) error {
	if s.Capacity.Cmp(&s.CurrenTrains) == -1 {
		return fmt.Errorf("station (%s): too many trains (capacity: %s, current: %s)", s.ID, s.Capacity.String(), s.CurrenTrains.String())
	}

	return nil
}

func (w *World) Validate() []error {
	e := make(chan error)
	var errs []error
	var wg sync.WaitGroup
	for k := range w.Stations {
		wg.Add(1)
		go func(index string) {
			defer wg.Done()
			err := w.Stations[index].IsValid(w)
			if err != nil {
				e <- fmt.Errorf("validation failed for station '%s': %s", index, err.Error())
			}
		}(k)
	}

	for k := range w.Lines {
		wg.Add(1)
		go func(index string) {
			defer wg.Done()
			err := w.Lines[index].IsValid(w)
			if err != nil {
				e <- fmt.Errorf("validation failed for line '%s': %s", index, err.Error())
			}
		}(k)
	}

	for k := range w.Trains {
		wg.Add(1)
		go func(index string) {
			defer wg.Done()
			err := w.Trains[index].IsValid(w)
			if err != nil {
				e <- fmt.Errorf("validation failed for train '%s': %s", index, err.Error())
			}
		}(k)
	}

	for k := range w.Passengers {
		wg.Add(1)
		go func(index string) {
			defer wg.Done()
			err := w.Passengers[index].IsValid(w)
			if err != nil {
				e <- fmt.Errorf("validation failed for passenger '%s': %s", index, err.Error())
			}
		}(k)
	}

	go func() {
		wg.Wait()
		close(e)
	}()

	for err := range e {
		errs = append(errs, err)
	}

	return errs
}

func (w *World) ValidateStart() []error {
	var errs []error

	e := make(chan error)
	var wg sync.WaitGroup
	for k := range w.Stations {
		wg.Add(1)
		go func(index string) {
			defer wg.Done()
			err := w.Stations[index].IsValidStart(w)
			if err != nil {
				e <- fmt.Errorf("validation failed for station '%s': %s", index, err.Error())
			}
		}(k)
	}

	for k := range w.Lines {
		wg.Add(1)
		go func(index string) {
			defer wg.Done()
			err := w.Lines[index].IsValidStart(w)
			if err != nil {
				e <- fmt.Errorf("validation failed for line '%s': %s", index, err.Error())
			}
		}(k)
	}

	for k := range w.Trains {
		wg.Add(1)
		go func(index string) {
			defer wg.Done()
			err := w.Trains[index].IsValidStart(w)
			if err != nil {
				e <- fmt.Errorf("validation failed for train '%s': %s", index, err.Error())
			}
		}(k)
	}

	for k := range w.Passengers {
		wg.Add(1)
		go func(index string) {
			defer wg.Done()
			err := w.Passengers[index].IsValidStart(w)
			if err != nil {
				e <- fmt.Errorf("validation failed for passenger '%s': %s", index, err.Error())
			}
		}(k)
	}

	go func() {
		wg.Wait()
		close(e)
	}()

	for err := range e {
		errs = append(errs, err)
	}

	if !w.CheckConnected() {
		errs = append(errs, fmt.Errorf("validation failed for world: not all stations are connected"))
	}

	return errs
}

func (w *World) CheckConnected() bool {
	if len(w.Stations) == 0 {
		return true
	}
	if len(w.Lines) == 0 {
		return false
	}
	marked := make(map[string]bool)
	bfs := make([]string, 0, len(w.Stations))
	first := true

	for k := range w.Stations {
		if first {
			bfs = append(bfs, k)
			marked[k] = true
			first = false
		} else {
			marked[k] = false
		}
	}

	for len(bfs) != 0 {
		current := bfs[0]
		bfs = bfs[1:]

		for k := range w.Lines {
			if w.Lines[k].End[0] == current && !marked[w.Lines[k].End[1]] {
				bfs = append(bfs, w.Lines[k].End[1])
				marked[w.Lines[k].End[1]] = true
			}
			if w.Lines[k].End[1] == current && !marked[w.Lines[k].End[0]] {
				bfs = append(bfs, w.Lines[k].End[0])
				marked[w.Lines[k].End[0]] = true
			}
		}
	}

	for k := range marked {
		if !marked[k] {
			return false
		}
	}

	return true
}
