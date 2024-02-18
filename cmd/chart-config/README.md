# chart-config

Build a chart plotting stream config file from delta meta & response information.

## Overview

The chart plotting software currently uses miniseed data as its input for plotting. To convert this to meaningful data, and to allow
annotations some form of meta data is required.

Two forms of metadata are compiles: the station and network descriptions; and the stream's scaling response to convert signal values.

Extra streams can be provided, so long as the station details have been loaded into delta, this requires the stream srcname (NN_SSSS_LL_CCC)
and the sampling frequency (e.g. `IU_RAR_10_BHZ:40`).

## Usage:

  ./chart-config [options]

## Options:

  -base string
        delta base files
  -channels value
        channel selection regexp (default .*)
  -exclude value
        station exclusion regexp
  -extra string
        extra streams to include
  -include value
        station inclusion regexp
  -locations value
        location selection regexp (default .*)
  -networks value
        network selection regexp (default .*)
  -output string
        output chart configuration file
  -primary string
        add phase constituent for tsunami streams (default "M2")
  -resp string
        base directory for response xml files on disk
  -single
        only add one stream per station
  -skip string
        extra streams to exclude
  -stations value
        station selection regexp (default .*)

## Example

This is the output for the station `CAW` and the channel `EHZ`:

```
[
  {
    "srcname": "NZ_CAW_10_EHZ",
    "network-code": "NZ",
    "station-code": "CAW",
    "location-code": "10",
    "channel-code": "EHZ",
    "station-name": "Cannon Point",
    "internal-network": "WL",
    "network-description": "Wellington regional seismic network",
    "latitude": -41.107194232,
    "longitude": 175.066438523,
    "sampling-period": 10000000,
    "sensitivity": 167772160,
    "gain": 1,
    "input_units": "m/s",
    "output_units": "count"
  }
]
```

