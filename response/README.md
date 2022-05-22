## RESPONSE ##

_Lists the equipment details useful for building response information._

### FILES ###

* `labels.csv` - Describes datalogger stream channel codes for given sensors.

#### _LABELS_ ####

| Field | Description |
| --- | --- |
| _Type_ | The _Sensor_ type of the recorded stream.
| _Sampling Rate_ | The _Stream_ sampling rate, or sampling period if negative.
| _Azimuth_ | The internal _Sensor_ component azimuth relative to North.
| _Dip_ | The internal _Sensor_ dip angle relative to the horizontal, with positive down.
| _Code_ | The code used for the __StationXML__ Channel for the _Stream_ and _Sensor_ pair.
| _Flags_ | Any specific _Sensor_ meta-data flags as expected by __StationXML__.

The _Flags_ should be single characters arranged alphabetically into a single word.

### CHECKS ###

Pre-commit checks will be made on these files to ensure:
* No duplicated type, sampling rate, azimuth, and dip - these need to be globally unique
* The code needs to be uppercase and have a length of three characters.
* The flags characters need to be sorted into a single word sequence.
