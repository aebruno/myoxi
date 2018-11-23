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
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	SessionSchema = `
		create table if not exists session 
		(id integer primary key, start_time datetime unique, model string, duration_seconds integer)
	`

	OxiRecordSchema = `
		create table if not exists oxi_record 
		(date_time datetime primary key, session_id integer not null, pulse integer, spo2 integer)
	`
)

var ErrNotFound = errors.New("Record not found in database")

type Datastore interface {
	Initialize() error
	SaveRecords(records []*OxiRecord) error
	FetchRecords(from, to time.Time) ([]*OxiRecord, error)
	SaveSession(session *Session) error
	FetchLastSession() (*Session, error)
	FetchSessionByStartTime(start time.Time) (*Session, error)
	FetchAllSessions() ([]*Session, error)
}

type DB struct {
	*sqlx.DB
}

func NewDB(driver, dsn string) (Datastore, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Initialize() error {
	_, err := db.Exec(SessionSchema)
	if err != nil {
		return err
	}

	_, err = db.Exec(OxiRecordSchema)
	if err != nil {
		return err
	}

	return nil
}
