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
	"regexp"
	"sync"
)

type TrainPosition int

const (
	TrainPositionUnknown TrainPosition = iota
	TrainPositionStation
	TrainPositionLine
	TrainPositionWildcard
)

type Train struct {
	ID               string
	Capacity         big.Int
	Passengers       big.Int
	Speed            big.Rat
	Position         []string
	PositionSince    big.Rat
	PositionType     TrainPosition
	Plan             map[string]string
	BoardingPossible bool
	L                sync.Mutex
}

var trainPlanRegexp = regexp.MustCompile(`\A(?P<time>[\d]+) (?P<action>(Start)|(Depart)) (?P<id>[a-zA-Z0-9_]+)[\s]*\z`)
var trainPlanRegexpTime = trainPlanRegexp.SubexpIndex("time")
var trainPlanRegexpAction = trainPlanRegexp.SubexpIndex("action")
var trainPlanRegexpID = trainPlanRegexp.SubexpIndex("id")

func (t *Train) IsValidStart(w *World) error {
	if t.Capacity.Sign() == -1 {
		return fmt.Errorf("train (%s): capacity must not be negative", t.ID)
	}

	if t.Speed.Cmp(big.NewRat(0, 1)) != +1 {
		return fmt.Errorf("train (%s): speed must be larger than 0", t.ID)
	}
	return t.IsValid(w)
}

func (t *Train) IsValid(w *World) error {
	t.L.Lock()
	defer t.L.Unlock()

	if t.PositionType == TrainPositionUnknown {
		return fmt.Errorf("train (%s): unknown position", t.ID)
	}

	if t.PositionType == TrainPositionStation {
		if len(t.Position) != 1 {
			return fmt.Errorf("[internal] train (%s): wrong number of positions", t.ID)
		}
		_, ok := w.Stations[t.Position[0]]
		if !ok {
			return fmt.Errorf("train (%s): unknown station '%s'", t.ID, t.Position[0])
		}
	}

	if t.PositionType == TrainPositionLine {
		if len(t.Position) != 2 {
			return fmt.Errorf("[internal] train (%s): wrong number of positions", t.ID)
		}
		_, ok := w.Lines[t.Position[0]]
		if !ok {
			return fmt.Errorf("train (%s): unknown line '%s'", t.ID, t.Position[0])
		}
		_, ok = w.Stations[t.Position[1]]
		if !ok {
			return fmt.Errorf("train (%s): unknown target station '%s'", t.ID, t.Position[0])
		}
	}

	if t.Passengers.Sign() == -1 {
		return fmt.Errorf("train (%s): passengers must not be negative", t.ID)
	}

	if t.Capacity.Cmp(&t.Passengers) == -1 {
		return fmt.Errorf("train (%s): too many passengers (capacity: %s, current: %s)", t.ID, t.Capacity.String(), t.Passengers.String())
	}
	return nil
}

func (t *Train) Update(w *World, e chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	t.L.Lock()
	defer t.L.Unlock()

	// Update position
	switch t.PositionType {
	case TrainPositionStation:
		t.BoardingPossible = true
	case TrainPositionLine:
		t.BoardingPossible = false
		err := t.advanceLinePosition(w)
		if err != nil {
			e <- err
			return
		}
	case TrainPositionWildcard:
		t.BoardingPossible = false
	default:
		e <- fmt.Errorf("train (%s): unknown position", t.ID)
		return
	}

	plan, ok := t.Plan[w.CurrentTime.String()]

	if !ok {
		// Nothing to do here
		return
	}

	err := t.processRule(w, plan)
	if err != nil {
		e <- err
	}
}

func (t *Train) advanceLinePosition(w *World) error {
	if t.PositionType != TrainPositionLine {
		return fmt.Errorf("train %s (internal): positionType must be %d  but is %d", t.ID, TrainPositionLine, t.PositionType)
	}
	t.PositionSince.Add(&t.PositionSince, big.NewRat(1, 1))
	distance := big.NewRat(1, 1).Mul(&t.PositionSince, &t.Speed)
	if len(t.Position) != 2 {
		return fmt.Errorf("train %s (internal): position %v can not be right (length must be 2)", t.ID, t.Position)
	}
	line, ok := w.Lines[t.Position[0]]
	if !ok {
		return fmt.Errorf("train %s: line %s does not exist", t.ID, t.Position[0])
	}

	if distance.Cmp(&line.Length) >= 0 {
		// Reached end of line
		st, ok := w.Stations[t.Position[1]]
		if !ok {
			return fmt.Errorf("train %s: reached non existing station %s", t.ID, t.Position[1])
		}
		t.Position = []string{t.Position[1]}
		t.PositionType = TrainPositionStation
		line.L.Lock()
		line.CurrentCapacity.Sub(&line.CurrentCapacity, big.NewInt(1))
		line.L.Unlock()
		st.L.Lock()
		st.CurrenTrains.Add(&st.CurrenTrains, big.NewInt(1))
		st.L.Unlock()
	}

	return nil
}

func (t *Train) processRule(w *World, rule string) error {
	if t.PositionType == TrainPositionLine {
		return fmt.Errorf("train (%s): new plan but train is still on line", t.ID)
	}
	if t.PositionType == TrainPositionWildcard {
		return fmt.Errorf("train (%s): new plan but train is still on '*' (no Start rule)", t.ID)
	}
	matches := trainPlanRegexp.FindStringSubmatch(rule)
	if matches == nil {
		return fmt.Errorf("train (%s): can not match rule '%s'", t.ID, rule)
	}
	switch matches[trainPlanRegexpAction] {
	case "Start":
		return fmt.Errorf("train (%s): at this point (%s) Start is not allowed", t.ID, w.CurrentTime.String())
	case "Depart":
		lineID := matches[trainPlanRegexpID]
		line, ok := w.Lines[lineID]
		if !ok {
			return fmt.Errorf("train (%s): unknown target line %s", t.ID, lineID)
		}
		currentPosition := t.Position[0]
		t.PositionType = TrainPositionLine
		t.Position = []string{line.ID, ""}
		t.BoardingPossible = false
		var foundPosition bool
		for i := range line.End {
			if line.End[i] == currentPosition {
				foundPosition = true
			} else {
				t.Position[1] = line.End[i]
			}
		}
		if t.Position[1] == "" {
			return fmt.Errorf("train (%s): target line %s does not have a target station", t.ID, lineID)
		}
		if !foundPosition {
			return fmt.Errorf("train (%s): target line %s does not connect to current station %s", t.ID, lineID, currentPosition)
		}
		line.L.Lock()
		line.CurrentCapacity.Add(&line.CurrentCapacity, big.NewInt(1))
		line.L.Unlock()
		st, ok := w.Stations[currentPosition]
		if !ok {
			return fmt.Errorf("train %s: depature from non existing station %s", t.ID, currentPosition)
		}
		st.L.Lock()
		st.CurrenTrains.Sub(&st.CurrenTrains, big.NewInt(1))
		st.L.Unlock()
		t.PositionSince = *big.NewRat(0, 1)
		t.advanceLinePosition(w)
	default:
		return fmt.Errorf("train (%s): unknown action '%s'", t.ID, matches[trainPlanRegexpAction])
	}
	return nil
}
