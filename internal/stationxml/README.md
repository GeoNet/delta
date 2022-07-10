# FDSNStationXML

Automatically generate _FDSN StationXML_ golang bindings from their respective Schema `xsd` files.

Documentation on the XML schema format can be found at https://docs.fdsn.org/projects/stationxml

## Schemas Available

- v1.0 -- https://www.fdsn.org/xml/station/fdsn-station-1.0.xsd
- v1.1 -- https://www.fdsn.org/xml/station/fdsn-station-1.1.xsd
- v1.2 -- https://www.fdsn.org/xml/station/fdsn-station-1.2.xsd

## Usage

The modules provide a `stationxml` package and can be imported using:

```golang
import (
        github.com/geonet/internal/stationxml/v1.2
)
```

## Update

The code to build the golang encodings is in the `generate` subdirectory with the actual
commands in the `doc.go` file.

To add a new version, update the `doc.go` code to add a new generate line.

The schema files are expected to be found at:

        https://www.fdsn.org/xml/station/fdsn-station-1.2.xsd

with the appropriate version substitued.

Then in the current directory run:

        `go generate`

The doc.go file and generated code, in the versioned directories, should be commited to git.
