## ENVIRONMENT ##

_Lists the geological and physical environment details of collection points._

### FILES ###

* `classes.csv` - Site Class Descriptions.

* `constituents.csv` - Constituents Descriptions.

* `features.csv` - Sensor Installation Descriptions.

* `gauges.csv` - Gauges Descriptions.

* `placenames.csv` - Placename Descriptions.

* `gts.csv` - WMO Numbers and Bulletin Headers.

* `visibility.csv` - Sky View Descriptions.

#### _CLASSES_ ####

Site classes as used for earthquake engineering applications and studies, for more details please see:

> __Kaiser, A., Van Houtte, C., Perrin, N., Wotherspoon, L., & McVerry, G. (2017).__
> Site characterisation of GeoNet stations for the New Zealand Strong Motion Database.
> Bulletin of the New Zealand Society for Earthquake Engineering, 50(1), 39–49. https://doi.org/10.5459/bnzsee.50.1.39-49

| Field | Description |
| --- | --- |
| _Station_        | Station _Code_ used to identify where the site class has been described
| _Site Class_     | Site class as defined by NZS1170, ranges between *A* and *E*
| _Vs30_           | The time-averaged shear-wave velocity to 30m depth in meters per second
| _Vs30 Quality_   | A quality estimate of the _Vs30_ estimate, ranges from *Q1* to *Q3*
| _Tsite_          | The low-strain fundamental site period in seconds, a range can be given as a prefix (e.g. &lt; or &gt;)
| _Tsite Method_   | An indication of the method used to estimate _Tsite_
| _Tsite Quality_  | A quality estimate of the _Tsite_ estimate, ranges from *Q1* to *Q3*
| _Basement Depth_ | Estimated depth to the basement layer in meters, this may also be estimated as the depth to 1000 meters per second.
| _Depth Quality_  | A quality estimate of the depth to basement, ranges from *Q1* to *Q3*
| _Link_           | A reference link to another site
| _Citations_      | A semicolon separated list of reference citations for the site class estimates
| _Notes_          | Extra notes which should be avoided in lieu of a better machine readable field

For reference the NZS1170.5 site classes are roughly

> | Class | Description |
> | --- | --- |
> | A   | __Strong Rock__
> | B   | __Rock__
> | C   | __Shallow Soil__
> | D   | __Deep or Soft Soil__
> | E   | __Very Soft Soil__

Quality estimates are roughly

> | Quality | Description |
> | --- | --- |
> | Q1  | &lt; 10 %
> | Q2  | 10 - 20 %
> | Q3  | &gt; 20 %

#### _CONSTITUENTS_ ####

| Field | Description |
| --- | --- |
| _Gauge_ | Code used to uniquely identify Tide Gauge _Station_
| _Number_ | Constituent number, used mainly for display and sorting
| _Constituent_ | Standard Constituent Name
| _Amplitude_ | Analysis Amplitude, in cm
| _Lag_ | Analysis Phase Lag, in degrees
| _Start Date_ | General date and time at which the _Constituent_ description was valid.
| _End Date_ | General date and time at which the _Constituent_ description was no longer valid.

#### _FEATURES_ ####

| Field | Description |
| --- | --- |
| _Station_ | Code used to uniquely identify Recording _Station_.
| _Location_ | Code used to uniquely identify the _Site_ at the  Recording _Station_.
| _Sublocation_ | Code used to uniquely identify the _Site_ Sublocation if applicable.
| _Property_ | Property being measured.
| _Description_ | A helpful description of the physical site location.
| _Aspect_ | Additional description of the physical site location, if applicable.
| _Start Date_ | General date and time at which the _Site_ description was valid.
| _End Date_ | General date and time at which the _Site_ description was no longer valid.

#### _GAUGES_ ####

| Field | Description |
| --- | --- |
| _Gauge_ | Code used to uniquely identify Tide Gauge _Stations_
| _Network_ | Code used to group Tide Gauge _Stations_ together by project or operator
| _LINZ Number_ | Code used by _LINZ_ to identify the Tide Gauge _Station_
| _Analysis Time Zone_ | Time-zone offset used in the Tidal Constituent Analysis
| _Analysis Latitude_ | Latitude used in the Tidal Constituent Analysis, usually positive
| _Analysis Longitude_ | Longitude used in the Tidal Constituent Analysis
| _Crex Tag_ | Tide gauge Crex format location
| _Start Date_ | General date and time at which the _Gauge_ description was valid.
| _End Date_ | General date and time at which the _Gauge_ description was no longer valid.

#### _GTS_ ####

| Field | Description |
| --- | --- |
| _Gauge_ | Code used to uniquely identify Tide Gauge _Stations_
| _Network_ | Code used to group Tide Gauge _Stations_ together by project or operator
| _Analysis Latitude_ | Latitude used in the Tidal Constituent Analysis, usually positive
| _Analysis Longitude_ | Longitude used in the Tidal Constituent Analysis
| _Bulletin Header_ | GTS Bulletin Header
| _WMO Number_ | GTS WMO Number
| _Start Date_ | General date and time at which the _Gauge_ description was valid.
| _End Date_ | General date and time at which the _Gauge_ description was no longer valid.

#### _PLACENAMES_ ####

| Field | Description |
| --- | --- |
| _Name_ | The name of the place to use for distance measurements.
| _Latitude_ | The latitude of the place.
| _Longitude_ | The longitude of the place, this should be within -180 to 180 degrees.
| _Level_ | The level is used to discard small places when finding the closest place at large distances.

The level should be from 0 to 3, with 3 used for the smallest places.

#### _VISIBILITY_ ####

| Field | Description |
| --- | --- |
| _Code_ | Code to uniquely identify GNSS _Mark_ (or recording _Station_)
| _Sky Visibility_ | Free form description of the site sky visibility and obstructions
| _Start Date_ | General date and time at which the visibility was accurate
| _End Date_ | General date and time at which the visibility was no longer accurate

### CHECKS ###

Pre-commit checks will be made on these files to ensure:
* No duplicated gauges - these need to be globally unique

### NOTES ###

Dates should be given as in _ISO 8601_ (i.e. `2016-09-18T02:24:26Z`), future dates should be given in the form: `9999-01-01T00:00:00Z`.

