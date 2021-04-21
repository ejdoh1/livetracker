# LiveTracker

A program that polls the Garmin LiveTrack API and creates a localhost API with a CSV version of the latest data point.

## Running

First, run the livetracker binary (MacOS) or livetracker.exe (Windows) with command line options below if required.

Then, access the CSV data on localhost:8080 (default).

## Command line options

```sh
mac:livetracker user$ ./livetracker -h
Usage of ./livetracker:
  -i int
        interval sec (default 1)
  -p string
        port (default "8080")
  -s string
        session (default "6949868c-b1f0-444d-a1d1-9e50b06dcd18")
```


