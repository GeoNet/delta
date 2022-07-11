## ENVIRONMENT ##

_Lists the geological and physical environment details of collection points._

### FILES ###

* `constituents.csv` - Constituents Descriptions.

* `features.csv` - Sensor Installation Descriptions.

* `gauges.csv` - Gauges Descriptions.

* `placenames.csv` - Placename Descriptions.

* `visibility.csv` - Sky View Descriptions.

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
| _SubLocation_ | Code used to uniquely identify the _Site_ SubLocation if applicable.
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

