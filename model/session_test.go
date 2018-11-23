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

func TestSession(t *testing.T) {
	db, err := newTestDB()
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()

	data := []*Session{
		&Session{StartTime: start, Model: "50F", Seconds: 3600},
		&Session{StartTime: start.Add(-time.Second * 86400), Model: "50F", Seconds: 28800},
	}

	for _, s := range data {
		err = db.SaveSession(s)
		if err != nil {
			t.Error(err)
		}
	}

	err = db.SaveSession(data[0])
	if err == nil {
		t.Errorf("Unique constraint failed. Should not be allowed to save session with same start time")
	}

	sessions, err := db.FetchAllSessions()
	if err != nil {
		t.Error(err)
	}

	if len(sessions) != len(data) {
		t.Errorf("Invalid number of sessions returned. Got %d wanted %d", len(sessions), len(data))
	}

	session, err := db.FetchLatestSession()
	if err != nil {
		t.Error(err)
	}

	if session.ID != 1 {
		t.Errorf("Invalid session ID for latest session returned. Got %d wanted %d", session.ID, 1)
	}

	if session.StartTime.UTC() != data[0].StartTime.UTC() {
		t.Errorf("Invalid start time for latest session returned. Got %s wanted %s", session.StartTime.UTC(), data[0].StartTime.UTC())
	}

	session, err = db.FetchSessionByStartTime(start)
	if err != nil {
		t.Error(err)
	}

	if session.ID != 1 {
		t.Errorf("Invalid session ID for session returned. Got %d wanted %d", session.ID, 1)
	}

	session, err = db.FetchPreviousSession()
	if err != nil {
		t.Error(err)
	}

	if session.ID != 2 {
		t.Errorf("Invalid session ID for previous session returned. Got %d wanted %d", session.ID, 1)
	}
}
