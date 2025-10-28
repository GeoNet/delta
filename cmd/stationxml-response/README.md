# stationxml-response

Build a simplified JSON response file from StationXML input.

## Overview

The stationxml-response application takes a StationXML input file and builds a simplified response file for standard seismic
site installations. It assumes a single sensor poles and zeros entry followed by a list of datalogger anti-alias filters used
to decimate the incoming signal.

It's main use is for efficient instrument response without the overhead of handling large StationXML files. It currently
assumes that the StationXML is at version 1.2

## Usage:

  ./stationxml-response [options]

## Options:

  -input string
        provide an input StationXML file
  -match value
        provide a regexp match for input channels (default ^[ESHBL][HN][ZNE12]$)
  -output string
        provide an output JSON simplified response file
