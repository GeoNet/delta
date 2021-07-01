# Seismic Sensor Orientation Methods

The methods used by GeoNet to determine seismic sensor orientation are outlined in this document.

Seismic sensor orientations and a description of the method used to determine these, when known, are captured in `install/sensors.csv` in the `Azimuth` and `Orientation Method` columns respectively.

Seismic sensor orientations are the bearing of the sensor's north component. Sensors that are not aligned at a 000 bearing (north component pointing north) will have their horizontal components labeled with `1` and `2` as the last character of the channel code, as per `docs/SEISMIC_CHANNEL_NAMING_CONVENTIONS.md`.

If a sensor is "north-aligned", i.e. orientation of 000, then no orientation method is noted in the metadata. Similarly, if a sensor is not north-aligned but GeoNet is unaware of the misorientation, the metadata will record correct sensor orientation. This is an issue currently being addressed.  

## Orientation Methods

### Compass

The compass method involves taking a bearing of the sensor's north component using a compass.

In *delta*, the compass method is noted as `compass` in the `Orientation Method` column of the `install/sensors.csv` file.

### Azimuth Pointing System

The Azimuth Pointing System (APS) method uses a differential GPS system with two GPS antennae at a short horizontal offset (~1 m) at either ends of an arm with one end afixed to a tripod. The APS instrument projects a laser downward from below its outer antenna which can be aligned with the sensor's north component marking. The relative position of the outer antenna compared to the inner antenna on the APS arm gives the sensor orientation.

In *delta*, the APS method is noted as `APS` in the `Orientation Method` column of the `install/sensors.csv` file.

### Signal Coherency

The signal coherency method uses a reference seismometer with a known orientation and the coherency of a given signal between the two seismic sensors to determine the sensor orientation. Signal coherency methods simulate a rotation of the sensor and find the rotation angle for which the signal recorded on both sensors is most similar. The sensor orientation is then determined from this rotation angle.

Signal coherency methods work best when the distance between sensors is small and the signals used are strongly coherent between the two sensors. Commonly the signal coherency method is used for orientating borehole sensors as neither of the other methods are applicable.

In *delta*, the signal coherency method is noted as either `local coherency` or `remote coherency` in the `Orientation Method` column of the `install/sensors.csv` file, depending on whether an oriented surface sensor was used (in the case of borehole sensors), or if an oriented remote sensor was used.
