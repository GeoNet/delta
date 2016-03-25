# delta

Meta information for the GeoNet sensor network.

## network

Geographical and equipment grouping meta data.

## equipment

Physical equipment asset management.

## install

Equipment installation and connections details.

### sensors

A list of _sensor_ installations, these include values for:
- Sensor Make, Model &amp; Serial Numbers
- Station &amp; Site Codes
- Installation Azimuth &amp; Dips [&deg;]
- Installation Depth [m]
- Sensor Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### gauges

A list of _gauge_ installations, these include values for:
- Gauge Make, Model &amp; Serial Numbers
- Station &amp; Site Codes
- Installation Dips [&deg;]
- Installation Offsets from the specified Site [m]
- Cable Length [m]
- Gauge Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### dataloggers

A list of _datalogger_ installations, these include values for:
- Datalogger Make, Model &amp; Serial Numbers
- Deployment Place and optional Roles
- Datalogger Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### connections

A list of _datalogger_ connections, these include values for:
- Station &amp; Site Codes
- Deployment Place and optional Roles
- Connection Start &amp; End dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### recorders

A list of _recorder_ installations, these include values for:
- Recorder Make, Model &amp; Serial Numbers
- Station &amp; Site Codes
- Installation Azimuth &amp; Dip [&deg;]
- Installation Depth [m]
- Recorder Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### antennas

A list of _antenna_ installations, these include values for:
- Antenna Make, Model &amp; Serial Numbers
- Mark Code
- Installation Height &amp; Offsets [m]
- Antenna Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### radmomes

A list of _radome_ installations, these include values for:
- Radome Make, Model &amp; Serial Numbers
- Mark Code
- Radome Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### metsensors

A list of _metsensor_ installations, these include values for:
- Met Sensor Make, Model &amp; Serial Numbers
- Mark Code &amp; IMS Comments
- Installation Location [&deg;], Datum &amp; Heights [m]
- Met Sensor Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### receivers

A list of _receiver_ installations, these include values for:
- Receiver Make, Model &amp; Serial Numbers
- Mark Code
- Receiver Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]

### firmware

A list of _firmware_ versions, these include values for:
- Device Make, Model &amp; Serial Numbers
- Version Number
- Firmware Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]
- Extra Notes

### cameras

A list of _camera_ installations, these include values for:
- Camera Make, Model &amp; Serial Numbers
- Camera Site Codes
- Installation Dip &amp; Azimiths [&deg;]
- Installation Height &amp; Offsets
- Camera Installation &amp; Removal dates [<yyyy>-<mm>-<dd>T<hh>-<mm>-<ss>Z]
- Installation Notes

## meta

Golang module to load network list files (csv).

## tests

Consistency checking of the network meta data.
