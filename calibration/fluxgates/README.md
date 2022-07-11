## FLUXGATES CALIBRATIONS ##

The `calibrations/fluxgates` contains individual calibration certificates of fluxgates used on Geomagnetic observatories.

A fluxgate sensor calibration certificate contains calibration constants as measured by the Danish Metereological Institute for individual fluxgate magnetometer sensors used across the GeoNet network.

ε0, ε1, ε2, ε3, ε4 misalignments are not used in the geomagnetic data collection and translation mechanisms, but are provided here to reflect the intrinsic error associated with each fluxgate sensor measurement.

The coil constants and scale resistors are used to convert measurements of the geomagnetic field (in nanotesla) to volts and counts (recorded at the datalogger).

Files in this directory are named based on sensor serial and serial number code.

Sensor specific calibrations are used in the `install/calibrations.csv`file.

In the `install/calibrations.csv` the scale factor is obtained by dividing the coil constant by the instrument scale resistor. 

Further calibrations are applied and are site specific, and those are available in the `install/gains.csv` file.

To obtain site specific fluxgate gains in each direction (based on sensor calibration and sensor alignment to magnetic north on site), the following formulas are used:

For the fluxgate:

__obs__<sub>i</sub> = __coil__<sub>i</sub> * __polarity__<sub>i</sub> * (__volts__<sub>i</sub> / __resolution__ + __step__ * __bias__<sub>i</sub>) + __polarity__<sub>i</sub> * __offset__<sub>i</sub>

For the temperatures:

__temp__<sub>i</sub> = __volts__<sub>i</sub> * __gain__ - 273.0

Those are used to obtain "scale factor" and "scale bias" in the `install/gains.csv` file.


nomenclature used in sensor calibration files with respect to `X|Y|Z` and `N|E|vert`

| Field | Units   | Description       |
| ----- | ------- | ----------------- |
|X-coil | (nT/mA) | X constant        |
|Y-coil | (nT/mA) | Y constant        |
|Z-coil | (nT/mA) | Z constant        |
|ε0     | (mrad)  | X-Y orthogonality |
|ε1     | (mrad)  | X misalignment    |
|ε2     | (mrad)  | Y misalignment    |
|ε3     | (mrad   | Z-N misalignment  |
|ε4     | (mrad)  | Z-E misalignment  |
|res    | (kohm)  | scale resistor    |



