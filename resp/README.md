## RESPONSE FILES

The goal of the `resp` directory is to build a `golang` source module that allows easy instrument
response information discovery for downstream programs but primarily to build `stationxml` representations
of the meta information.

### helper files

The main `golang` output is the automatically built `auto.go` file. There are some helper functions and
structure definitions in the `response.go` and `streams.go` files. The later provides
a mechanism to discover all possibly configured streams for a datalogger and sensor pair.

### generate

The code needed to build the `auto.go` file using the raw configuration files can be found in the
`generate` sub-directory. This code is usually run in the main `resp` directory via a call to
`go generate`.

This is managed via the header line in the `response.go` file, i.e.

```
//go:generate bash -c "go run generate/*.go | gofmt > auto.go"
```

If for some reason this command fails, there is a likelihood that the next run will also fail.
It will try to run `go generate` and when it finds an empty `auto.go` file it will stop and
complain. The solution is simply to remove the empty `auto.go` file.

The generated `auto.go` file should be committed into the repo as per the configuration files
or other source code.

### configuration files

The configuration files are stored in `YAML` format under the `responses` directory. There is
no requirement on the directory or file layout in this sub-directory, other than the files
need to have a suffix of `.yaml`. All files have the same overall format
although not all configuration sections will be present.
The contents of each file are merged together prior to processing and building the `auto.go` file.

For maintainability the files have been split into areas related to sensors, dataloggers, configurations
related to sensor and datalogger pairs, and filters which presently tend to be FIR filter definitions
as used in the dataloggers.

The various sections are as follows, the actual definition of the `paz`, `polynomial`, and `fir` sections
map fairly well to the `SEED` definitions for the relative response details.

#### paz

Used to describe a Poles & Zeros filter.

``` yaml
paz:
  1500Hz Bessel 3P-LP:
    code: A
    type: Laplace transform analog stage response, in rad/sec.
    notes: ""
    poles:
    - (-9904.799805+3786i)
    - (-9904.799805-3786i)
    - (-12507+0i)
    zeros: []
```

Note the mechanism to insert a complex number, this is needed to overcome the lack of complex numbers in the `YAML` importing mechanism.

#### polynomial

Used to describe a polynomial response. 

``` yaml
polynomial:
  Druck PTX-1830-A:
    gain: 0.0008
    approximationtype: MACLAURIN
    frequencylowerbound: 0
    frequencyupperbound: 0
    approximationlowerbound: 0
    approximationupperbound: 20
    maximumerror: 0
    coefficients:
    - 0.004
    - 0.0008
```

#### fir

Used to describe a FIR filter response.

``` yaml
---
fir:
  Q330_FLbelow100-1:
    causal: true
    symmetry: none
    decimation: 1
    gain: 1
    factors:
    - 1.2199295e-16
    - 3.1619205e-10
    - -4.3146524e-08
    - -5.6355576e-07
    - ...
``` 

#### sensor-model

Provides a basic description of a sensor model, which also includes
the internal components and their relative orientation.

``` yaml
sensor-model:
  LE-3Dlite:
    type: Short Period Seismometer
    description: LE-3Dlite
    manufacturer: "Lennartz"
    vendor: ""
    components:
    - azimuth: 0
      dip: -90
    - azimuth: 0
      dip: 0
    - azimuth: 90
      dip: 0
``` 

#### datalogger-model

Provides a basic description of a datalogger model.

``` yaml
datalogger-model:
  Q330/3:
    type: Datalogger
    description: Q330
    manufacturer: Quanterra
    vendor: ""
```

#### filter

A filter is used to join generic filter responses (e.g. `paz`, `polynomial`, or `fir`) into
a set of response stages combining the responses with gain, sampling rate information, and units
as required.

``` yaml
filter:
  LE-3Dlite:
  - type: paz
    lookup: LE-3Dlite
    frequency: 15
    samplerate: 0
    decimate: 0
    gain: 400
    scale: 1
    correction: 0
    delay: 0
    inputunits: m/s
    outputunits: V
```

#### response

The `response` section joins everything together. Given a list of sensors, and their associated filters, together with
a list of datalogger sampling rates, and filters, a full set of possible combinations of input streams can be built.

This is the primary information that goes into the `auto.go` file, together with the simpler sensor and datalogger details.

There is a built in mechanism to handle small variations in instrument capabilities, using a list allows similar dataloggers or sensors
models to be configured together rather than each needing a separate entry.
This also constrains which sensors can be attached to which dataloggers and will be needed to be updated for new equipment
as well as the individual instrument details.

``` yaml
response:
  Quanterra Dataloggers Connected to STS-2 Sensors:
    sensors:
    - sensors:
      - STS-2
      filters:
      - STS-2
      channels: ZNE
      reversed: false
    dataloggers:
    - dataloggers:
      - Q330HR/6
      type: CG
      label: HH
      samplerate: 100
      frequency: 1
      storageformat: Steim2
      clockdrift: 0.0001
      filters:
      - Q330HR_FLbelow100-100
      reversed: false
    - dataloggers:
      - Q330HR/6
      type: CG
      label: LH
      samplerate: 1
      frequency: 0.1
      storageformat: Steim2
      clockdrift: 0.0001
      filters:
      - Q330HR_FLbelow100-1
      reversed: false
```
