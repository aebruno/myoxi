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

package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/aebruno/myoxi/device"
	"github.com/aebruno/myoxi/model"
	"github.com/aebruno/myoxi/tools"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func connectDevice(port string) (device.Device, error) {
	log.Infof("Using device port: %s", port)

	device := &device.CMS50{}

	err := device.Connect(port)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"port":  port,
		}).Error("Failed to open device port")

		return nil, err
	}

	log.Infof("Successfully connected to device at %s", port)

	return device, nil
}

func initDB(dbpath string) (model.Datastore, error) {
	if len(dbpath) == 0 {
		home := os.Getenv("HOME")
		if runtime.GOOS == "windows" {
			home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
			if home == "" {
				home = os.Getenv("USERPROFILE")
			}
		}
		dbpath = filepath.Join(home, ".myoxi.db")
	}

	log.Infof("Database path: %s", dbpath)

	db, err := model.NewDB("sqlite3", dbpath)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"dbpath": dbpath,
		}).Error("Failed to open database file")

		return nil, err
	}

	err = db.Initialize()
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"dbpath": dbpath,
		}).Error("Failed to initialize database")

		return nil, err
	}

	log.Infof("Successfully opened myoxi database")

	return db, nil
}

func setup(dbpath, port string) (model.Datastore, device.Device, error) {
	db, err := initDB(dbpath)
	if err != nil {
		return nil, nil, err
	}

	device, err := connectDevice(port)
	if err != nil {
		return nil, nil, err
	}

	return db, device, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "myoxi"
	app.Authors = []cli.Author{cli.Author{Name: "Andrew E. Bruno", Email: "aeb@qnot.org"}}
	app.Usage = "myoxi"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "debug,d", Usage: "Print debug messages"},
		&cli.StringFlag{Name: "port, p", Usage: "Path to device port", Value: "/dev/ttyUSB0"},
		&cli.StringFlag{Name: "dbpath, x", Usage: "Path to database file"},
	}
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}

		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "import",
			Usage: "Import data from device",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "noop, n", Usage: "Dump data only. Don't save to database"},
				&cli.BoolFlag{Name: "force, f", Usage: "Force overwrite session if exists"},
			},
			Action: func(c *cli.Context) error {
				db, device, err := setup(c.GlobalString("dbpath"), c.GlobalString("port"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				err = tools.Import(db, device, c.Bool("noop"), c.Bool("force"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		},
		{
			Name:  "stats",
			Usage: "Display database stats",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "all, a", Usage: "Display stats for all data"},
				&cli.BoolFlag{Name: "prev, p", Usage: "Display stats for previous session"},
				&cli.BoolFlag{Name: "week, w", Usage: "Display stats for last week"},
				&cli.BoolFlag{Name: "month, m", Usage: "Display stats for last month"},
				&cli.BoolFlag{Name: "quarter, q", Usage: "Display stats for last quarter"},
				&cli.BoolFlag{Name: "year, y", Usage: "Display stats for last year"},
			},
			Action: func(c *cli.Context) error {
				db, err := initDB(c.GlobalString("dbpath"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				now := time.Now()
				var records []*model.OxiRecord
				if c.Bool("all") {
					records, err = db.FetchRecords(time.Time{}, time.Time{})
				} else if c.Bool("prev") {
					session, err := db.FetchPreviousSession()
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					records, err = db.FetchRecordsBySessionID(session.ID)
				} else if c.Bool("week") {
					records, err = db.FetchRecords(now.Add(-24*7*time.Hour), now)
				} else if c.Bool("month") {
					records, err = db.FetchRecords(now.Add(-24*30*time.Hour), now)
				} else if c.Bool("year") {
					records, err = db.FetchRecords(now.Add(-24*365*time.Hour), now)
				} else {
					session, err := db.FetchLatestSession()
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					records, err = db.FetchRecordsBySessionID(session.ID)
				}

				if err != nil {
					return cli.NewExitError(err, 1)
				}

				tools.ComputeAndPrintStats(records)

				return nil
			},
		},
		{
			Name:  "device",
			Usage: "Display information about device",
			Action: func(c *cli.Context) error {
				device, err := connectDevice(c.GlobalString("port"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				err = tools.DeviceInfo(device)
				if err != nil {
					return cli.NewExitError(err, 1)
				}

				return nil
			},
		}}

	app.RunAndExitOnError()
}
