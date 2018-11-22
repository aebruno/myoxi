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

package tools

import (
	"fmt"
	"time"

	"github.com/aebruno/myoxi/device"
	"github.com/aebruno/myoxi/model"
	log "github.com/sirupsen/logrus"
)

func Import(db model.Datastore, device device.Device, noop bool) error {
	err := device.ResetDevice()
	if err != nil {
		return fmt.Errorf("Failed to reset device: %s", err)
	}

	count, err := device.GetSessionCount()
	if err != nil {
		return fmt.Errorf("Failed to reset device: %s", err)
	}

	log.Infof("Found %d sessions", count)

	if count == 0 {
		log.Warn("No sessions found. Nothing to import")
		return nil
	}

	for i := uint8(0); i < count; i++ {
		duration, err := device.GetSessionDuration(i)
		if err != nil {
			return fmt.Errorf("Failed to fetch session duration: %s", err)
		}

		startTime, err := device.GetSessionTime(i)
		if err != nil {
			return fmt.Errorf("Failed to fetch session time: %s", err)
		}

		log.Infof("Importing data for session %d - %s (%s)", i, startTime, duration)
		data, err := device.GetSessionData(i)
		if err != nil {
			return fmt.Errorf("Failed to connect to device: %s", err)
		}

		log.Infof("Downloaded %d records. Total duration in seconds %0.2f", len(data), duration.Seconds())

		total := int(duration.Seconds())
		if total > len(data) {
			log.WithFields(log.Fields{
				"numRecords":      len(data),
				"durationSeconds": total,
			}).Warn("Not enough records found for the session duration")
			total = len(data)
		}

		for i, rec := range data {
			rec.DateTime = startTime.Add(time.Second * time.Duration(i))
			if noop {
				fmt.Printf("Record %d - %s\n", i, rec)
			} else {
				log.Debugf("Record %d - %s", i, rec)
			}
		}

		if noop {
			return nil
		}

		log.Infof("Saving records to database")
		err = db.SaveRecords(data[:total])
		if err != nil {
			return fmt.Errorf("Failed to save records to database: %s", err)
		}
	}

	return nil
}
