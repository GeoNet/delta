## FLUXGATES CALIBRATIONS ##

The `calibrations/fluxgates` contains individual calibration certificates of fluxgates used on Geomagnetic observatories.

A fluxgate sensor calibration certificate contains calibration constants as measured by the Danish Metereological Institute for individual fluxgate magnetometer sensors used across the GeoNet network.

Files in this directory are named based on sensor serial and serial number code.

Further calibrations are applied and are site specific, and those are available in the `install/sensitivity.csv` file.

To obtain site specific fluxgate sensitivity (based on sensor calibration and sensor alignment to magnetic north on site), the following formulas are used:

For the fluxgate:
```
obs[i] = p->coil * p->polarity * (f->volts[i] / p->resolution + s->step * (double) p->bias) + p->polarity * p->offset;
```

For the temperatures:
```
sensor_temp= f->volts[i] * p->gain - 273
```

Those are used to obtain "scale" and "bias" in the `install/sensitivity.csv` file.


nomenclature used in sensor calibration files | w.r.t to X|Y|Z and N|E|vert
--|--
X-coil (nT/mA) | X constant
Y-coil (nT/mA) | Y constant
Z-coil (nT/mA) | Z constant
ε0 (mrad) | X-Y orthogonality
ε1 (mrad) | X misalignment
ε2 (mrad) | Y misalignment
ε3 (mrad | Z-N misalignment
ε4 (mrad) | Z-E misalignment

