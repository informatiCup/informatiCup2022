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

type PassengerPosition int

var InvalidDelay = big.NewInt(-1)

const (
	PassengerPositionUnknown PassengerPosition = iota
	PassengerPositionStation
	PassengerPositionTrain
)

type Passenger struct {
	ID            string
	Start         string
	Target        string
	Size          big.Int
	TargetTime    big.Int
	TargetReached big.Int
	PositionType  PassengerPosition
	Position      string
	Plan          map[string]string
}

var passengerPlanRegexp = regexp.MustCompile(`\A(?P<time>[\d]+) (?P<action>(Board)|(Detrain)) ?(?P<id>[a-zA-Z0-9_]+)?[\s]*\z`)
var passengerPlanRegexpTime = trainPlanRegexp.SubexpIndex("time")
var passengerPlanRegexpAction = trainPlanRegexp.SubexpIndex("action")
var passengerPlanRegexpID = trainPlanRegexp.SubexpIndex("id")

func (p *Passenger) Delay() *big.Int {
	if p.TargetReached.Cmp(big.NewInt(0)) == 0 {
		return InvalidDelay
	}
	delay := big.NewInt(0)
	delay.Sub(&p.TargetReached, &p.TargetTime)
	if delay.Cmp(big.NewInt(0)) == -1 {
		return big.NewInt(0)
	}
	delay.Mul(delay, &p.Size)
	return delay
}

func (p *Passenger) IsValidStart(w *World) error {
	_, ok := w.Stations[p.Start]
	if !ok {
		return fmt.Errorf("passenger (%s): unknown start '%s'", p.ID, p.Start)
	}
	_, ok = w.Stations[p.Target]
	if !ok {
		return fmt.Errorf("passenger (%s): unknown start '%s'", p.ID, p.Target)
	}
	if p.Size.Cmp(big.NewInt(0)) != +1 {
		return fmt.Errorf("passenger (%s): size '%s' must be positive", p.ID, p.Size.String())
	}

	if p.TargetTime.Cmp(big.NewInt(0)) != +1 {
		return fmt.Errorf("passenger (%s): target time '%s' must be positive", p.ID, p.Size.String())
	}
	return p.IsValid(w)
}

func (p *Passenger) IsValid(w *World) error {
	if p.PositionType == PassengerPositionUnknown {
		return fmt.Errorf("passenger (%s): unknown position", p.ID)
	}
	if p.PositionType == PassengerPositionTrain {
		_, ok := w.Trains[p.Position]
		if !ok {
			return fmt.Errorf("passenger (%s): unknown train '%s'", p.ID, p.Position)
		}
	}
	if p.PositionType == PassengerPositionStation {
		_, ok := w.Stations[p.Position]
		if !ok {
			return fmt.Errorf("passenger (%s): unknown station '%s'", p.ID, p.Position)
		}
	}
	return nil
}

func (p *Passenger) Update(w *World, e chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	plan, ok := p.Plan[w.CurrentTime.String()]

	if !ok {
		// Nothing to do here
		return
	}

	p.TargetReached = big.Int{}
	matches := passengerPlanRegexp.FindStringSubmatch(plan)
	if matches == nil {
		e <- fmt.Errorf("passenger (%s): can not match rule '%s'", p.ID, plan)
		return
	}

	switch matches[passengerPlanRegexpAction] {
	case "Board":
		if p.PositionType == PassengerPositionTrain {
			e <- fmt.Errorf("passenger (%s): can not board, already on a train", p.ID)
			return
		}
		train, ok := w.Trains[matches[passengerPlanRegexpID]]
		if !ok {
			e <- fmt.Errorf("passenger (%s): can not find train %s", p.ID, matches[passengerPlanRegexpID])
			return
		}
		train.L.Lock()
		defer train.L.Unlock()
		if !train.BoardingPossible {
			e <- fmt.Errorf("passenger (%s): boarding not possible at %s", p.ID, train.ID)
			return
		}
		if train.Position[0] != p.Position {
			e <- fmt.Errorf("passenger (%s): train is not at station %s (currently: %s)", p.ID, p.Position, train.Position[0])
			return
		}
		train.Passengers.Add(&train.Passengers, &p.Size)
		p.PositionType = PassengerPositionTrain
		p.Position = train.ID
		p.TargetReached = big.Int{}
	case "Detrain":
		if p.PositionType == PassengerPositionStation {
			e <- fmt.Errorf("passenger (%s): can not detrain, already at a station", p.ID)
			return
		}
		train, ok := w.Trains[p.Position]
		if !ok {
			e <- fmt.Errorf("passenger (%s): can not find train '%s'", p.ID, p.Position)
			return
		}
		train.L.Lock()
		defer train.L.Unlock()
		if !train.BoardingPossible {
			e <- fmt.Errorf("passenger (%s): boarding not possible at %s", p.ID, train.ID)
			return
		}
		station, ok := w.Stations[train.Position[0]]
		if !ok {
			e <- fmt.Errorf("passenger (%s): can not find station %s", p.ID, train.Position[0])
			return
		}
		train.Passengers.Sub(&train.Passengers, &p.Size)
		p.PositionType = PassengerPositionStation
		p.Position = station.ID
		if p.Position == p.Target {
			targetTime := new(big.Int).Set(&w.CurrentTime)
			p.TargetReached = *targetTime
		}
	default:
		e <- fmt.Errorf("passenger (%s): unknown action '%s'", p.ID, matches[passengerPlanRegexpAction])
		return
	}
}
