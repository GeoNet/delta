## FLUXGATES CALIBRATIONS ##

The `calibrations/fluxgates` contains individual calibration certificates of fluxgates used on Geomagnetic observatories.

Within delta, a fluxgate sensor can be made by a "geomagnetic sensor" and a "driver", that are treated here as an individual sensor with a composite serial number, where applicable.

A fluxgate sensor calibration certificate contains calibration constants as measured by the Danish Metereological Institute for individual fluxgate magnetometer sensors used across the GeoNet network.

ε0, ε1, ε2, ε3, ε4 misalignments are not used in the geomagnetic data collection and translation mechanisms, but are provided here to reflect the intrinsic error associated with each fluxgate sensor measurement.

The coil constants and scale resistors are used to convert measurements of the geomagnetic field (in nanotesla) to volts and counts (recorded by the datalogger). 

The temperature constants are used to convert measurements of the temperature (in degree kelvin) to volts and counts (recorded by the datalogger)

Coil and resistor constants are different for different sensor serial numbers. Temperature constants are the same for a given manufacturer and model.

Files in this directory are named based on sensor serial and serial number code.

Sensor specific calibrations are used in the `install/calibrations.csv`file. In the `install/calibrations.csv` the scale factor is described via the coil constant divided by the instrument scale resistor. 

Further calibrations are applied and are site specific, and those are available in the `install/gains.csv` file. 

To obtain site specific fluxgate gains in each direction (based on sensor calibration and sensor alignment to magnetic north on site), the following formulas are used:

For the fluxgate:

__obs__<sub>i</sub> = __coil__<sub>i</sub> * __polarity__<sub>i</sub> * (__volts__<sub>i</sub> / __resolution__ + __step__ * __bias__<sub>i</sub>) + __polarity__<sub>i</sub> * __offset__<sub>i</sub>

For the temperatures:

__temp__<sub>i</sub> = __volts__<sub>i</sub> * __gain__ - 273.0


Nomenclature used in sensor calibration files with respect to `X|Y|Z` and `N|E|vert`

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

If calibrations are provided for Sensor and Electronic (driver), those are reported in the calibration file, alongside the component serial number.


Nomenclature used in sensor calibration files for each component ("number" in `install/components.csv`) 

| channel number | Fluxgate Component | Geographical Component |
| -------------- | ------------------ | ---------------------- |
| 0              | X                  | N                      |
| 1              | Y                  | E                      |
| 2              | Z                  | vertical               |
| 3              | temperature sensor |                        |
| 4              | temperature driver |                        |
