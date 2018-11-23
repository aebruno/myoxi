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
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Session struct {
	ID        int64     `db:"id" json:"id"`
	StartTime time.Time `db:"start_time" json:"start_time"`
	Model     string    `db:"model" json:"model"`
	Seconds   int       `db:"duration_seconds" json:"duration_seconds"`
}

func (s *Session) String() string {
	return fmt.Sprintf(
		"ID=%d StartTime=%s Model=%s Duration=%s",
		s.ID,
		s.StartTime.Format("2006-01-02 15:04:05"),
		s.Model,
		time.Duration(time.Second*time.Duration(s.Seconds)))
}

func (db *DB) SaveSession(session *Session) error {
	res, err := db.NamedExec(`
        insert into session (start_time, model, duration_seconds) 
        values (:start_time, :model, :duration_seconds)`, session)
	if err != nil {
		return err
	}

	session.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) FetchSessionByStartTime(start time.Time) (*Session, error) {
	query := `
        select
			id,
			start_time,
			model,
            duration_seconds
        from session
        where start_time = ?
	`

	log.Debugf("Fetch Last Session by start time query: %s", query)

	session := &Session{}
	err := db.Get(session, query, start)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return session, nil
}

func (db *DB) FetchLastSession() (*Session, error) {
	query := `
        select
			id,
			start_time,
			model,
            duration_seconds
        from session
        order by start_time desc
        limit 1
	`

	log.Debugf("Fetch Last Session query: %s", query)

	session := &Session{}
	err := db.Get(session, query)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return session, nil
}

func (db *DB) FetchAllSessions() ([]*Session, error) {
	query := `
        select
			id,
			start_time,
			model,
            duration_seconds
        from session
	`
	log.Debugf("Fetch All Sessions query: %s", query)

	sessions := []*Session{}
	err := db.Select(&sessions, query)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}
