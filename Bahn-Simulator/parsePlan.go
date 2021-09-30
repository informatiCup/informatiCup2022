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
	"strings"
)

type PlanState int

const (
	PlanUnknown PlanState = iota
	PlanPassenger
	PlanTrain
)

func ParsePlan(w *World, path string) error {
	currentState := PlanUnknown
	currentID := ""

	f, err := os.Open(path)
	if err != nil {
		return nil
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
		if strings.HasPrefix(s, "[") {
			s = strings.TrimSpace(s)
			s = strings.TrimPrefix(s, "[")
			s = strings.TrimSuffix(s, "]")
			split := strings.Split(s, ":")
			if len(split) != 2 {
				return fmt.Errorf("can not parse %s", s)
			}
			currentID = split[1]
			if split[0] == "Train" {
				currentState = PlanTrain
			} else if split[0] == "Passenger" {
				currentState = PlanPassenger
			} else {
				return fmt.Errorf("unknown type '%s'", split[0])
			}
			continue
		}

		switch currentState {
		case PlanUnknown:
			return fmt.Errorf("can not parse '%s': no prior definition found", s)
		case PlanPassenger:
			matches := passengerPlanRegexp.FindStringSubmatch(s)
			if matches == nil {
				return fmt.Errorf("can not parse '%s': not matching definition for line", s)
			}
			time, ok := new(big.Int).SetString(matches[passengerPlanRegexpTime], 10)
			if !ok {
				return fmt.Errorf("can not parse time '%s'", matches[passengerPlanRegexpTime])
			}
			if time.Cmp(big.NewInt(0)) != +1 {
				return fmt.Errorf("can not parse '%s': time '%s' must be positive", s, matches[passengerPlanRegexpTime])
			}
			p, ok := w.Passengers[currentID]
			if !ok {
				return fmt.Errorf("can not parse '%s': no valid passenger id (%s)", s, currentID)
			}
			_, ok = p.Plan[time.String()]
			if ok {
				return fmt.Errorf("can not parse '%s': time %s already in plan", s, time.String())
			}
			p.Plan[time.String()] = s

			if time.Cmp(&w.MaxTime) == +1 {
				maxtime := new(big.Int).Set(time)
				maxtime.Add(maxtime, big.NewInt(1))
				w.MaxTime = *maxtime
			}
		case PlanTrain:
			matches := trainPlanRegexp.FindStringSubmatch(s)
			if matches == nil {
				return fmt.Errorf("can not parse '%s': not matching definition for line", s)
			}
			time, ok := new(big.Int).SetString(matches[trainPlanRegexpTime], 10)
			if !ok {
				return fmt.Errorf("can not parse time '%s'", matches[trainPlanRegexpTime])
			}
			t, ok := w.Trains[currentID]
			if !ok {
				return fmt.Errorf("can not parse '%s': no valid train id (%s)", s, currentID)
			}
			switch time.Cmp(big.NewInt(0)) {
			case +1:
				_, ok = t.Plan[time.String()]
				if ok {
					return fmt.Errorf("can not parse '%s': time %s already in plan", s, time.String())
				}
				t.Plan[time.String()] = s

				if time.Cmp(&w.MaxTime) == +1 {
					maxtime := new(big.Int).Set(time)
					maxtime.Add(maxtime, big.NewInt(1))
					w.MaxTime = *maxtime
				}
			case 0:
				if matches[trainPlanRegexpAction] != "Start" {
					return fmt.Errorf("can not parse '%s': time %s must be 'Start'", s, time.String())
				}
				if t.PositionType != TrainPositionWildcard || len(t.Position) != 1 || t.Position[0] != "*" {
					return fmt.Errorf("can not parse '%s': train must be at '*' for 'Start'", s)
				}
				st, ok := w.Stations[matches[trainPlanRegexpID]]
				if !ok {
					return fmt.Errorf("can not parse '%s': station '%s' does not exist'", s, matches[trainPlanRegexpID])
				}
				t.Position = []string{st.ID}
				t.PositionType = TrainPositionStation
				st.CurrenTrains.Add(&st.CurrenTrains, big.NewInt(1))
			case -1:
				return fmt.Errorf("can not parse '%s': time '%s' must be positive", s, matches[trainPlanRegexpTime])
			}
		default:
			return fmt.Errorf("[internal] unknown current state")
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
