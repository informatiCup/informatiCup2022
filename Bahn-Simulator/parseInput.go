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
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"
)

type InputMode int

const (
	InputUnknown = iota
	InputLines
	InputStations
	InputTrains
	InputPassengers
)

var (
	inputStationsRegexp               = regexp.MustCompile(`\A(?P<id>[a-zA-Z0-9_]+) (?P<kapazitaet>[\d]+)[\s]*\z`)
	inputStationsRegexpID             = inputStationsRegexp.SubexpIndex("id")
	inputStationsRegexpKapazität      = inputStationsRegexp.SubexpIndex("kapazitaet")
	inputLinesRegexp                  = regexp.MustCompile(`\A(?P<id>[a-zA-Z0-9_]+) (?P<anfang>[a-zA-Z0-9_]+) (?P<ende>[a-zA-Z0-9_]+) (?P<laenge>[\d]+\.?[\d]*) (?P<kapazitaet>[\d]+)[\s]*\z`)
	inputLinesRegexpID                = inputLinesRegexp.SubexpIndex("id")
	inputLinesRegexpAnfang            = inputLinesRegexp.SubexpIndex("anfang")
	inputLinesRegexpEnde              = inputLinesRegexp.SubexpIndex("ende")
	inputLinesRegexpLänge             = inputLinesRegexp.SubexpIndex("laenge")
	inputLinesRegexpKapazität         = inputLinesRegexp.SubexpIndex("kapazitaet")
	inputTrainsRegexp                 = regexp.MustCompile(`\A(?P<id>[a-zA-Z0-9_]+) (?P<start>([a-zA-Z0-9_]+|\*)) (?P<geschwindigkeit>[\d]+\.?[\d]*) (?P<kapazitaet>[\d]+)[\s]*\z`)
	inputTrainsRegexpID               = inputTrainsRegexp.SubexpIndex("id")
	inputTrainsRegexpStart            = inputTrainsRegexp.SubexpIndex("start")
	inputTrainsRegexpGeschwindigkeit  = inputTrainsRegexp.SubexpIndex("geschwindigkeit")
	inputTrainsRegexpKapazität        = inputTrainsRegexp.SubexpIndex("kapazitaet")
	inputPassengersRegexp             = regexp.MustCompile(`\A(?P<id>[a-zA-Z0-9_]+) (?P<startbahnhof>[a-zA-Z0-9_]+) (?P<zielbahnhof>[a-zA-Z0-9_]+) (?P<gruppengroesse>[\d]+) (?P<ankunftszeit>[\d]+)[\s]*\z`)
	inputPassengersRegexpID           = inputPassengersRegexp.SubexpIndex("id")
	inputPassengersRegexpStartbahnhof = inputPassengersRegexp.SubexpIndex("startbahnhof")
	inputPassengersRegexpZielbahnhof  = inputPassengersRegexp.SubexpIndex("zielbahnhof")
	inputPassengersRegexpGruppengröße = inputPassengersRegexp.SubexpIndex("gruppengroesse")
	inputPassengersRegexpAnkunftszeit = inputPassengersRegexp.SubexpIndex("ankunftszeit")
)

func ParseInput(path string) (*World, error) {
	w := World{
		Lines:      make(map[string]*Line),
		Stations:   make(map[string]*Station),
		Trains:     make(map[string]*Train),
		Passengers: make(map[string]*Passenger),
	}

	tempStationCurrentCount := make(map[string]*big.Int)

	currentInputMode := InputUnknown

	f, err := os.Open(path)
	if err != nil {
		return &World{}, nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.HasPrefix(s, "#") {
			// Comment
			continue
		}
		if s == "" {
			// Empty line
			continue
		}
		if s == "[Stations]" {
			currentInputMode = InputStations
			continue
		}
		if s == "[Lines]" {
			currentInputMode = InputLines
			continue
		}
		if s == "[Trains]" {
			currentInputMode = InputTrains
			continue
		}
		if s == "[Passengers]" {
			currentInputMode = InputPassengers
			continue
		}

		switch currentInputMode {
		case InputUnknown:
			return &w, fmt.Errorf("can not parse '%s': no prior definition found", s)
		case InputLines:
			matches := inputLinesRegexp.FindStringSubmatch(s)
			if matches == nil {
				return &w, fmt.Errorf("can not parse '%s': not matching definition for line", s)
			}
			var l Line
			l.ID = matches[inputLinesRegexpID]
			if l.ID == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid id", s)
			}
			l.End = make([]string, 2)
			l.End[0] = matches[inputLinesRegexpAnfang]
			l.End[1] = matches[inputLinesRegexpEnde]
			if l.End[0] == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid start", s)
			}
			if l.End[1] == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid end", s)
			}
			if l.End[0] == l.End[1] {
				return &w, fmt.Errorf("can not parse '%s': start and end same station", s)
			}
			length, ok := new(big.Rat).SetString(matches[inputLinesRegexpLänge])
			if !ok {
				return &w, fmt.Errorf("can not parse '%s': can not parse length", s)
			}
			l.Length = *length
			capacity, ok := new(big.Int).SetString(matches[inputLinesRegexpKapazität], 10)
			if !ok {
				return &w, fmt.Errorf("can not parse '%s': can not parse capacity", s)
			}
			l.MaxCapacity = *capacity

			_, ok = w.Lines[l.ID]
			if ok {
				return &w, fmt.Errorf("can not parse '%s': id found twice", s)
			}

			w.Lines[l.ID] = &l
		case InputStations:
			var st Station
			matches := inputStationsRegexp.FindStringSubmatch(s)
			if matches == nil {
				return &w, fmt.Errorf("can not parse '%s': not matching definition for line", s)
			}
			st.ID = matches[inputStationsRegexpID]
			if st.ID == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid id", s)
			}
			capacity, ok := new(big.Int).SetString(matches[inputStationsRegexpKapazität], 10)
			if !ok {
				return &w, fmt.Errorf("can not parse '%s': invalid capacity", s)
			}
			st.Capacity = *capacity
			_, ok = w.Stations[st.ID]
			if ok {
				return &w, fmt.Errorf("stations %s found twice", s)
			}
			w.Stations[st.ID] = &st
		case InputTrains:
			matches := inputTrainsRegexp.FindStringSubmatch(s)
			if matches == nil {
				return &w, fmt.Errorf("can not parse '%s': not matching definition for line", s)
			}
			var t Train

			t.ID = matches[inputTrainsRegexpID]
			if t.ID == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid id", s)
			}
			capacity, ok := new(big.Int).SetString(matches[inputTrainsRegexpKapazität], 10)
			if !ok {
				return &w, fmt.Errorf("can not parse '%s': invalid capacity", s)
			}
			t.Capacity = *capacity
			speed, ok := new(big.Rat).SetString(matches[inputTrainsRegexpGeschwindigkeit])
			if !ok {
				return &w, fmt.Errorf("can not parse '%s': invalid speed", s)
			}
			t.Speed = *speed
			t.Position = []string{matches[inputTrainsRegexpStart]}
			if t.Position[0] == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid position", s)
			}
			if t.Position[0] == "*" {
				t.PositionType = TrainPositionWildcard
			} else {
				t.PositionType = TrainPositionStation
				c := tempStationCurrentCount[t.Position[0]]
				if c == nil {
					c = big.NewInt(0)
					tempStationCurrentCount[t.Position[0]] = c
				}
				c.Add(c, big.NewInt(1))
			}
			t.Plan = make(map[string]string)
			_, ok = w.Trains[t.ID]
			if ok {
				return &w, fmt.Errorf("can not parse '%s': id found twice", s)
			}
			w.Trains[t.ID] = &t
		case InputPassengers:
			matches := inputPassengersRegexp.FindStringSubmatch(s)
			if matches == nil {
				return &w, fmt.Errorf("can not parse '%s': not matching definition for line", s)
			}
			var p Passenger

			p.ID = matches[inputPassengersRegexpID]
			if p.ID == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid id", s)
			}

			p.Start = matches[inputPassengersRegexpStartbahnhof]
			if p.Start == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid start", s)
			}

			p.Target = matches[inputPassengersRegexpZielbahnhof]
			if p.Target == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid target", s)
			}

			size, ok := new(big.Int).SetString(matches[inputPassengersRegexpGruppengröße], 10)
			if !ok {
				return &w, fmt.Errorf("can not parse '%s': invalid size", s)
			}
			p.Size = *size

			targetTime, ok := new(big.Int).SetString(matches[inputPassengersRegexpAnkunftszeit], 10)
			if !ok {
				return &w, fmt.Errorf("can not parse '%s': invalid target time", s)
			}
			p.TargetTime = *targetTime

			p.PositionType = PassengerPositionStation

			p.Position = matches[inputPassengersRegexpStartbahnhof]
			if p.Position == "" {
				return &w, fmt.Errorf("can not parse '%s': invalid start", s)
			}

			p.Plan = make(map[string]string)

			_, ok = w.Passengers[p.ID]
			if ok {
				return &w, fmt.Errorf("can not parse '%s': id found twice", s)
			}
			w.Passengers[p.ID] = &p
		default:
			return &w, fmt.Errorf("[internal] unknown current state")
		}
	}
	if err := scanner.Err(); err != nil {
		return &w, err
	}

	for k := range tempStationCurrentCount {
		s, ok := w.Stations[k]
		if !ok {
			return &w, fmt.Errorf("can not assign current number of trains to non-existing station '%s'", k)
		}
		s.CurrenTrains.Add(&s.CurrenTrains, tempStationCurrentCount[k])
	}

	return &w, nil
}
