## RESPONSE FILES

The goal of the `resp` directory is to build a `golang` source module that allows easy instrument
response information discovery for downstream programs but primarily to build `stationxml` representations
of the meta information.

### files

The `files` directory holds XML response snippets. These are short representations of a StationXML style set of response stages for
a given sensor or datalogger. The full StationXML response is found by joining these stages together as required.
It is heavily influenced by the NRL (Nominal Response Library) currently located at https://ds.iris.edu/ds/nrl/.

The naming convention is not as strict, but it generally follows the form of:

#### sensors

The sensor response is used for standalone sensors which are independent of sampling mechanism.

The files follow the format of:

```
sensor_<MAKE>_<MODEL>.xml
```

Where `MAKE` and `MODEL` can be upper and lower case letters, numbers, and the dash symbol to represent spaces.

#### dataloggers

The datalogger response is used for standalone dataloggers which are independent of attached sensors.

The files follow the format of:

```
datalogger_<MAKE>_<MODEL>_<SCALE>_<RATE>.xml
```

Where `MAKE` and `MODEL` can be upper and lower case letters, numbers, and the dash symbol to represent spaces and
`RATE` represents the expected input sampling rate or period, either as `sps` (samples per second), or `s` (seconds per sample).
The `<SCALE>` can be used to represent different gain settings on the instrument, and is usually of the form `24bits`

#### combined

The combined response stages are for use when the sensor and datalogger are combined into a single unit (such as for a digital sensor).

The files follow the format of:

```
combined_<MAKE>_<MODEL>_<RATE>.xml
```

Where `MAKE` and `MODEL` can be upper and lower case letters, numbers, and the dash symbol to represent spaces and
`RATE` represents the expected input sampling rate or period, either as `sps` (samples per second), or `s`
(seconds per sample).

#### derived

For response stages that are derived from a previous stage and are not directly related to an instrument or a sensor.
Examples include simple unit conversions.

The files follow the format of:

```
derived_<LABEL>_<RATE>.xml
```

Where `LABEL` can be upper and lower case letters, numbers, and the dash symbol to represent spaces and
`RATE` represents the expected input sampling rate or period, either as `sps` (samples per second), or `s`
(seconds per sample).
