# Seismic Sensor Orientation Conventions and Methods

The conventions and methods used by GeoNet to measure and record seismic sensor orientation are outlined in this document.

Seismic sensor orientations indicate the bearing of the sensor's north component. Seismic sensor orientations and the method used to determine them are presented in `install/sensors.csv` in the `Azimuth` and `Method` columns. 

Sensors that are not intended to be north-aligned will have `1` and `2` as the last character of their horizontal component channel code, and sensors that are intended to be north-aligned will be `N` and `E` as the last character of their horizontal component channel code - as per `docs/SEISMIC_CHANNEL_NAMING_CONVENTIONS.md`. It is important to mark this distinction: that all sensors have an orientation, but not all sensors are intended to be north-oriented.

## Orientation Methods

### Unknown

When an orientation is unknown, the `Azimuth` value will be `360` and the `Method` will be `unknown`.

### Onsite

When an orientation has been measured at the station, the `Azimuth` value will be this measurement and the `Method` will be `onsite`.

The most common onsite orientation methods are using a compass and an azimuith pointing system. A compass measurement involves taking a bearing of the sensor's north component using a compass. An azimuth pointing system (APS) projects a laser onto the sensor's north component marking and uses differential GPS positioning of two antenna at either end of its laser to measure the sensor orientation. Both of these approaches are considered reliable and which method has been used to measure sensor orientation is not noted.

### Offsite

When an orientation is calculated using the data from a station or - as in the case of strong motion sensors in building arrays - using geographic references, the `Azimuth` value is the calculated value and the `Method` is `offsite`.

There are many offsite orientation methods. Besides geographic referencing, the most common method is to use the polarisation of seismic waves from a known source to calculate the sensor orientation. Here both teleseismic earthquake phases and the Rayleigh waves of the secondary microseism recoverable using ambient noise techniques are used. Historically comparisons between surface and borehole seismic sensors have been used to calculate borehole sensor orientation, as has the approach of minimising moment tensor solution residuals by orienting the data from sensors with unknown orientations. To avoid presenting the full description of orientation techniques here, we simply note that offsite orientation methods are always the best possible and orientations are verified by the agreement of two or more methods where possible.

Often offsite orientation methods are used for borehole sensors and other sensors that cannot be oriented using onsite methods, i.e. those in building arrays.

### Method Accuracy

Uncertainty is not measured for orientation values. Unless the sensor orientation is unknown, the orientation uncertainty can be safely assumed to be +/- 5 degrees, if not much less.

## Legacy Sensor Orientations

In the GeoNet archive there are data from a long legacy of sensors. Unless the orientations of these sensors are stated as `360` or `unknown`, these are considered reliable.

## Reporting Sensor Orientation Inaccuracies

There is always the possibility that sensor orientations in delta are not accuruate. Through its operations, GeoNet will ocassionally discover inaccuracies and update its metadata accordingly. If you find an inaccuracy through your work, please reach out to info@geonet.org.nz.
