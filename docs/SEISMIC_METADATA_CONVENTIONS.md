# Seismic Metadata Conventions

In this document we describe the metadata conventions present in the GeoNet seismic data and metadata.

## Overview

GeoNet has two types of seismic stations: weak motion and broadband stations, and strong motion stations. The distinction is made by the primary purpose of the station; whether it exists principally to record weak ground motion (via a seismometer) or strong ground motion (via an accelerometer).

Within each type, stations are grouped by network codes.

At each station there can be many sites, and sites are often referred to by their station code.

Data are sampled at sites, which are located at stations. 

## Stream Naming Convention

A station can host many data streams, and each data stream contains a unique set of metadata describing itself and where it was collected as the combination of:

__&lt;NETWORK&gt; &lt;STATION&gt; &lt;LOCATION&gt; &lt;CHANNEL&gt;__  

The seismic, and related, data stream naming conventions are based on historical usage together with recommendations from the [SEED manual](https://www.fdsn.org/seed_manual/SEEDManual_V2.4.pdf). Sometimes codes were created and used where no appropriate conventions applied at the time, and these have generally been left as is even after later conventions were developed.

### Network Code

As weak motion seismic data is expected to be globally distributed, its network and station codes need to be internationally registered.  The current approach is to use the registered `NZ` network code for all recorded public weak motion data. Otherwise the internal `XX` code is used for temporary or private data that will not, or cannot, be exported internationally. Data of both network codes are available via GeoNet data services.

The full set of network codes can be found in the `networks/network.csv` file. While these are not relevant for modern data access, much of GeoNet's history is contained in these codes and they are still in use for metadata management in _delta_ and for data operations within GeoNet. For example, a station can belong to a regional seismic network like the Taranaki volcano seismic network `TR` and be distributed with the network code `NZ`. 

### Station Codes

Station codes are assigned at the first installation of a seismic sensor at a location and do not change.

Station codes are unique within a given reference set. Weak motion station codes are unique for the global set of station codes managed by the ISC and are registered with the ISC to safeguard this uniqueness. Strong motion station codes are unique only within the set of such codes used by GeoNet.     

#### Weak Motion Station Code Conventions

Station codes are found in the `networks/stations.csv` file.

_National Seismograph Network_ station codes use a three letter code with the last letter being a `Z`. The exceptions being the very oldest installations such as `WEL`.

_Regional Seismograph Networks_ use a four letter code, again with the last letter being a `Z`.  Exceptions to this are generally the very oldest stations which pre-date the addition of the trailing `Z`. In these cases, the trailing letters can represent the regional network the station is in, such as in the Auckland Volcano Seismic Network where many stations have `AK` suffix, or otherwise indicate the network or geographical location of the station.  For many older stations the suffixes tend to be the initial letter of the original network code.

For both national and regional network station codes, the first two letters try to give an indication of where the station is (i.e. they will be an abbreviation of a close town or farm station name).

While the legacy of network and station codes introduces complexity to the seismic metadata, changes to either code for a site are not anticipated unless physical changes require it.

#### Strong Motion Station Code Conventions

In the past, the National Strong Motion Network recording sites tended to have a numbering system with a three digit prefix and a trailing letter. The current National Strong Motion Network site code naming convention is to use four letter codes describing where the station is, e.g. LPLS is near Lake Paringa, PRNS is near Paringa.

As with weak motion, strong motion stations end in a particular character: 'S'.

Due to the utility of co-locating strong motion sensors with sensors of other types, many strong motion sites exist at weak motion stations. Such sites will have the weak motion station code but a different location code.

### Location Codes

The location code is primarily used to distinguish between multiple sensors installed at a single recording station where the same station code is used.

Location codes are associated to station codes and we refer to location codes when we make reference to a "site", even though we often use the host station code to make this reference.

There are two types of location codes used and they are related to either sensor placement or the recording datalogger. The location code can also be used to distinguish between the state of health (SOH) records taken from any dataloggers which may be installed at the same site.

In their role as metadata detailing sensor placement, location codes convey:
 1. Sensor positions for one or more sensors installed coincidentally or in sequence.
 2. Sensor types for one or more sensors installed coincidentally or in sequence.
 
Location codes are two characters, with the first character denoting the sensor (or sensor data) type, and the second character denoting the sensor position. The first character follows groupings as:

- 0? - Reserved for datalogger SOH channels
- 1? - Reserved for weak motion sensors
- 2? - Reserved for strong motion sensors
- 3? - Reserved for acoustic or pressure sensors
- 4? - Reserved for water level pressure sensors
- 5? - Reserved for geomagnetic sensors
- 6? - Reserved for strain or displacement sensors
- 7? - Reserved for wind measuring sensors
- 8? - Reserved for temperature sensors

There is an informal convention of using `01` for the primary datalogger (generally weak motion) and `02` for the secondary datalogger (generally strong motion).  This setup has been maintained for sites with only strong-motion recorders as it makes maintaining instrument configurations easier.

Testing, or non-production, dataloggers will have codes using the sequence: 0Z, 0Y, 0X, ... etc.  They should also use a similar sensor location sequence depending on sensor type, e.g. 1Z, 1Y, 1X ... etc.

## Naming Conventions When Moving or Installating Different Sensors

Sites are associated with stations at the start of data collection from the site. Data collection is from a datalogger, which is connected to a sensor. 

When a sensor is moved at a station or a new sensor is installed, the station and location code describing that installation follows these conventions:

1. If the sensor is of the same type and in the same position as the previous sensor, neither station nor location code changes.
1. If the sensor is of a different type but in the same position as the previous sensor, the station code remains the same but the location code changes .
1. If the sensor is more than 1 m from the position of the previous sensor, the station code remains the same but the location code changes. 
1. If the sensor is more than 200 m from the position of the previous sensor, the station code changes. Here a new station may need to be made with a new set of location codes describing data collection at the station. 

These conventions reflect GeoNet's understanding of the purpose of its seismic station and location codes. Where location codes are used to distinguish between different sensor types or positions at a station, station codes are used to distinguish between the different nodes of a sensor network. What makes a station distinct is the relationship of its data to that of other stations in the same network. What makes a station continuous is the relationship of its data to that recorded previously at the same station. 
 When a sensor position or type changes, we cannot assume its data is comparable to what was produced previously under the same location or station code. Our conventions try to capture those sensor changes that alter a station or site's data beyond the point of comparability or continuity.    

These conventions are ultimately used at the discretion of the person(s) responsible for the metadata describing equipment changes. If, for example, a sensor moved less than 200 m but the geologic or local site conditions changed substantially, a new station may be established to reflect this change.

### Borehole Naming Conventions

Sensors in a borehole invoke a special set of naming conventions.

For borehole sensors, the station code will always be the same regardless of sensor position in the borehole, so long as the sensor is in the borehole. However, the location code changes as normal; for sensor position changes of more than 1 metre, the location code changes; for sensor type changes, the location code changes.

Sensors at the surface adjacent to boreholes are covered by these special conventions; borehole surface sensors are considered borehole sensors regarding naming conventions.

## Channel Codes

Channel codes generally follow the SEED conventions, although some channel codes preceded the conventions and have not been updated.

Apart from a small number of SOH channels, the first letter of the code represents a combination of sampling rate and sensor bandwidth, e.g.

- `W` (Ulta-ultra long period; used for DART data)
- `U` (Ultra-long period broadband sampled every 100s, or SOH sampled every 100s)
- `V` (Broadband sampled every 10s, or SOH sampled every 10s)
- `L` (Broadband sampled at 1Hz, or SOH sampled at 1Hz)
- `B` (Broadband sampled at between 10 and 80 Hz, usually 10 or 50 Hz)
- `S` (Short-period sampled at between 10 and 80 Hz, usually 50 Hz)
- `H` (Broadband sampled at or above 80Hz, generally 100 or 200 Hz)
- `E` (Short-period sampled at or above 80Hz, generally 100 Hz)

The second letter represents the sensor type, e.g..

- `H` (Weak motion sensor, e.g. measuring velocity)
- `N` (Strong motion sensor, e.g. measuring acceleration)
- `L` (Low gain sensor, usually velocity)
- `M` (Mass position, used for monitoring broadband sensors)
- `D` (Barometer or pressure sensor)
- `K` (Temperature sensor)
- `A` (Tiltmeter)
- `F` (Geomagnetic sensor)
- `T` (Water level sensor)

The third letter either represents the sensor orientation or a processing stage.

- `Z,N,E` (Three component sensor aligned to North)
- `Z,1,2` (Three component sensor with non-aligned orientation, generally used for borehole sensors or strong motion recorders)
- `X,Y,Z` (Three component sensors with site specific orientations, generally used for building arrays)
- `U,V,W` (Three component sensors with non-standard orientations, generally used for mass positions of broadband sensors)
- `Z` (Single component vertical sensor, or pressure sensor used to measure height).
- `F` (Sensors that have no orientation, such as pressure sensors or geomagnetic fields)
- `H` (Sea level streams which have been corrected to height, generally used for processing)
- `T` (Sea level streams which have been de-tided, generally used for processing)

## Example Channels

### Data Channels

| Channel               | Measurement
|-----------------------|-------------
| `VHZ VHN VHE VH1 VH2` | velocity -- broadband
| `LHZ LHN LHE LH1 LH2` | velocity -- broadband
| `BHZ BHN BHE BH1 BH2` | velocity -- broadband
| `HHZ HHN HHE HH1 HH2` | velocity -- broadband
| `SHZ SHN SHE SH1 SH2` | velocity -- short period
| `EHZ EHN EHE EH1 EH2` | velocity -- short period
| `SLE SLN SLZ` | velocity -- low gain sensor
| `HNZ HNN HNE HN1 HN2` | acceleration -- strong motion
| `BNZ BNN BNE BN1 BN2` | acceleration -- strong motion
| `HDF BDF LDF` | pressure -- barometer)
| `HDH BDH LDH` | pressure -- hydrophone)
| `LDA HDA` | pressure -- microphone
| `VTZ LTZ BTZ` | pressure - water
| `VTH LTH BTH` | pressure - corrected water depth
| `VTT LTT BTT` | pressure - water depth with tide removed
| `LAX LAY` | tilt
| `BKO LKO VKO` | temperature (e.g. lake)
| `LKD LKS` | temperature (e.g. geomag sensor)
| `LFX LFY LFZ` | geomagnetic -- field values
| `LFF LFD` | geomagnetic -- full field values
| `CRX` | tidal height CREX encoded messages

### Data Quality Channels

| Channel               | Measurement
|-----------------------|-------------
| `VMZ VM1 VM2 VMN VME` | mass position - broadband sensor offset
| `LMZ LM1 LM2 LMN LME` | mass position - broadband sensor offset
| `VMU VMV VMW`         | mass position - broadband sensor offset
| `LMU LMV LMW`         | mass position - broadband sensor offset
| `LEQ VEQ`             | geomagnetic - absolute field observation quality
| `CAL`                 | sensor calibration details
| `BTL`                 | packet latency times

### State of Health Channels

| Channel       | Measurement
|---------------|-------------
| `LOG`         | General log messages in encoded text format
| `ACE`         | Clock timing messages in encoded text format
| `UEP VEP LEP` | Instrument voltage
| `LEB`         | Instrument internal battery voltage
| `VEC LEC`     | Instrument current
| `UK1 UK2`     | Internal instrument temperatures
| `VKI LKI`     | Internal instrument temperature
| `LII`         | Instrument humidity
| `LEU`         | Instrument CPU load
| `VEM LEM`     | Instrument percent disk free or buffer full
| `UCQ VCQ LCQ` | Clock quality
| `LCE`         | Clock phase error
| `VEA`         | Clock antenna current
| `VCO`         | Clock VCO frequency control
| `UCD`         | Clock drift
 
