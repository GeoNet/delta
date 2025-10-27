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
* `polarities.csv` - site specific polarity settings which indicate when a site may have reversed polarity, or otherwise
* `gains.csv` - site specific settings applied to individual sensors that may impact overall sensitivities.
* `calibrations.csv` - Individual sensor sensitivity values that can be used rather than default values.
* `components.csv` - Individual sensor elements including measurement position and responses.
* `channels.csv` - Individual datalogger recording elements including digitiser position, sampling rate, and responses.
* [`preamps.csv`](#preamps) - site specific settings applied to individual datalogger pre-amplification that may impact overall sensitivities.
* [`telemetries.csv`](#telemetries) - site specific settings applied to datalogger and sensor connections that may use analogue telemetry.
* [`timings.csv`](#timings) - site specific settings to indicate time corrections that may need to be applied to archived raw data.

* `cameras.csv` - Installed field cameras.
* `doases.csv` - Installed field DOAS (Differential Optical Absorption Spectrometer) equipment.

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
 
A list of _GNSS_ _Receiver_ and _Antenna_ recording sessions

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
For dataloggers that have different response characteristics depending on
what channel number (or pin) is used, an offset can be given for the sensor
channel start.

| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Recording _station_
| _Location_ | Sensor _site_ location
| _Place_ | Datalogger deployment _place_
| _Role_ | Datalogger deployment _role_
| _Number_ | Initial datalogger pin, or channel, offset number used for the sensor
| _Start_ | Connection start time
| _Stop_ | Connection stop time

#### _STREAMS_ ####
 
A list of _datalogger_ sampling configurations for a given _station_ and recording _site_.

| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Recording _Station_|
| _Location_ | Recording locations _Site_|
| _Band_ | Channel _Band_ code|
| _Source_ | Channel _Source_ code|
| _Sampling Rate_ | Nominal stream sampling rate | samples per second (_Hz_)
| _Axial_ | Whether the stream is configured for</br>axial coordinates (_Z12_) or geographic (_ZNE_) |_"yes"_ or _"no"_
| _Triggered_ | Whether the stream represents</br>triggered recordings|_"yes"_ or _"no"_
| _Start_ | Stream start time|
| _Stop_ | Stream stop time|

The band and source codes are representatives of the FDSN channel naming convention as found at:

[FDSN Source Identifiers: Channel codes](http://docs.fdsn.org/projects/source-identifiers/en/v1.0/channel-codes.html)

#### _POLARITIES_ ####

Site specific times when the recorded values may have a reversed, or otherwise, polarity.
This is often difficult to track with the known instrument responses and installation details as it can
often be caused by in-field wiring or cabling changes from the expected standard.

There is also the possibility that there may be conflicting information or studies, where this is
the case the _Primary_ field can be used to indicate which one should be used for downstream processing.

To reduce the number valid _Method_ entries these are defined as a set of "known" values, and an "unknown" one.
Where possibly the code should be updated to add any standard polarity detection methods, there is also
room to provide a reference citation for any associated studies.


| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Recording _Station_|
| _Location_ | Recording location _Site_|
| _Sublocation_ | Recording location _Sublocation_|
| _Subsource_ | Recording location _Subsource_|
| _Primary_ | Whether the entry takes precedence| _"yes"_, _"no"_, or blank.
| _Reversed_ | Whether the stream in the time window should be considered reversed or not| _"yes"_ or _"no"_
| _Method_ | How this information was obtained| _"study"_, _"compass"_, or _"unknown"_
| _Citation_ | A reference citation for the method if appropriate| _"key"_
| _Start_ | Time window start time|
| _Stop_ | Time window stop time|

Notes:

- For the _Subsource_, this can either be individual entries (e.g., "Z", "N"), or multiple entries ("ZNE"), or it can be empty, which will be interpreted as all subsource values.
(The subsource is the last character in the standard SEED channel naming convention, e.g. EHZ).

- If a citation is given the actual reference information should be given in the _citations.csv_ file.

- An empty, or blank, _Method_ is treated as _unknown_.

Example:

    Station,Location,Sublocation,Subsource,Primary,Reversed,Method,Citation,Start Date,End Date
    WEL,10,,Z,,true,compass,,2022-07-28T01:59:00Z,9999-01-01T00:00:00Z

#### _GAINS_ ####
 
Site specific gain settings applied to correct for local conditions. A list of installation times where gains need to be applied to datalogger or sensor settings.
For the scale factor and bias either a value can be given directly or an expression can be used if that makes it clearer where the number has come from.

| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Datalogger recording _Station_|
| _Location_ | Recording sensor site _Location_ |
| _Sublocation_ | additional location identifier for multi-parametric sensors installations, if applicable |
| _Subsource_ | The sensor channel(s), as defined in the response configuration, which requires a gain adjustment, multiple subsource channels can be joined (e.g _"Z"_ or _"ZNE"_).
| _Scale Factor_ | Scale, or gain factor, that the input signal is multiplied by prior to digitisation, or for polynomial responses it is the factor used to convert Volts into the signal units. If this field is empty, it should be assumed to have a value of __1.0__ which in theory should have no impact.
| _Scale Bias_ | An offset value that needs to be added to the signal prior to digitisation or to raw digital data. The offset indicates a polynomial response is expected, if this field is blank it is assumed that the value is __0.0__.
| _Absolute Bias_ | An offset value that needs to be added to the signal after the scale factors have been applied to the polynomial response, if this field is blank it is assumed that the value is __0.0__.
| _Start_ | Gain start time|
| _Stop_ | Gain stop time|

For a second order polynomial response, the output is expected to be `Y = a * X + b` where `X` is normally the input voltage, and Y the corrected signal.
The terms `a` and `b` are the factor and bias respectively. The gain adjustments (`a'`, `b'`, `c'`, the scale factor, scale bias, and absolute bias respectively)
update this via `Y = (a * a') * X + (b * a') + (a * b') + c'`

#### _CALIBRATIONS_ ####
 
Sensor specific calibrations that may impact overall sensitivity. A list of installation times where calibrated values of the _Sensor_ sensitivity are known and can be used to override 
the default _Model_ sensitivities.
For the component, sensitivity, and frequency either a value can be given directly or an expression can be used if that is more readable.

| Field | Description | Units |
| --- | --- | --- |
| _Make_ | Sensor make
| _Model_ | Sensor model name
| _Serial_ | Sensor serial number
| _Number_ | The sensor component or datalogger digitiser channel, as defined in the response configuration or elsewhere, which overrides the default values, a blank value is interpreted as the first sensor component, or __"pin"__ zero.
| _Scale Factor_ | Sensitivity, or scale factor, that the input signal is generally multiplied by to convert to Volts, or for polynomial responses the value used to convert Volts into the signal units. A blank value is expected to be read as __1.0__, an explicit value of zero is required to be entered if intended.
| _Scale Bias_ | A scale bias factor, for polynomial responses that is multiplied to any _Gain_ bias values before adding to the converted volts to give the signal values. If this field is blank it should be assumed that the value is __1.0__.
| _Scale Absolute_ | An offset, or bias, for polynomial responses that is added to the converted volts to give the signal values. If this field is blank it should be assumed that the value is __0.0__.
| _Frequency_ | Frequency at which the calibration value is correct for if appropriate.
| _Start_ | Calibration start time|
| _Stop_ | Calibration stop time|

For a second order polynomial response, the output is expected to be `Y = a * X + b` where `X` is normally the input voltage, and Y the corrected signal. The terms `a` and `b` are the scale factor and the absolute bias respectively. The scale bias `c` can be used with gain biases, e.g. for a gain adjustment of the equivalent (`a'`, `b'`, `c'`) values, then the polynomial adjustment will follow: `Y = a' * a * X + b * a' + a * b' * c + c'`.

#### _COMPONENTS_ ####

Sensor model component descriptions. The type is generally of the form "Accelerometer, Short Period Seismometer" etc.
The number represents the order of sensor components, this generally maps to the sensor cable and how it is connected
into the datalogger.
Subsource is the general term used for labelling the sensor component and is usually the last character in the SEED channel convention.
Dip and Azimuth are used to indicate the relative position of the sensor component within the sensor package and will be used with the
overall sensor installation values to provide component dips and azimuths.

For derived streams, such as a simple gain or unit conversion, can be indicated by providing an input sampling rate. This is matched by
the equivalent _Stream_ and allows for the response to be generated with the provided reference response only.

| Field           | Description |
| --------------- | ----------- |
| _Make_          | Sensor make
| _Model_         | Sensor model name
| _Type_          | Sensor type
| _Number_        | Sensor component offset
| _Source_        | Sensor source as used for the matching streams 
| _Subsource_     | Sensor component label
| _Dip_           | Internal dip of the component relative to whole sensor
| _Azimuth_       | Internal azimuth of the component relative to whole sensor
| _Types_         | A shorthand reference to the SEED type labels
| _Sampling Rate_ | An input sampling rate which can be used to indicate a _derived_ stream
| _Response_      | A reference to the nominal StationXML response

#### _CHANNELS_ ####

The individual channels configured for a given datalogger model, these include the channel numbers and sampling rates.
The channel number is an offset into the digitiser or digitisers and are used to match the connected sensor component
and the expected response. Some digitisers have different nominal responses for different groups of digitiser channels.

| Field           | Description | 
| --------------- | ----------- |
| _Make_          | Datalogger make
| _Model_         | Datalogger model name
| _Type_          | Datalogger type
| _Number_        | Datalogger channel offset, an empty value will map to zero
| _Sampling Rate_ | Configured Channel sampling rate
| _Response_      | A reference to the nominal StationXML response 

#### _PREAMPS_ ####

| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Datalogger recording _Station_|
| _Location_ | Recording sensor site _Location_ |
| _Subsource_ | The sensor channel _Subsource_ which has the preamp configured (e.g _"Z"_). An empty value indicates all channels have this setting for the provided _Location_.
| _Scale Factor_ | The datalogger pre-amp scale factor used for this time span. These tend to be integer steps and may be referenced as **gain** settings.
| _Start_ | Preamp start time|
| _Stop_ | Preamp stop time|

#### _TELEMETRIES_ ####

Sometimes the datalogger and the sensor are not at the same location. Usually this means there is some form of analogue link between the two, either a dedicated
telephone line, or an FM radio link. This table allows this to be documented, and provides a mechanism to adjust the signal gains if known.

| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Datalogger recording _Station_|
| _Location_ | Recording sensor site _Location_ |
| _Scale Factor_ | The telemetry gain factor for the analogue link, this represents the amplification of the signal if appropriate, an empty value is assumed to be 1.0
| _Start_ | Telemetry start time|
| _Stop_ | Telemetry stop time|

### CAMERA ###

#### _CAMERAS_ ####

A list of _camera_ installations, these include values for:

| Field | Description | Units |
| --- | --- | --- |
|  _Make_ | Installed camera make |
|  _Model_ | Installed camera model name |
|  _Serial_ | Installed camera serial number |
|  _Mount_ | Camera _mount_ code |
|  _View_ | Camera _view_ code |
|  _Dip_ | Installed camera dip | _degrees_ down from horizontal
|  _Azimuth_ | Installed camera azimuth | _degrees_ clockwise from north
|  _Height_ | Installed camera vertical offset | _metres_  positive upwards
|  _North_ | Installed camera offset north | _metres_
|  _East_ | Installed camera offset east | _metres_
|  _Start_ | Installed camera start time | 
|  _Stop_ | Installed camera stop time | 
|  _Notes_ | Extra installation information,</br>currently the photo caption.

### DOAS ###

#### _DOASES_ ####

A list of _doas_ installations, these include values for:

| Field | Description | Units |
| --- | --- | --- |
|  _Make_ | Installed DOAS make |
|  _Model_ | Installed DOAS model name |
|  _Serial_ | Installed DOAS serial number |
|  _Mount_ | DOAS _mount_ code |
|  _View_ | DOAS _view_ code |
|  _Dip_ | Installed DOAS dip | _degrees_ down from horizontal
|  _Azimuth_ | Installed DOAS azimuth | _degrees_ clockwise from north
|  _Height_ | Installed DOAS vertical offset | _metres_  positive upwards
|  _North_ | Installed DOAS offset north | _metres_
|  _East_ | Installed DOAS offset east | _metres_
|  _Start_ | Installed DOAS start time | 
|  _Stop_ | Installed DOAS stop time | 

#### _TIMINGS_ ####

There are some data sources which may have timing issues and the type of recording format means
that the archived raw data will include the bad timing and a correction will need to be applied.

Likewise, some sites may have a timing issue that cannot easily be corrected, this table allows
this correction to be stored for use with processing systems or otherwise.

The correction should be "added" to the sensor sample times, for when the clock is fast this value
should be negative. The format can follow a standard _go_ convention: 

    A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix,
    such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d", "w", "y".


| Field | Description | Units |
| --- | --- | --- |
| _Station_ | Recording _Station_|
| _Location_ | Recording location _Site_|
| _Correction_ | Recorded data required time correction |
| _Start_ | Time window start time|
| _Stop_ | Time window stop time|

Example:

    Station,Location,Correction,Start Date,End Date
    NZB,42,-24h,2023-12-07T14:15:00Z,9999-01-01T00:00:00Z

### NOTES ###

Dates should be given as in _ISO 8601_ (i.e. `2016-09-18T02:24:26Z`), future dates should be given in the form: `9999-01-01T00:00:00Z`.
