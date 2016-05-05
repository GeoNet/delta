# DELTA

Meta information for the GeoNet sensor network.

## network

Geographical and equipment grouping meta data.

## equipment

Physical equipment asset management.

## install

Equipment installation and connections details.

### antennas

A list of _antenna_ installations, these include values for:
- Antenna Make, Model &amp; Serial Numbers
- Mark Code
- Installation Height, Offsets [_m_] &amp; Azimuth [_&deg;_]
- Antenna Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### cameras

A list of _camera_ installations, these include values for:
- Camera Make, Model &amp; Serial Numbers
- Camera Mount Codes
- Installation Dip &amp; Azimiths [_&deg;_]
- Installation Height &amp; Offsets
- Camera Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]
- Installation Notes

### connections

A list of _datalogger_ connections, these include values for:
- Station &amp; Site Location Codes
- Deployment Place and optional Roles
- Connection Start &amp; End dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### dataloggers

A list of _datalogger_ installations, these include values for:
- Datalogger Make, Model &amp; Serial Numbers
- Deployment Place and optional Roles
- Datalogger Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### firmware

A list of _firmware_ versions, these include values for:
- Device Make, Model &amp; Serial Numbers
- Version Number
- Firmware Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]
- Extra Notes

### gauges

A list of _gauge_ installations, these include values for:
- Gauge Make, Model &amp; Serial Numbers
- Station &amp; Site Codes
- Installation Dips [_&deg;_]
- Installation Offsets from the specified Site [_m_]
- Cable Length [_m_]
- Gauge Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### metsensors

A list of _metsensor_ installations, these include values for:
- Met Sensor Make, Model &amp; Serial Numbers
- Mark Code &amp; IMS Comments
- Installation Location [_&deg;_], Datum &amp; Heights [_m_]
- Met Sensor Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### radomes

A list of _radome_ installations, these include values for:
- Radome Make, Model &amp; Serial Numbers
- Mark Code
- Radome Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### receivers

A list of _receiver_ installations, these include values for:
- Receiver Make, Model &amp; Serial Numbers
- Mark Code
- Receiver Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### recorders

A list of _recorder_ installations, these include values for:
- Recorder Make, Model &amp; Serial Numbers
- Station &amp; Site Codes
- Installation Azimuth &amp; Dip [_&deg;_]
- Installation Depth [_m_]
- Recorder Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### sensors

A list of _sensor_ installations, these include values for:
- Sensor Make, Model &amp; Serial Numbers
- Station &amp; Site Codes
- Installation Azimuth &amp; Dips [_&deg;_]
- Installation Depth [_m_]
- Sensor Installation &amp; Removal dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

### sessions

A list of _gnss_ receiver &amp; antenna recording sessions, these include values for:
- Mark Code
- Network Operator
- Network Agency
- Satellite System
- Session Interval
- Elevation Mask
- Header Comment
- Session Start &amp; Stop dates [_&lt;yyyy&gt;-&lt;mm&gt;-&lt;dd&gt;T&lt;hh&gt;-&lt;mm&gt;-&lt;ss&gt;Z_]

## meta

Golang module to load network list files (csv).

## tests

Consistency checking of the network meta data.
