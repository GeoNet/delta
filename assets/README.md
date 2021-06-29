## ASSETS ##

_Lists the equipment used to collect the various observations._

Each file represents an equipment type and is listed with the following attributes:

* _Make_ - Equipment manufacturer or commercial brand name.
* _Model_ - Equipment model name or general identification label.
* _Serial_ - Equipment serial number, needs to be unique within its set of equipment make &amp; model.
* _Number_ - Equipment asset number, this is optional and is used primarily for cross-checking
* _Notes_ - Any extra details which may be relevant to the piece of equipment.

### FILES ###

* `antennas.csv` - GNSS Antennas.
* `cameras.csv` - Network Cameras used for volcano and building monitoring.
* `dataloggers.csv` - Dataloggers for recording analogue signals.
* `doases.csv` - DOAS Spectrometers used for volcano monitoring.
* `metsensors.csv` - Meteorological Sensors used primarily with GNSS systems.
* `radomes.csv` - GNSS Antenna Radomes.
* `receivers.csv` - GNSS Satellite Receivers
* `recorders.csv` - Combined datalogger and sensors used primarily for strong motion recording.
* `sensors.csv` - Analogue sensors which are usually attached to dataloggers.

### CHECKS ###

Pre-commit checks will be made on these files to ensure:
* No duplicated asset numbers - if given these need to be globally unique
* No duplicated serial numbers - these need to be unique within each equipment make &amp; model set.
