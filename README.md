# myoxi - Record data from Pulse Oximeters

myoxi is a command line tool for storing and analyzing data from Pulse
Oximeters. Currently the only supported device is the Contec CMS50F.

## Features

- Stores heart rate and Oxygen Saturation (SpO2) in sqlite database
- Report statistics from previous sessions including Average Pulse, SpO2, and
  oxygen desaturation index.
- Export data in plain text or json format

## Getting started

The following assumes you have a Contec CMS50F Pulse Oximeter device. Each time
you record a session using the CMS50F it erases your previous session so be sure
to download the data before starting to record your next session.

- Download the myoxi binary release for your platform
  [here](https://github.com/aebruno/myoxi/releases)

- Be sure to set the correct time on your CMS50F before starting a session

- Record your session by turning on the "Record" option from the main menu. When
  finished stop recording.

- Plugin your CMS50F to your computer using the USB cable

- Find the path to the device file. On linux this will typically be
  `/dev/ttyUSB0`.

- Check the connectivity and device status by running the following command.
  This will connect to the device and display the model info and last session
  info:

```
	$ ./myoxi --port /dev/ttyUSB0 device 
	INFO[0000] Using device port: /dev/ttyUSB0              
	INFO[0000] Successfully connected to device at /dev/ttyUSB0 
	Device model: 50F
	Userinfo: user
	Session count: 1
	------------------------------
	- 2018-11-24 00:23:46 -0500 EST (7h37m8s)
```

- Download the latest session data from the device into the myoxi database
  run the following command:

```
	$ ./myoxi --port /dev/ttyUSB0 import 
	INFO[0000] Database path: /home/username/.myoxi.db           
	INFO[0000] Successfully opened myoxi database           
	INFO[0000] Using device port: /dev/ttyUSB0              
	INFO[0000] Successfully connected to device at /dev/ttyUSB0 
	INFO[0000] Found 1 sessions                             
	INFO[0000] Importing data for session 3 - 2018-11-24 00:23:46 -0500 EST (7h37m8s) 
	INFO[0015] Downloaded 27429 records. Total duration in seconds 27428.00 
	INFO[0015] Saving records to database 
```

- View the statistics from the last session run:

```
	$ ./myoxi stats
	INFO[0000] Database path: /home/username/.myoxi.db           
	INFO[0000] Successfully opened myoxi database           
	------------------------------------------------------
	Start: 2018-11-24 00:23:46 End: 2018-11-24 08:00:53
	------------------------------------------------------
	Total Records: 27428 (n = 27426, bad data = 2)
	Average SpO2 %: 95.94 (min: 88 max: 100 sd: 1.70)
	Average Pulse Rate: 61.91 (min: 50 max: 103 sd: 5.94)
	ODI: 10.38
	CT90: 7s
	Oxygen Desaturation Events = 7
	------------------------------------------------------
	11-24 04:26:55 lasting 8s desaturation 94.89 to 90.00
	11-24 04:39:09 lasting 10s desaturation 94.94 to 90.00
	11-24 04:57:15 lasting 11s desaturation 94.20 to 90.00
	11-24 05:13:22 lasting 6s desaturation 95.26 to 91.00
	11-24 05:13:38 lasting 8s desaturation 95.26 to 91.00
	11-24 05:33:07 lasting 5s desaturation 94.27 to 90.00
	11-24 07:58:52 lasting 35s desaturation 96.66 to 90.46
```

- See help for more reporting options:

```
	$ ./myoxi stats --help
	NAME:
	   myoxi stats - Display database stats

	USAGE:
	   myoxi stats [command options] [arguments...]

	OPTIONS:
	   --all, -a      Display stats for all data
	   --prev, -p     Display stats for previous session
	   --week, -w     Display stats for last week
	   --month, -m    Display stats for last month
	   --quarter, -q  Display stats for last quarter
	   --year, -y     Display stats for last year
```

## Building from source

myoxi is written in Go and requires v1.11 or greater. Clone the repository:

```
    $ git clone https://github.com/aebruno/myoxi
    $ cd myoxi
    $ go build ./...
```

## Acknowledgements

The code in device/cms50f.go was adopted from the SleepLib oximeter loader
plugin found in the free and open-source software SleepyHead available from
http://sleepyhead.jedimark.net, developed and copyright by Mark Watkins (C)
2011-2018.

## License

myoxi is released under the GPLv3 license. See the LICENSE file.
