## CALIBRATION ##

_Sensor calibration information for scientific equipment installed at collection points._

### FOLDERS ###

* `antennas` - individual GNSS receiver antenna calibration files.
* `multigas` - gas detection limit information.
* `fluxgates` - individual Geomagnetic Fluxgate sensor calibration files.

### FILES ###

#### multigas.csv - gas calibration details

|Field|Description|Units|
|--|--|--|
|Station|Installed recording station||
|Location|Installed sensor site location||
|Gas|Calibration gas||
|Concentration|Of calibration gas|ppm (parts per million)|
|Frequency|How often calibration performed, D=daily, W=weekly, 4W=4 weekly, etc||
|Day Of Week|Calibration day, Monday, Tuesday, etc||
|Calibration Time|Time calibration performed|HH:MM:SS|
|Zero Time|Time observation should be zero|HH:MM:SS|
|Start Date|Of calibrations|YYYY-mm-ddTHH:MM:SSZ|
|End Date|Of calibrations|YYYY-mm-ddTHH:MM:SSZ|
