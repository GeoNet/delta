# Strong Motion Deployment Best Practices (azimuth/orientation and miniSEED channel names)

This document acts to serve as a guide for field best practices when deploying and orienting strong motion sensors and recorders in the GeoNet sensor network. For the purpose of this document both analogue strong motion FBA sensors and FBA/MEMS recorders will be referred to as sensors. All GeoNet strong motion sites currently operate with a 200Hz sampling rate therefore their miniSEED channel naming convention will always have a prefix of HN*.  

## Strong motion site types 

Regarding instrument installation, there are three types of strong motion site (excluding structural arrays):

1. Sites in buildings
2. Sites in GeoNet field infrastructure (cabinet, vault, VSAT hut)
3. Sites in boreholes

### 1) Sites in buildings

Sensors at sites in buildings are installed inside cases which are typically aligned with the walls of the building. The orientation of these sensors should be in line with the edges of their host cases, which are aligned with the walls of the host building. In short, these sensors should be aligned with the walls of the host building. Their orientation can then be taken from the orientation of these walls. The miniSEED channel naming for these sites will be HNZ|1|2 (axial). If the building structure happens to be aligned to true north, there will be no change to this naming schema and the azimuth will be recorded as 0 (to maintain metadata consistency/continuity).

### 2) Sites in GeoNet field infrastructure (cabinet, vault, VSAT hut)

Sensors at sites in GeoNet infrastructure are installed in light structures without cases. There is varying amounts of room for orienting sensors in these installations. Here it is desirable that sensors are true north-aligned, but if this is not possible or cannot be done accurately due to space constraints, an alignment with the host structure’s walls is acceptable. The orientation can then be taken from the orientation of these walls. For the majority that are true north-aligned the miniSEED channel naming will be HNZ|N|E (non-axial) and those in the rare instances that are aligned to the structure will be named HNZ|1|2 (axial).

### 3) Sites in boreholes

Sensors at sites in boreholes generally cannot be aligned to true north if the depth of the borehole is >10m as rotating the sensor is impratical at these depths. The orientation can then be taken from either an onsite method (surface sensor comparison) or offsite method (polarisation analysis), please see SEISMIC_SENSOR_ORIENTATION_POLARITIES.md for further information regarding these methods. The miniSEED channel naming for these sites will be HNZ|1|2 (axial). For shallow boreholes <10m rotating the sensor may be achievable and in these circumstances (if visbility allows an accurate measurement) the sensor will be oriented true north with miniSEED channel naming HNZ|N|E.     

### Find out more

Further information regarding strong motion metadata an be found via the following links:

Strong motion sensor azimuths and depths - https://github.com/GeoNet/delta/blob/main/install/sensors.csv
Strong motion recorder azimuths - https://github.com/GeoNet/delta/blob/main/install/recorders.csv
Strong motion axial or non-axial streams - https://github.com/GeoNet/delta/blob/main/install/streams.csv