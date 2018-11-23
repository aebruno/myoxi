// Copyright 2018 Andrew E. Bruno
//
// This file is part of myoxi.
//
// myoxi is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// myoxi is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with myoxi.  If not, see <http://www.gnu.org/licenses/>.

package model

import (
	"testing"
	"time"
)

func TestOxiRecord(t *testing.T) {
	db, err := newTestDB()
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()

	data := []*OxiRecord{
		&OxiRecord{DateTime: start.Add(time.Second * 1), Pulse: 77, Spo2: 98, SessionID: 1},
		&OxiRecord{DateTime: start.Add(time.Second * 2), Pulse: 78, Spo2: 96, SessionID: 1},
		&OxiRecord{DateTime: start.Add(time.Second * 3), Pulse: 76, Spo2: 95, SessionID: 1},
		&OxiRecord{DateTime: start.Add(time.Second * 4), Pulse: 79, Spo2: 98, SessionID: 1},
		&OxiRecord{DateTime: start.Add(time.Second * 5), Pulse: 77, Spo2: 99, SessionID: 1},
		&OxiRecord{DateTime: start.Add(time.Second * 6), Pulse: 79, Spo2: 94, SessionID: 1},
	}

	err = db.SaveRecords(data)
	if err != nil {
		t.Error(err)
	}

	queries := [][]time.Time{
		[]time.Time{time.Time{}, time.Time{}},
		[]time.Time{start, start.Add(time.Second * 7)},
	}

	for _, query := range queries {
		records, err := db.FetchRecords(query[0], query[1])
		if err != nil {
			t.Error(err)
		}

		if len(records) != len(data) {
			t.Errorf("Invalid number of records returned. Got %d wanted %d", len(records), len(data))
		}

		for i := range records {
			if records[i].DateTime.UTC() != data[i].DateTime.UTC() {
				t.Errorf("Invalid datetime for record %d. Got %s wanted %s", i, records[i].DateTime.UTC(), data[i].DateTime.UTC())
			}
			if records[i].Pulse != data[i].Pulse {
				t.Errorf("Invalid pulse for record %d. Got %d wanted %d", i, records[i].Pulse, data[i].Pulse)
			}
			if records[i].Spo2 != data[i].Spo2 {
				t.Errorf("Invalid spo2 for record %d. Got %d wanted %d", i, records[i].Spo2, data[i].Spo2)
			}
		}
	}

	records, err := db.FetchRecords(start.Add(time.Second*100), time.Time{})
	if err != nil {
		t.Error(err)
	}

	if len(records) != 0 {
		t.Errorf("Invalid number of records returned. Got %d wanted %d", len(records), 0)
	}

	records, err = db.FetchRecords(time.Time{}, start.Add(-time.Second*100))
	if err != nil {
		t.Error(err)
	}

	if len(records) != 0 {
		t.Errorf("Invalid number of records returned. Got %d wanted %d", len(records), 0)
	}

	records, err = db.FetchRecordsBySessionID(1)
	if err != nil {
		t.Error(err)
	}

	if len(records) != len(data) {
		t.Errorf("Invalid number of records returned for fetch by sessionID. Got %d wanted %d", len(records), len(data))
	}
}
