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

	"github.com/aebruno/myoxi/model"
	. "github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
)

type DesaturationEvent struct {
	start   time.Time
	end     time.Time
	avg120  float64
	records []*model.OxiRecord
}

type Stats struct {
	totalRecords int
	badRecords   int
	spo2Min      uint8
	spo2Max      uint8
	spo2Mean     float64
	spo2SD       float64
	pulseMin     uint8
	pulseMax     uint8
	pulseMean    float64
	pulseSD      float64
	odi          float64
	ct90         time.Duration
	events       []*DesaturationEvent
}

func (d *DesaturationEvent) String() string {
	sum := 0
	for _, rec := range d.records {
		sum += int(rec.Spo2)
	}
	mean := float64(sum) / float64(len(d.records))
	return fmt.Sprintf("%s lasting %s desaturation %.2f to %.2f", d.start.Format("01-02 15:04:05"), d.end.Sub(d.start), d.avg120, mean)
}

func computeODI(data []*model.OxiRecord) (float64, time.Duration, []*DesaturationEvent) {
	nullTime := time.Time{}
	avg120 := float64(95)
	sum120 := 0
	count120 := 1

	sumODI := 0
	countODI := 0

	events := make([]*DesaturationEvent, 0)
	curEvent := &DesaturationEvent{}
	ct90 := 0

	idx := 0
	processHours := true

	for processHours {
		hourODI := 0
		end := idx + 3600

		if end > len(data) {
			end = len(data)
			processHours = false
		}

		for _, rec := range data[idx:end] {
			// Throw away data that's not within reasonable physical limits
			if rec.Pulse < 40 {
				continue
			}
			if rec.Spo2 < 65 {
				continue
			}

			if rec.Spo2 < 90 {
				ct90++
			}

			if avg120-float64(rec.Spo2) >= 4 {
				log.Debugf("Oxygen desaturation event at %s: %d (%.2f 120s avg)", rec.DateTime.Format("01-02 15:04:05"), rec.Spo2, avg120)
				hourODI++
				if curEvent.start == nullTime {
					curEvent.start = rec.DateTime
					curEvent.records = append(curEvent.records, rec)
					curEvent.avg120 = avg120
				} else {
					curEvent.records = append(curEvent.records, rec)
				}
			} else if curEvent.start != nullTime {
				curEvent.end = rec.DateTime
				events = append(events, curEvent)
				curEvent = &DesaturationEvent{}
			}

			sum120 += int(rec.Spo2)
			count120++

			if count120 == 120 {
				avg120 = float64(sum120) / 120
				count120 = 1
				sum120 = 0
				log.Debugf("Avg 120: %.2f", avg120)
			}
		}

		idx += 3600
		sumODI += hourODI
		countODI++
	}

	odi := float64(sumODI) / float64(countODI)
	log.Debugf("Hours: %d, ODI: %.2f sum: %d", countODI, odi, sumODI)
	log.Debugf("Number of events: %d", len(events))
	for i, ev := range events {
		log.Debugf("Event %d: avg: %.2f start: %s end: %s", i, ev.avg120, ev.start, ev.end)
		for _, rec := range ev.records {
			log.Debugf("     - Record %s", rec)
		}
	}

	return odi, time.Duration(time.Second * time.Duration(ct90)), events
}

func ComputeStats(records []*model.OxiRecord) *Stats {
	var n, pulseSum, spo2Sum float64
	stats := &Stats{}
	stats.spo2Min, stats.pulseMin = math.MaxUint8, math.MaxUint8

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
		if rec.Pulse > stats.pulseMax {
			stats.pulseMax = rec.Pulse
		}
		if rec.Spo2 > stats.spo2Max {
			stats.spo2Max = rec.Spo2
		}
		if rec.Pulse < stats.pulseMin {
			stats.pulseMin = rec.Pulse
		}
		if rec.Spo2 < stats.spo2Min {
			stats.spo2Min = rec.Spo2
		}
		n++
	}

	stats.pulseMean = pulseSum / n
	stats.spo2Mean = spo2Sum / n

	for _, rec := range records {
		stats.pulseSD += math.Pow(float64(rec.Pulse)-stats.pulseMean, 2)
		stats.spo2SD += math.Pow(float64(rec.Spo2)-stats.spo2Mean, 2)
	}

	stats.pulseSD = math.Sqrt(stats.pulseSD / n)
	stats.spo2SD = math.Sqrt(stats.spo2SD / n)

	stats.odi, stats.ct90, stats.events = computeODI(records)
	stats.totalRecords = int(n)
	stats.badRecords = len(records) - int(n)

	return stats
}

func ComputeAndPrintStats(records []*model.OxiRecord) {
	stats := ComputeStats(records)

	fmt.Printf("------------------------------------------------------\n")
	fmt.Printf("Start: %s End: %s\n", records[0].DateTime.Format("2006-01-02 15:04:05"), records[len(records)-1].DateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("------------------------------------------------------\n")
	fmt.Printf("Total Records: %d (n = %d, bad data = %d)\n", len(records), stats.totalRecords, stats.badRecords)
	fmt.Printf("Average SpO2 %%: %.2f (min: %d max: %d sd: %.2f)\n", Bold(Blue(stats.spo2Mean)), stats.spo2Min, stats.spo2Max, stats.spo2SD)
	fmt.Printf("Average Pulse Rate: %.2f (min: %d max: %d sd: %.2f)\n", Bold(Red(stats.pulseMean)), stats.pulseMin, stats.pulseMax, stats.pulseSD)
	fmt.Printf("ODI: %.2f\n", Bold(Blue(stats.odi)))
	fmt.Printf("CT90: %s\n", Bold(stats.ct90))
	fmt.Printf("Oxygen Desaturation Events = %d\n", len(stats.events))
	fmt.Printf("------------------------------------------------------\n")
	for _, e := range stats.events {
		fmt.Printf("%s\n", e)
	}
	fmt.Printf("\n")
}
