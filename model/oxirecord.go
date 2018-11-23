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
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type OxiRecord struct {
	DateTime  time.Time `db:"date_time" json:"date_time"`
	SessionID int64     `db:"session_id" json:"session_id"`
	Pulse     uint8     `db:"pulse" json:"pulse"`
	Spo2      uint8     `db:"spo2" json:"spo2"`
}

func (r *OxiRecord) String() string {
	return fmt.Sprintf("DateTime=%s Pulse=%d SPO2=%d", r.DateTime.Format("2006-01-02 15:04:05"), r.Pulse, r.Spo2)
}

func (db *DB) SaveRecords(records []*OxiRecord) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Commit()

	for _, record := range records {
		_, err := tx.NamedExec(`
            replace into oxi_record (date_time, session_id, pulse, spo2) 
            values (:date_time, :session_id, :pulse, :spo2)`, record)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) FetchRecords(from, to time.Time) ([]*OxiRecord, error) {
	args := make([]interface{}, 0)
	query := `
        select
			date_time,
            session_id,
			pulse,
            spo2
        from oxi_record
	`

	nullTime := time.Time{}

	if from != nullTime && to != nullTime {
		query += ` where date_time > ? and date_time < ?`
		args = append(args, from)
		args = append(args, to)
	} else if from != nullTime {
		query += ` where date_time > ?`
		args = append(args, from)
	} else if to != nullTime {
		query += ` where date_time < ?`
		args = append(args, to)
	}

	query += ` order by date_time asc`

	log.Debugf("Fetch Records args: %v query: %s", args, query)

	data := []*OxiRecord{}
	err := db.Select(&data, query, args...)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (db *DB) FetchRecordsBySessionID(sessionID int64) ([]*OxiRecord, error) {
	query := `
        select
			date_time,
            session_id,
			pulse,
            spo2
        from oxi_record
        where session_id = ?
	`

	log.Debugf("Fetch Records by session id query: %s", query)

	data := []*OxiRecord{}
	err := db.Select(&data, query, sessionID)
	if err != nil {
		return nil, err
	}

	return data, nil
}
