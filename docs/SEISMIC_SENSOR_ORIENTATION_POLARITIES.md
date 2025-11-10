# Seismic Sensor Orientation and Polarity Conventions and Methods

The conventions and methods used by GeoNet to measure and represent seismic sensor orientations and polarities are outlined in this document.

Data users encounter orientations and polarities through the component orientations in StationXML. This document describes the underlying conventions and processes which define these values.

## Orientation Conventions and Measurement Methods

Sensor orientations are contained in `delta.git/install/sensors.csv` under the `Azimuth` column. All sensor installations have an orientation entry, although it is not always defined.

Seismic sensor orientations indicate the bearing of the sensor's north component and assume a left-hand convention for the orientation of the vertical and east components with respect to the north component (i.e., E component is at 90 degrees clockwise of N in the N-E plane; Z component is 90 degrees anticlockwise in N-Z plane).

Sensors that are not *intended* to be north-aligned will have `1` and `2` as the last character of their horizontal component channel code, and sensors that are intended to be north-aligned will be `N` and `E` as the last character of their horizontal component channel code - as per `docs/MINISEED_METADATA_CONVENTIONS.md`. It is important to mark this distinction: that all sensors have an orientation, but not all sensors are intended to be north-oriented.

### Unknown

When an orientation is unknown, the `Azimuth` value will be *empty* and the `Method` will be *empty*. This ensures that azimuth information is not provided to data users who may otherwise take values indicative of an unknown orientation as the orientation itself.

### Onsite

When an orientation has been measured at the station, the `Azimuth` value will be this measurement and the `Method` will be `onsite`.

The most common onsite orientation methods are using a compass and an azimuith pointing system. A compass measurement involves taking a bearing of the sensor's north component using a compass. An azimuth pointing system (APS) projects a laser onto the sensor's north component marking and uses differential GPS positioning of two antenna at either end of its laser to measure the sensor orientation. Both of these approaches are considered reliable and which method has been used to measure sensor orientation is not noted.

### Offsite

When an orientation is calculated using the data from a station or - as in the case of strong motion sensors in building arrays - using geographic references, the `Azimuth` value is the calculated value and the `Method` is `offsite`.

There are many offsite orientation methods. Besides geographic referencing, the most common method is to use the polarisation of seismic waves from a known source to calculate the sensor orientation. Here both teleseismic earthquake phases and the Rayleigh waves of the secondary microseism recoverable using ambient noise techniques are used. Historically comparisons between surface and borehole seismic sensors have been used to calculate borehole sensor orientation, as has the approach of minimising moment tensor solution residuals by orienting the data from sensors with unknown orientations. To avoid presenting the full description of orientation techniques here, we simply note that offsite orientation methods are always the best possible and orientations are verified by the agreement of two or more methods where possible.

Often offsite orientation methods are used for borehole sensors and other sensors that cannot be oriented using onsite methods, i.e. those in building arrays.

### Method Accuracy

Uncertainty is not recorded for orientation values. However, orientation uncertainties typically fall within +/- 5 degrees, if not a much smaller range.

## Polarity Conventions and Measurement Methods

Polarity values are contained in `delta.git/install/polarities.csv` under the `Reversed` column.

When polarity values are known, they are presented as either `true` or `false` in this column, where `true` indicates a reversed polarity. The absence of a polarity entry for a particular data acquisition system does not indicate that polarities follow the default left-hand convention - just that polarities are unknown. However, in the absence of polarity information polarities are assumed to follow the left-hand convention.

Polarities describe data aquisition system configurations: the combination of sensor, cable, and datalogger. As such, they are not necessarily representative of a particular sensor's internal component polarities. To accomodate this, polarity values are associated to a time period for a site that may not be coincident with particular equipment at that site.

Polarities may describe the polarity of a single data stream (e.g., a Z component) or many data streams (e.g., all components). What the polarity value applies to is distinguished by the `Station,Location,Sublocation,Subsource` columns in the `polarities.csv`.

To illustrate the differences between conventional and reversed polarities, the effects on digitized seismic waveforms, and how the polarities table links into StationXML metadata, please refer to the schematic below (the SVG can be downloaded and opened in Excalidraw).

![SEISMIC_SENSOR_ORIENTATION_POLARITIES-reversed](/docs/SEISMIC_SENSOR_ORIENTATION_POLARITIES-reversed.excalidraw.svg)

### Measurement

Polarity measurements are all data-based. Often polarities are communicated to GeoNet by data users noticing non-default values in their applications. The `Method,Citation` columns in `polarities.csv` captures the source of polarity information.

## Combining and Communication Azimith and Polarity Information to Data Users

Each seismic sensor component has an orientation and a polarity. GeoNet combines these as component orientations in StationXML as follows:
1. The N, X or 1 component is described by the sensor azimuth, with the E/Y/2 and Z component in a left-hand rule with respect to the N/X/1 component.
1. If a component polarity is reversed, a 180 degree phase shift is applied to its orientation.
1. If a component polarity is unknown, no phase shift is applied to its orientation.
1. If azimuth is unknown, no orientation information is provided, meaning polarity information is not conveyed in StationXML even if it exists in delta.git

## Reporting Sensor Orientation and Polarity Information and Inaccuracies

There is always the possibility that sensor orientations and polarities in delta are not accuruate or are absent. Through its operations, GeoNet will ocassionally discover such cases themselves and update its metadata accordingly. If you find an inaccuracy or have orientations or polarities that could help to complete GeoNet's metadata archive, please reach out to info@geonet.org.nz.
