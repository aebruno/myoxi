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
	"math"
	"time"

    . "github.com/logrusorgru/aurora"
	"github.com/aebruno/myoxi/model"
)

func Stats(db model.Datastore) error {
    var pulseMin, pulseMax, spo2Min, spo2Max uint8
    var n, pulseSum, spo2Sum float64
    var pulseMean, pulseSD, spo2Mean, spo2SD float64

    spo2Min, pulseMin = math.MaxUint8, math.MaxUint8

	records, err := db.FetchRecords(time.Time{}, time.Time{})
	if err != nil {
        return err
	}

    for _, rec := range records {
        // Throw away data that's not within reasonable physical limits
        // TODO: perhaps make these configurable?
        if rec.Pulse < 40 {
            continue
        }
        if rec.Spo2 < 65 {
            continue
        }

        pulseSum += float64(rec.Pulse)
        spo2Sum += float64(rec.Spo2)
        if rec.Pulse > pulseMax {
            pulseMax = rec.Pulse
        }
        if rec.Spo2 > spo2Max {
            spo2Max = rec.Spo2
        }
        if rec.Pulse < pulseMin {
            pulseMin = rec.Pulse
        }
        if rec.Spo2 < spo2Min {
            spo2Min = rec.Spo2
        }
        n++
    }

    pulseMean = pulseSum/n
    spo2Mean = spo2Sum/n

    for _, rec := range records {
        pulseSD += math.Pow(float64(rec.Pulse) - pulseMean, 2)
        spo2SD += math.Pow(float64(rec.Spo2) - spo2Mean, 2)
    }

    pulseSD = math.Sqrt(pulseSD/n)
    spo2SD = math.Sqrt(spo2SD/n)

    fmt.Printf("------------------------------------------------------\n")
    fmt.Printf("Start: %s End: %s\n", records[0].DateTime.Format(time.Stamp), records[len(records)-1].DateTime.Format(time.Stamp))
    fmt.Printf("------------------------------------------------------\n")
    fmt.Printf("Total Records: %d (n = %d, bad data = %d)\n", len(records), int(n), len(records)-int(n))
    fmt.Printf("Average SpO2 %%: %.2f (min: %d max: %d sd: %.2f)\n", Bold(Blue(spo2Mean)), spo2Min, spo2Max, spo2SD)
    fmt.Printf("Average Pulse Rate: %.2f (min: %d max: %d sd: %.2f)\n", Bold(Red(pulseMean)), pulseMin, pulseMax, pulseSD)
    fmt.Printf("\n")

	return nil
}
