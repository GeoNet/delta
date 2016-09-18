## NETWORK ##

Lists the geographical and phsyical characteristics of collection points.

### FILES ###

* `marks.csv` - GNSS Observation Points.
* `monuments.csv` - GNSS Observation Monuments.
* `mounts.csv` - Network Camera Mount Points.
* `networks.csv` - Grouping of GNSS Marks & Recording Stations.
* `sites.csv` - Specific Observation Points at a specific Recording Station.
* `stations.csv` - Location for a group of Recording Sites.

#### MARK ####

 * _Mark_ - Code used to uniquely identify GNSS _Mark_.
 * _Network_ - Code used to group marks together by project or operator.
 * _Name_ - Used to describe the general geographical location of the _Mark_.
 * _Latitude_ - Geographical latitude of the _Mark_ for the given _Datum_.
 * _Longitude_ - Geographical longitude of the _Mark_ for the given _Datum_.
 * _Elevation_ - Height in meters of the _Mark_ for the given _Datum_.
 * _Datum_ - Geographical reference system used for the latitude, longitude & elevation.
 * _Start Date_ - General date and time at which the _Mark_ was operational.
 * _End Date_ - General date and time at which the _Mark_ was no longer operational.

#### MONUMENTS ####

 * _Mark_ - Code used to uniquely identify GNSS _Mark_.
 * _Domes Number_ -
 * _Type_ - Monument construction, see below for valid types.
 * _Mark Type_ -
 * _Ground Relationship_ -
 * _Foundation Type_ -
 * _Foundation Depth_ -
 * _Start Date_ - General date and time at which the _Monument_ was operational.
 * _End Date_ - General date and time at which the _Monument_ was operational.
 * _Bedrock_ -
 * _Geology_ -

Valid Monument Types are:

 * `Shallow Rod / Braced Antenna Mount`
 * `Wyatt/Agnew Drilled-Braced`
 * `Pillar`
 * `Steel Mast`
 * `Unknown`

#### MOUNTS ####

 * _Mount__ - Code used to uniquely identify Camera _Mount_.
 * _Network_ - Code used to group marks together by project or operator.
 * _Name_ - Used to describe the general geographical location of the _Mount_.
 * _Latitude_ - Geographical latitude of the _Mount_ for the given _Datum_.
 * _Longitude_ - Geographical longitude of the _Mount_ for the given _Datum_.
 * _Elevation_ - Height in meters of the _Mount_ for the given _Datum_.
 * _Datum_ - Geographical reference system used for the latitude, longitude & elevation.
 * _Description_ - Caption used for the camera _Mount_.
 * _Start Date_ - General date and time at which the _Mount_ was operational.
 * _End Date_ - General date and time at which the _Mount_ was no longer operational.

#### NETWORKS ####

Network,External,Description,Restricted

 * _Network_ - Code used to group GNSS _Marks_ and Recording _Stations_ together by project or operator.
 * _External_ - Alternative code used to externally represent this _Network_.
 * _Description_ - Information about the _Network_.
 * _Restricted_ - Whether the _Network_ has restrictions, a boolean value [`true` or `false`].

#### SITES ####

#### STATIONS ####

### CHECKS ###

Pre-commit checks will be made on these files to ensure:
* No duplicated marks - these need to be globally unique
* No duplicated monuments - these will have a new _Mark_ code if rebuilt
* That all monument types are valid

### NOTES ###

Dates should be given as in _ISO 8601_ (i.e. `2016-09-18T02:24:26Z`), future dates should be given in the form: `9999-01-01T00:00:00Z`.

