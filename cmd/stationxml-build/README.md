# stationxml-build

Build a network StationXML file from delta meta & response information.

## Overview

The stationxml-build application can be run either with the compiled in delta files, or by pointing at a delta file base directory.

The first option is designed for automatic configuration generators running in a CI/CD like environment whereas the second is aimed
at manual runs, it is important to provide a base in this instance as otherwise the built in, and likely to be old, files will be used.

The output of the program can be either of two types,

- A complete StationXML file given by the _output_ flag, or `stdout`.
- A set of StationXML files per station, this is aimed at configuring applications such as SeisComP using individual StationXML files.

For the second option, this is _purge_ option which will maintain the contents of the given _directory_ used for the _single_ option. It
will automatically remove any non-active stations.

Sometimes, e.g. for automatic systems, it is better to not create an XML _Created_ entry, as this will make the file unique and will always
indicate a new file has been created, even if the metadata content hasn't changed. 

The response information is read via small XML files representing dataloggers and sensors, these are joined together on the fly as needed.
The application will come with compiled in versions, again this may need to be updated if new responses are added. Alternatively this can
be overridden using a standalone directory.

The versioning mechanism is based around a single structure that incorporates all parameter options, and when building a particular version
this structure is mapped into the actual version layout and saved via XML encoding.

It is assumed that the small response files will not vary with version changes.


## Usage:

  ./stationxml-build [options]

## Options:

  -verbose
        _add operational info_

  -debug
        _add extra operational info_

  -base string
        _base of delta files on disk_

  -resp string
        _base for response xml files on disk_

  -version string
        _create a specific StationXML version_

  -create
        _add a root XML "Created" entry_

  -module string
        _stationxml module (default "Delta")_

  -sender string
        _stationxml sender (default "WEL(GNS_Test)")_

  -source string
        _stationxml source (default "GeoNet")_

  -output string
        _output xml file, use "-" for stdout_

  -directory string
        _where to store station xml files (default "xml")_

  -single
        _produce single station xml files_

  -purge
        _remove unknown single xml files_

  -template string
        _how to name the single station xml files (default "station_{{.ExternalCode}}_{{.StationCode}}.xml")_

  -ignore string
        _list of stations to skip_

  -exclude value
        _regexp selection of networks to exclude (default ^()$)_

  -external value
        _regexp selection of external networks (default ^(NZ)$)_

  -station value
        _regexp selection of stations (default [A-Z0-9]+)_

  -network value
        _regexp selection of networks (default [A-Z0-9]+)_

  -location value
        _regexp selection of locations (default [A-Z0-9]+)_

  -channel value
        _regexp selection of channels (default [A-Z0-9]+)_

