# Stream Naming Conventions

The seismic, and related, data stream naming conventions are based on historical usage together with recommendations from the [SEED manual](https://www.fdsn.org/seed_manual/SEEDManual_V2.4.pdf). Sometimes codes were created and used were no appropriate conventions applied at the time, these were generally left as is even after later conventions were developed.

Each recorded stream is identified by a unique set of codes representing:

__&lt;NETWORK&gt; &lt;STATION&gt; &lt;LOCATION&gt; &lt;CHANNEL&gt;__

## Network Code

The original New Zealand digital seismic networks had three letter network codes related to where the data was being collected, either in the office or where the analogue telemetry was actually digitised. In some cases, mostly related to volcano monitoring, the same data feeds could be digitised and recorded at multiple places. These original codes can be found in the `networks/network.csv` file with the trailing `N` character removed, together with those added subsequently. These network codes are generally used only for display and organisational purposes.

The seismic data is expected to be globally distributed, which means the codes need to be internationally registered.  The current approach is to use the already registered `NZ` network code for all recorded public data. Otherwise the internal `XX` code is used for temporary or private data that will not, or cannot, be exported.

## Station Code

The _National Seismograph Network_ station codes have, over time, developed into three letter codes with generally the last letter being a `Z`. The exceptions being the very oldest installations such as `WEL`. The trailing `Z` has generally made registering station codes easier.

The _National Strong Motion Network_ recording sites tended to have a numbering system with a three digit prefix and a trailing letter.  Equipment at these sites are mostly analogue, or manually digitised data.

Station codes are found in the `networks/stations.csv` file.  For _National Seismograph Network_ sites they generally follow the convention of having three letters, with a last being a `Z`.  Other _Regional Seismograph Networks_ use a four letter code, again with the last letter being a `Z`.  The letter prior to this can also represent the regional or geographic network the station is in.  Exceptions to this are generally the very oldest sites which pre-date the addition of the trailing `Z`.  All other sites use four letter codes with the trailing letter giving some indication of the network or geographical location.  For older sites the suffixes tend to be the initial letter of the original network code as outlined above.

For any future major instrument expansions the use of five letter station codes may be needed.

## Location Code

There are two types of location codes used and are related to either sensor placement or the recording datalogger.  The location code is primarily used to distinguish between multiple sensors installed at a single recording site where the same station code is used. The location code is also used to distinguish between the SOH records taken from any dataloggers which may be installed at the same site.

The SEED format requires that the location code, if used, be case independent (upper-case), and made up only of the standard letters A-Z, numbers 0-9, or the space character.  There are no related SEED usage conventions although originally the global networks tended to use the location code to distinguish between sensors that were installed at various depths down boreholes from those installed on the surface.

The current convention is used to both indicate sensor type and location.

- 0? - Reserved for datalogger SOH channels
- 1? - Reserved for weak motion sensors
- 2? - Reserved for strong motion sensors
- 3? - Reserved for acoustic or pressure sensors
- 4? - Reserved for water level pressure sensors
- 5? - Reserved for geomagnetic sensors
- 6? - Reserved for strain or displacement sensors
- 7? - Reserved for wind measuring sensors
- 8? - Reserved for temperature sensors

There is an informal convention of using `01` for the primary datalogger (generally weak-motion) and `02` for the secondary datalogger (generally strong motion).  This setup has been maintained for sites with only strong-motion recorders as it makes maintaining instrument configurations easier.

Testing, or non-production, dataloggers will have codes using the sequence: 0Z, 0Y, 0X, ... etc.  They should also use a similar sensor location sequence depending on sensor type, e.g. 1Z, 1Y, 1X ... etc.

## Channel Code

Channel codes generally follow the SEED conventions, although some channel codes preceded the conventions and have not been updated.

Apart from a small number of state of health (SOH) channels, the first letter of the code represents a combination of sampling rate and sensor bandwidth, e.g.

- `U` (Ultra broadband sampled every 100s, or SOH sampled every 100s)
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
