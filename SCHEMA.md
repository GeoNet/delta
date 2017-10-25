This is a notional schema as not all entries are represented by separate files.

There is an implied join over the _Start_ and _End_ times where suitable,
the overlapping time span should be used. For _Station_, _Mark_, _Mount_, _Site_ the
start and end times are attributes and should encompass all installed equipment.

# General

## Asset
Entry|Format|Unique
---|---|---
Make   | string| &#10003;
Model  | string| &#10003;
Serial | string| &#10003;
Number | string
Notes  | string

<div style="page-break-after: always;"></div>

## Network
Entry|Format|Unique
---|---|---
Code        | string| &#10003;
External    | string
Description | string
Restricted  | bool

<div style="page-break-after: always;"></div>

## Station
Entry|Format|Unique|Join
---|---|---|---
Network   | | &#10003;| **Network** **[** Code **]**
Code      | string|&#10003;
Name      | string
Latitude  | float64
Longitude | float64
Elevation | float64
Datum     | string
Start     | time
End       | time

<div style="page-break-after: always;"></div>

# GNSS

## Mark
Entry|Format|Unique|Join
---|---|---|---
Network   | |&#10003; | **Network** **[** Code **]**
Code      | string|&#10003;
Name      | string
Latitude  | float64
Longitude | float64
Elevation | float64
Datum     | string
IGS       | bool
Start     | time
End       | time

<div style="page-break-after: always;"></div>

## Monument
Entry|Format|Unique|Join
---|---|---|---
Mark               | |&#10003; | **Mark** **[** Code **]**
DomesNumber        | string
MarkType           | string
Type               | string
GroundRelationship | float64
FoundationType     | string
FoundationDepth    | float64
Bedrock            | string
Geology            | string
Start              | time
End                | time

<div style="page-break-after: always;"></div>

## Antenna
Entry|Format|Unique|Join
---|---|---|---
Asset | | &#10003; | **Asset** **[** Make, Model, Serial **]**

## Installed Antenna
Entry|Format|Unique|Join
---|---|---|---
Mark     | | &#10003; | **Mark** **[** Code **]**
Antenna  | | &#10003; | **Antenna** **[** Make, Model, Serial **]**
Vertical | float64
North    | float64
East     | float64
Azimuth  | float64
Start    | time |+
End      | time |+

<div style="page-break-after: always;"></div>

## Receiver
Entry|Format|Unique|Join
---|---|---|---
Asset | | | **Asset** **[** Make, Model, Serial **]**

## Installed Receiver
Entry|Format|Unique|Join
---|---|---|---
Mark     | | &#10003; | **Mark** **[** Code **]**
Receiver | | &#10003; | **Receiver** **[** Make, Model, Serial **]**
Start    | time |+
End      | time |+

<div style="page-break-after: always;"></div>

## Radome
Entry|Format|Unique|Join
---|---|---|---
Asset | | &#10003; | **Asset** **[** Make, Model, Serial **]**

## Installed Radome
Entry|Format|Unique|Join
---|---|---|---
Mark   | | &#10003; | **Mark** **[** Code **]**
Radome | | &#10003; | **Radome** **[** Make, Model, Serial **]**
Start  | time |+
End    | time |+

<div style="page-break-after: always;"></div>

## MetSensor
Entry|Format|Unique|Join
---|---|---|---
Asset | | &#10003; | **Asset** **[** Make, Model, Serial **]**

## Installed MetSensor
Entry|Format|Unique|Join
---|---|---|---
Mark       | | &#10003; | **Mark** **[** Code **]**
MetSensor  | | &#10003; | **MetSensor** **[** Make, Model, Serial **]**
Latitude   | float64
Longitude  | float64
Elevation  | float64
Datum      | string
IMSComment | string
Start      | time |+
End        | time |+

<div style="page-break-after: always;"></div>

## Firmware
Entry|Format|Unique|Join
---|---|---|---
Receiver | | &#10003; | **Receiver** **[** Make, Model, Serial **]**
Version  | string
Notes    | string
Start    | time |+
End      | time |+

<div style="page-break-after: always;"></div>

## Session
Entry|Format|Unique|Join
---|---|---|---
Mark            | | &#10003; | **Mark** **[** Code **]**
Operator        | string
Agency          | string
Model           | string
SatelliteSystem | string|&#10003;
Interval        | duration|&#10003;
ElevationMask   | float64
HeaderComment   | string
Format          | string
Start           | time|+
End             | time|+

<div style="page-break-after: always;"></div>

# Camera

## Mount
Entry|Format|Unique|Join
---|---|---|---
Network     | | &#10003; | **Network** **[** Code **]**
Code        | string|&#10003;
Name        | string
Latitude    | float64
Longitude   | float64
Elevation   | float64
Datum       | string
Description | string
Start       | time
End         | time

<div style="page-break-after: always;"></div>

## Camera
Entry|Format|Unique|Join
---|---|---|---
Asset | | | **Asset** **[** Make, Model, Serial **]**

## Installed Camera
Entry|Format|Unique|Join
---|---|---|---
Mount    | | &#10003; | **Mount** **[** Code **]**
Camera   | | &#10003; | **Camera** **[** Make, Model, Serial **]**
Dip      | float64
Azimuth  | float64
Vertical | float64
North    | float64
East     | float64
Notes    | string
Start    | time|+
End      | time|+

<div style="page-break-after: always;"></div>

# Tidal Modelling

## Gauge
Entry|Format|Unique|Join
---|---|---|---
Station   | | &#10003; | **Station** **[** Code **]**
Latitude  | float64
Longitude | float64
Elevation | float64
Datum     | string
Number    | string
TimeZone  | float64
Crex      | string

<div style="page-break-after: always;"></div>

## Constituent
Entry|Format|Unique|Join
---|---|---|---
Gauge     | | &#10003; | **Gauge** **[** Station **]**
Number    | int     |&#10003;
Name      | string  |&#10003;
Amplitude | float64
Lag       | float64

<div style="page-break-after: always;"></div>

# Signal Recording

## Site
Entry|Format|Unique|Join
---|---|---|---
Station   | | &#10003; | **Station** **[** Code **]**
Latitude  | float64
Longitude | float64
Elevation | float64
Datum     | string
Location  | string |&#10003;
Survey    | string
Start     | time
End       | time

<div style="page-break-after: always;"></div>

## Sensor
Entry|Format|Unique|Join
---|---|---|---
Asset | | | **Asset** **[** Make, Model, Serial **]**

## Installed Sensor
Entry|Format|Unique|Join
---|---|---|---
Site     | | &#10003; | **Site** **[** Station, Location **]**
Sensor   | | &#10003; | **Sensor** **[** Make, Model, Serial **]**
Dip      | float64
Azimuth  | float64
Vertical | float64
North    | float64
East     | float64
Factor   | float64
Bias     | float64
Start    | time |+
End      | time |+

<div style="page-break-after: always;"></div>

## Recorder
Entry|Format|Unique|Join
---|---|---|---
Asset | | | **Asset** **[** Make, Model, Serial **]**

## Installed Recorder
Entry|Format|Unique|Join
---|---|---|---
Site      | | &#10003; | **Site** **[** Station, Location **]**
Recorder  | | &#10003; | **Recorder** **[** Make, Model, Serial **]**
Dip       | float64
Azimuth   | float64
Vertical  | float64
North     | float64
East      | float64
Factor    | float64
Bias      | float64
Start     | time |+
End       | time |+

<div style="page-break-after: always;"></div>

## Datalogger
Entry|Format|Unique|Join
---|---|---|---
Asset | | | **Asset** **[** Make, Model, Serial **]**

## Deployed Datalogger
Entry|Format|Unique|Join
---|---|---|---
Datalogger  | | &#10003; | **Datalogger** **[** Make, Model, Serial **]**
Operation   | | &#10003; | **Operation** **[** Place, Role **]**
Start      | time |+
End        | time |+

<div style="page-break-after: always;"></div>

## Operation
Entry|Format|Unique
---|---|---
Place|string|&#10003;
Role|string|&#10003;

<div style="page-break-after: always;"></div>

## Connection
Entry|Format|Unique|Join
---|---|---|---
Site        | | &#10003; | **Site** **[** Station, Location **]**
Operation   | | &#10003; | **Operation** **[** Place, Role **]**
Start | time |+
End   | time |+

<div style="page-break-after: always;"></div>

## Stream
Entry|Format|Unique|Join
---|---|---|---
Site        | | &#10003; | **Site** **[** Station, Location **]**
SamplingRate | float64|&#10003;
Axial        | bool
Reversed     | bool
Start        | time |+
End          | time |+
