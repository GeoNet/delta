## NETWORK ##

_Lists the geographical and physical characteristics of collection points._

### FILES ###

* `marks.csv` - GNSS Observation Points.
* `monuments.csv` - GNSS Observation Monuments.
* `mounts.csv` - Network Camera Mount Points.
* `networks.csv` - Grouping of GNSS Marks & Recording Stations.
* `sites.csv` - Specific Observation Points at a Specific Recording Station.
* `stations.csv` - Location for a Group of Recording Sites.
* `views.csv` - Network Camera View codes.

#### _MARKS_ ####

* _Mark_ - Code used to uniquely identify GNSS _Mark_.
* _Network_ - Code used to group marks together by project or operator.
* _Igs_ - Whether the _Mark_ is used by the IGS service, *yes* or *no*.
* _Name_ - Used to describe the general geographical location of the _Mark_.
* _Latitude_ - Geographical latitude of the _Mark_ for the given _Datum_.
* _Longitude_ - Geographical longitude of the _Mark_ for the given _Datum_.
* _Elevation_ - Height in meters of the _Mark_ for the given _Datum_.
* _Datum_ - Geographical reference system used for the latitude, longitude & elevation.
* _Start Date_ - General date and time at which the _Mark_ was operational.
* _End Date_ - General date and time at which the _Mark_ was no longer operational.

#### _MONUMENTS_ ####

* _Mark_ - Code used to uniquely identify GNSS _Mark_.
* _Domes Number_ - Globally unique identifier used for a continuous GNSS _Mark_ assigned by IERS.
* _Type_ - Type of Antenna mount construction, at the top of a geodetic monument.
* _Mark Type_ - Type of Monument construction, see below for valid monument types.
* _Ground Relationship_ - Height of the geodetic monument from the ground surface.
* _Foundation Type_ - Type of monument foundation.
* _Foundation Depth_ - Depth of monument foundation.
* _Start Date_ - General date and time at which the _Monument_ was operational.
* _End Date_ - General date and time at which the _Monument_ was operational.
* _Bedrock_ - Type of bedrock.
* _Geology_ - Description of geologic conditions of the ground surface.

Valid Monument Types are:

* `Shallow Rod / Braced Antenna Mount`
* `Wyatt/Agnew Drilled-Braced`
* `Pillar`
* `Steel Mast`
* `Unknown`

#### _MOUNTS_ ####

* _Mount_ - Code used to uniquely identify Camera _Mount_.
* _Network_ - Code used to group marks together by project or operator.
* _Name_ - Used to describe the general geographical location of the _Mount_.
* _Latitude_ - Geographical latitude of the _Mount_ for the given _Datum_.
* _Longitude_ - Geographical longitude of the _Mount_ for the given _Datum_.
* _Elevation_ - Height in meters of the _Mount_ for the given _Datum_.
* _Datum_ - Geographical reference system used for the latitude, longitude & elevation.
* _Description_ - Caption used for the camera _Mount_.
* _Start Date_ - General date and time at which the _Mount_ was operational.
* _End Date_ - General date and time at which the _Mount_ was no longer operational.

#### _NETWORKS_ ####

* _Network_ - Code used to group GNSS _Marks_ and Recording _Stations_ together by project or operator.
* _External_ - Alternative code used to externally represent this _Network_.
* _Description_ - Information about the _Network_.
* _Restricted_ - Whether the _Network_ has restrictions, a Boolean value [`true` or `false`].

#### _SITES_ ####

* _Station_ - Code used to uniquely identify Recording _Station_.
* _Location_ - Code used to uniquely identify the _Site_ at the  Recording _Stations_.
* _Latitude_ - Geographical latitude of the _Site_ for the given _Datum_.
* _Longitude_ - Geographical longitude of the _Site_ for the given _Datum_.
* _Elevation_ - Height in meters of the _Site_ above the free surface for the given _Datum_.
* _Depth_ - Depth of water in meters above the _Site_ if installed underwater.
* _Datum_ - Geographical reference system used for the latitude, longitude & elevation.
* _Start Date_ - General date and time at which the _Site_ was operational.
* _End Date_ - General date and time at which the _Site_ was no longer operational.

#### _STATIONS_ ####

* _Station_ - Code used to uniquely identify Recording _Station_.
* _Network_ - Code used to group Recording _Stations_ together by project or operator.
* _Name_ - Used to describe the general geographical location of the _Station_.
* _Latitude_ - Geographical latitude of the _Station_ for the given _Datum_.
* _Longitude_ - Geographical longitude of the _Station_ for the given _Datum_.
* _Elevation_ - Height in meters of the _Station_ above the free surface for the given _Datum_.
* _Depth_ - Depth of water in meters above the _Station_ if installed underwater.
* _Datum_ - Geographical reference system used for the latitude, longitude & elevation.
* _Start Date_ - General date and time at which the _Station_ was operational.
* _End Date_ - General date and time at which the _Station_ was no longer operational.

#### _VIEWS_ ####

* _Mount_ - Code used to uniquely identify a Camera _Mount_.
* _View_ - Code used to uniquely identify a Camera _View_.
* _Label_ - Feature the Camera is pointing to.
* _Azimuth_ - Camera _View_ azimuth, measured in _degrees_ clockwise from north.
* _Method_ - 
* _Dip_ - Camera _dip_ angle, measured in _degrees_ down from the horizontal.
* _Description_ - Detailed description of the feature and Camera _View_.
* _Start Date_ - General date and time at which the _Station_ was operational.
* _End Date_ - General date and time at which the _Station_ was no longer operational.


### CHECKS ###

Pre-commit checks will be made on these files to ensure:

* No duplicated marks - these need to be globally unique
* No duplicated monuments - these will have a new _Mark_ code if rebuilt
* That all monument types are valid
* No duplicated networks - these need to be globally unique
* No duplicated stations - these need to be globally unique
* No duplicated sites - these need to be unique at each station

### NOTES ###

Dates should be given as in _ISO 8601_ (i.e. `2016-09-18T02:24:26Z`), future dates should be given in the form: `9999-01-01T00:00:00Z`.

