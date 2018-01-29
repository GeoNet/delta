## INSTALL ##

_Equipment installation, configuration and connection details._

### FILES ###

Meta information for the GeoNet equipment network.
 
* `antennas.csv` - GNSS observation antennas
* `radomes.csv` -  GNSS installed antenna radomes
* `metsensors.csv` -  GNSS met sensors attached to receivers
* `receivers.csv` -  GNSS observation receivers
* `firmwares.csv` - GNSS receiver firmware versions
* `sessions.csv` -  GNSS receiver session configurations

* `sensors.csv` - Recording sensors
* `recorders.csv` - Combined sensor and datalogger recorders
* `dataloggers.csv` - Recording dataloggers
* `connections.csv` - Datalogger and sensor connection details
* `streams.csv` - Datalogger and recorder sampling configurations

* `cameras.csv` - Installed field cameras.

### GNSS ###

#### _ANTENNAS_ ####

A list of _antenna_ installations.

| Field | Description | Units |
| --- | --- | --- |
| _Make_ | Installed antenna make
| _Model_ | Installed antenna model name
| _Serial_ | Installed antenna serial number
| _Mark_ | Installed _Mark_ code
| _Height_ | Installed height | _metres_ above the mark
| _North_ | Installed offset north | _metres_
| _East_ | Installed offset east | _metres_
| _Azimuth_ | Installed azimuth [ _degrees_ clockwise from north
| _Start_ | Antenna installation start time
| _Stop_ | Antenna installation stop time

#### _RADOMES_ ####

A list of _radome_ installations associated with GNSS antenna installations.

| Field | Description | Units |
| --- | --- | --- |
| _Make_ | Installed radome make
| _Model_ | Installed radome model name
| _Serial_ | Installed radome serial number
| _Mark_ | Installed radome associated GNSS _Mark_
| _Start_ | Radome installation start time
| _Stop_ | Radome installation stop time

#### _METSENSORS_ ####

A list of _metsensor_ installations attached
to GNSS receivers

| Field | Description | Units |
| --- | --- | --- |
| _Make_ | Installed met sensor make
| _Model_ | Installed met sensor model name
| _Serial_ | Installed met sensor serial number
| _Mark_ | Installed met sensor associated GNSS _Mark_
| _IMS Comment_ | Header comments
| _Humidity_ | Humidity sensor accuracy | % rel H
| _Pressure_ | Pressure sensor accuracy | hPa
| _Temperature_ | Temperature sensor accuracy | deg C
| _Latitude_ | Installed met sensor latitude | degrees north
| _Longitude_ | Installed met sensor longitude | degrees east
| _Elevation_ | Installed met sensor elevation | metres
| _Datum_ | Installed met sensor datum used</br>to define the latitude, longitude, and elevation
| _Start_ | Installed met sensor installation start time
| _Stop_ | Installed met sensor installation stop time

#### _RECEIVERS_ ####

A list of GNSS _receiver_ installations.

| Field | Description | Units |
| --- | --- | --- |
| _Make_ | Deployed GNSS receiver make
| _Model_ | Deployed GNSS receiver model name
| _Serial_ | Deployed GNSS receiver serial number
| _Mark_ | Associated deployment GNSS _Mark_
| _Start_ | Receiver deployment start time
| _Stop_ | Receiver deployment stop time

#### _FIRMWARES_ ####

A list of GNSS receiver _firmware_ versions.

| Field | Description | Units |
| --- | --- | --- |
| _Make_ | GNSS receiver make
| _Model_ | GNSS receiver model name
| _Serial_ | GNSS receiver model serial number
| _Version_ | Installed receiver firmware version
| _Start_ | Receiver firmware start time
| _Stop_ | Receiver firmware stop time
| _Notes_ | Extra firmware notes

#### _SESSIONS_ ####
 
A list of _GNSS_ _Receiver_ &amp; _Antenna_ recording sessions

| Field | Description | Units |
| --- | --- | --- |
| _Mark_ | Session GNSS _Mark_
| _Operator_ | Deployed equipment operator information
| _Agency_ | Deployed equipment agency information
| _Model_ | Configuration model details
| _Satellite System_ | Configured receiver satellite settings
| _Interval_ | Configured receiver sampling interval|_seconds_
| _Elevation Mask_ | Configured receiver elevation mask| _degrees_ above the horizon
| _Header Comment_ | Configuration comments
| _Start_ | Session start time
| _Stop_ | Session stop time

### SIGNAL RECORDING ###

#### _SENSORS_ ####
 
A list of _sensor_ installations. The scale _factor_ and _bias_ can
be used for external adjustments, such as a pressure sensor being
used to measure water depth in salt water.

| Field | Description | Units |
| --- | --- | --- |
| _Make_  | Installed sensor make
| _Model_ | Installed sensor model name
| _Serial_ | Installed sensor serial number
| _Station_ | Installed recording _station_
| _Location_ | Installed sensor _site_ location
| _Azimuth_ | Installed sensor azimuth | _degrees_ clockwise from north
| _Dip_ | Installed sensor dip | _degrees_ down from horizontal
| _Depth_ | Installed sensor vertical offset | _metres_  positive downwards
| _North_ | Installed sensor offset north | _metres_
| _East_ | Installed sensor offset east | _metres_
| _Scale Factor_ | Optional installation specific</br>gain adjustment|defaults to _"1.0"_
| _Scale Bias_ |  Optional installation specific</br>level adjustment|defaults to _"0.0"_
| _Start_ | Sensor installation start time
| _Stop_ | Sensor installation stop time

#### _RECORDERS_ ####
 
A list of _recorder_ installations, these are considered to be
a combination of a sensor and a datalogger that are always
installed as a set.

| Field | Description | Units |
| --- | --- | --- |
| _Make_ | Installed recorder make
| _Sensor_ | Installed recorder sensor model name
| _Datalogger_ | Installed recorder datalogger model name
| _Serial_ | Installed recorder serial number
| _Station_ | Installed recording _station_
| _Location_ | Installed recording _site_ location
| _Azimuth_ | Installed recorder azimuth | _degrees_ clockwise from north
| _Dip_ | Installed recorder dip | _degrees_ down from horizontal
| _Depth_ | Installed recorder vertical offset | _metres_  positive downwards
| _Start_ | Installation start time
| _Stop_ | Installation stop time
 
#### _DATALOGGERS_ ####

A list of _datalogger_ deployments at a given _place_. Multiple
dataloggers at a _place_ can be distinguished by using an option
_role_ description. These are attached to the datalogger's 
associated _sensors_ via _connection_ records.

| Field | Description | Units |
| --- | --- | --- |
| _Make_ | Deployed datalogger make
| _Model_ | Deployed datalogger model name
| _Serial_ | Deployed datalogger serial number
| _Place_ | Deployed datalogger place
| _Role_ | Optional datalogger role at the place
| _Start_ | Deployment start time
| _Stop_ | Deployment stop time

#### _CONNECTIONS_ ####

A list of _datalogger_ connections, these are used to attach the sensors
at a given _site_ location to the dataloggers deployed at the associated
_place_. Multiple _dataloggers_ installed at the same place are distinguished
by an operational _role_, if required.

| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Recording _station_
| _Location_ | Sensor _site_ location
| _Place_ | Datalogger deployment _place_
| _Role_ | Datalogger deployment _role_
| _Start_ | Connection start time
| _Stop_ | Connection stop time

#### _STREAMS_ ####
 
A list of _datalogger_ sampling configurations for a given _station_ and recording _site_.

| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Recording _Station_|
| _Location_ | Recording locations _Site_|
| _Sampling Rate_ | Nominal stream sampling rate | samples per second (_Hz_)
| _Axial_ | Whether the stream is configured for</br>axial coordinates (_Z12_) or geographic (_ZNE_) |_"yes"_ or _"no"_
| _Reversed_ | Whether the recorded signal should</br>be reversed over the time window|_"yes"_ or _"no"_
| _Start_ | Stream start time|
| _Stop_ | Stream stop time|

### CAMERA ###

#### _CAMERAS_ ####

A list of _camera_ installations, these include values for:

| Field | Description | Units |
| --- | --- | --- |
|  _Make_ | Installed camera make |
|  _Model_ | Installed camera model name |
|  _Serial_ | Installed camera serial number |
|  _Mount_ | Camera _mount_ code |
|  _Dip_ | Installed camera dip | _degrees_ down from horizontal
|  _Azimuth_ | Installed camera azimuth | _degrees_ clockwise from north
|  _Height_ | Installed camera vertical offset | _metres_  positive upwards
|  _North_ | Installed camera offset north | _metres_
|  _East_ | Installed camera offset east | _metres_
|  _Start_ | Installed camera start time | 
|  _Stop_ | Installed camera stop time | 
|  _Notes_ | Extra installation information,</br>currently the photo caption.

### NOTES ###

Dates should be given as in _ISO 8601_ (i.e. `2016-09-18T02:24:26Z`), future dates should be given in the form: `9999-01-01T00:00:00Z`.
