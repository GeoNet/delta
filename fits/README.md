## FITS ##

_For routine ad hoc measurements or processed data observations._

### FILES ###

* `locations.csv` - Measurement Location Sites.
* `measurements.csv` - Measurement Reading Descriptions.

#### _LOCATIONS_ ####

* _Code_ - Code used to uniquely identify _Measurement_ location or site.
* _Description_ - Used to describe the general _Measurement_ feature.
* _Latitude_ - Geographical latitude of the _Location_ for the given _Datum_.
* _Longitude_ - Geographical longitude of the _Location_ for the given _Datum_.
* _Elevation_ - Height in meters of the _Location_ for the given _Datum_.
* _Datum_ - Geographical reference system used for the latitude, longitude & elevation.
* _Height_ - A vertical adjustment to the _Elevation_ to where the _Measurement_ is taken from.

#### _MEASUREMENTS_ ####

* _Code_ - Code used to associated a _Measurement_ with a _Location_ site.
* _Name_ - An overall term for the  _Measurement_ being made.
* _Sensor_ - A specific measurement label for distinguishing multiple _Measurement_ readings at a _Location_.
* _Type_ - A reference to the method being used for making the _Measurement_.
* _Unit_ - A reference to the _Measurement_ unit type.
* _Description_ - A general description of the _Measurement_.
