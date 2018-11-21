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

	log "github.com/sirupsen/logrus"
)

func newTestDB() (Datastore, error) {
	db, err := NewDB("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	err = db.Initialize()
	if err != nil {
		return nil, err
	}

	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}

	return db, nil
}

func TestDB(t *testing.T) {
	_, err := newTestDB()
	if err != nil {
		t.Fatal(err)
	}
}
