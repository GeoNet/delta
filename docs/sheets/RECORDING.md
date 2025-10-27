# COMMON SIGNAL RECORDING TASKS

* :hammer: [Combined editing steps](#steps)
* :scroll: :one: [Creating or updating a recording  __station__](#station)
* :scroll: :two: [Creating or updating a recording __site__](#site)
* :scroll: :three: [Installing or updating a __sensor__](#sensor)
* :scroll: :four: [Installing or updating a __recorder__](#recorder)
* :scroll: :five: [Installing or updating a __datalogger__](#datalogger)
* :scroll: :six: [Installing or updating a __datalogger__ to __sensor__ connection](#connection) :dragon:
* :scroll: :seven: [Installing or updating a __datalogger__ to __recorder__ stream](#stream) :dragon:

The dataloggers and recorders produce waveform data that is tagged by recording __station__ code,
recording __site__ code, a "stream" code indicating orientation, and a sampling rate.
The __install__ files are used to map this combination back to actual sensors and dataloggers and their
configurations.

Installed __sensors__ and __recorders__ have their installation times and orientations described in the
__install/sensors.csv__ and __install/recorders.csv__ files respectively.
Deployed __dataloggers__ are described in the __install/dataloggers.csv__ file.

Attaching sensors to dataloggers is managed through linking the __install/connections.csv__ file
to the __install/dataloggers.csv__ file. This is not required for __recorders__ as these are assumed
to be part of the sensor being installed at a recording __site__.
Linking is done via matching a datalogger to a notional "place", and optional operational "role". These
are presently not constrained, other than that they should match when a connection is needed.

Datalogger configuration, through the __install/streams.csv__ file, is meant to provide a mechanism
for mapping recorded data to a datalogger or recorder configuration. This is an overlay
template which isn't dependent on the actual equipment being installed, only the broad settings.

## <a name="steps">_Overall steps_

> * :file_folder: Using a suitable mechanism create a new _git_ branch for the changes.
> * :one: :pencil2: Update the __network/stations.csv__ file to add or update any recording __station__ information, the entries will need to be in order of __station__ code.
> * :two: :pencil2: Update the __network/sites.csv__ file to add or update any recording __site__ information, the entries will need to be in order of __station__ code then __site__ code.
> * :three: :pencil2: Update the __install/sensors.csv__ file to add any new __sensor__ installations, the entries will need to be in order of make, model, serial number, and installation time.
> * :four: :pencil2: Update the __install/recorders.csv__ file to add any new __recorder__ installations, the entries will need to be in order of make, model, serial number, and installation time.
> * :five: :pencil2: Update the __install/dataloggers.csv__ file to add any new __datalogger__ deployments, the entries will need to be in order of make, model, serial number, and deployment time.
> * :six: :pencil2: Update the __install/connections.csv__ file to add any new __datalogger__ to __sensor__ connection details, the entries will need to be in order of __station__ code, __site__ code, and connection time.
> * :seven: :pencil2: Update the __install/streams.csv__ file to add any new __datalogger__ or __recorder__ stream configuration details, the entries will need to be in order of __station__ code, __site__ code, sampling rate, and stream operational times.
> * :open_file_folder: Build the pull request with a meaningful title.
> * :link: Assign suitable tags and set reviewers.
> * :repeat: If the tests fail, the above changes may need some iteration until they pass.
> * :sos: If the tests are still failing, escalate as this may indicate some inconsistency within the network configuration.
> * :ok: Once the tests have passed and the pull request reviewed, depending on policy, the pull request can be merged

## :one: <a name="station">_Creating or updating a recording station_

> ### :page_with_curl: Files to update
>
> * __network/stations.csv__
>
> ### :page_with_curl: Reference files
>
> * __network/networks.csv__
>

New stations will require a station "code" to be assigned using the appropriate mechanism.

It is intended that recording __stations__ be fixed, any adjustments will be applied retrospectively.
The operation time should span all installations at the __station__.

> ### :information_source: General requirements
> 
> * The recording __station__ "code".
> * The __network__ "code" attached to the recording __station__.
> * The "name" of the recording __station__.
> 
> ### :information_source: Field requirements
> 
> * The location of the recording __station__.
> * The operational time of the recording __station__.
>
> ### :heavy_check_mark: Delta prerequisites
>
> * There can only be one __station__ for any given code.
> * There must be an entry in the __networks.csv__ file for the __station's__ __network__ code.
>
> ### :small_orange_diamond: Delta checks
>
> * The operational start time must be before the end time.

## :two: <a name="site">_Creating or updating a recording site_

> ### :page_with_curl: Files to update
>
> * __network/sites.csv__
>
> ### :page_with_curl: Reference files
>
> * __network/stations.csv__
>

New recording __sites__ will require a location "code" to be assigned using the appropriate mechanism.
There is a convention used for __location__ codes related to the type of signal being recorded:

* 1X -- weak motion sensors
* 2X -- strong motion sensors
* 3X -- acoustic sensors
* 4X -- pressure sensors
* 5X -- geomagnetic sensors
* 6X -- displacement sensors
* 7X -- weather style measurements
* 8X -- temperature sensors
* 9X -- tilt sensors

Each __site__ location must be different at a given __station__, although there is some leeway
when installing the actual sensor, such as depth down a borehole.

It is intended that recording __sites__ be fixed, any adjustments will be applied retrospectively.

> ### :information_source: General requirements
> 
> * The recording __station__ "code".
> * The recording __site__ location "code".
> 
> ### :information_source: Field requirements
> 
> * The location of the recording __site__.
> * The operational time of the recording __site__.
> 
> ### :heavy_check_mark: Delta prerequisites
>
> * There must be an entry in the __stations.csv__ file for the __site's__ __station__ code.
>
> ### :small_orange_diamond: Delta checks
> 
> * Location "codes" must be unique at any given recording __site__.
> * The operational start time must be before the end time.
> * Recording __site__ operational times must be within the associated __station's__ operational times.

## :three: <a name="sensor"></a>_Installing or replacing a sensor_

> ### :page_with_curl: Files to update
>
> * __install/sensors.csv__
>
> ### :page_with_curl: Reference files
>
> * __assets/sensors.csv__
> * __network/stations.csv__
> * __network/sites.csv__
>
> ### :information_source: General requirements
>
> * The recording __station__ "code".
> * The recording __site__ location "code".
>
> ### :information_source: Field requirements
>
> * The make, model, and serial number of the sensor being installed or replaced.
> * The azimuth and dip of the installed sensor.
> * Any offsets from the given recording __site__ location.
> * Any __sensor__ biases or factors related to the __sensor__ installation.
> * The times when the sensor is installed or removed.
>
> ### :heavy_check_mark: Delta prerequisites
> 
> * The sensor needs to be listed in the __assets/sensors.csv__ file.
> * The recording __station__ need to be present in the __network/stations.csv__ file.
> * The recording __site__ need to be present in the __network/sites.csv__ file.
> 
> ### :small_orange_diamond: Delta checks
> 
> * A given __sensor__ can only be installed at a single __site__ at any given time.
> * An individual __site__ can only have one __sensor__ or __recorder__ installed at any given time.
> * Installed __sensor__ times must be within the recording __site's__ operational time window.

## :four: <a name="recorder"></a>_Installing or replacing a recorder_

> ### :page_with_curl: Files to update
>
> * __install/recorders.csv__
>
> ### :page_with_curl: Reference files
>
> * __assets/recorders.csv__
> * __network/stations.csv__
> * __network/sites.csv__
>

Recorders are managed as a __sensor__ and __datalogger__ pair.
They are installed and removed at the same times and can only be
installed at a single recording __site__. This places a constraint
that there can only be one __sensor__ attached to the __recorder__.

> ### :information_source: General requirements
>
> * The recording __station__ "code".
> * The recording __site__ location "code".
>
> ### :information_source: Field requirements
>
> * The azimuth and dip of the installed __sensor__.
> * The make, model, sensor model, and serial number of the __recorder__ being installed or replaced.
> * The times when the recorder is installed or removed.
>
> ### :heavy_check_mark: Delta prerequisites
> 
> * The __recorder__ needs to be listed in the __assets/recorders.csv__ file.
> * The recording __station__ need to be present in the __network/stations.csv__ file.
> * The recording __site__ need to be present in the __network/sites.csv__ file.
> 
> ### :small_orange_diamond: Delta checks
> 
> * A given __recorder__ can only be installed at a single __site__ at any given time.
> * An individual __site__ can only have one __sensor__ or __recorder__ installed at any given time.
> * Installed __recorder__ times must be within the recording __site's__ operational time window.

## :five: <a name="datalogger"></a>_Deploying or replacing a datalogger_

> ### :page_with_curl: Files to update
>
> * __install/dataloggers.csv__
>
> ### :page_with_curl: Reference files
>
> * __assets/dataloggers.csv__
> * __network/stations.csv__
> * __network/sites.csv__
>
> ### :information_source: General requirements
>
> * The notional "place", and optional operational "role", of the datalogger being installed or replaced.
>
> ### :information_source: Field requirements
>
> * The model, make, and serial number of the datalogger being installed or replaced.
> * The times of the datalogger installation or removal.
>
> ### :heavy_check_mark: Delta prerequisites
> 
> * The __datalogger__ needs to be listed in the __assets/dataloggers.csv__ file.
> 
> ### :small_orange_diamond: Delta checks
> 
> * A given __datalogger__ can only be installed at a single "place" with a given "role" at any given time.
> * An individual datalogger "place" and "role" can only have one __datalogger__ deployed at any given time.

## :six: <a name="connection"></a>_Connecting or disconnecting a sensor to a datalogger_

> ### :page_with_curl: Files to update
>
> * __install/connections.csv__
>
> ### :page_with_curl: Reference files
>
> * __network/stations.csv__
> * __network/sites.csv__
> * __install/sensors.csv__
> * __install/dataloggers.csv__
>

Attaching sensors to dataloggers is managed through linking the __install/connections.csv__ file
to the __install/dataloggers.csv__ file. These are linked by using a notional datalogger "place" and
an optional "role". The "role" can be used to distinguish the datalogger functions when there are more
than one at a given "place". This linkage can be thought of as representing the "sensor cable", ignoring
physical constraints, such as plugs, this link will be independent of that actual datalogger or sensor.

> ### :information_source: General requirements
>
> * The recording __station__ "code".
> * The recording __site__ location "code".
> * The notional "place", and optional operational "role", of the datalogger being installed or replaced.
>
> ### :information_source: Field requirements
>
> * The times when the given datalogger "place" and "role" is connected to the recording site.
>
> ### :heavy_check_mark: Delta prerequisites
> 
> * The recording __station__ need to be present in the __network/stations.csv__ file.
> * The recording __site__ need to be present in the __network/sites.csv__ file.
> 
> ### :small_orange_diamond: Delta checks
> 
> * An individual connection "place" and "role" can only be assigned to one recording __site__ at any given time.

## :seven: <a name="stream"></a>_Creating recorder or datalogger configuration stream_

> ### :page_with_curl: Files to update
>
> * __install/streams.csv__
>
> ### :page_with_curl: Reference files
>
> * __network/stations.csv__
> * __network/sites.csv__
>

Building the meta-data associated with recorded data requires some general datalogger or recorder configuration information.
This is generally independent of the actual datalogger or sensor used and can span multiple installations.

This information is meant to be used to build a response for a given set of recordings, rather than to define what
data is expected or available.

Use an end date of `9999-01-01T00:00:00Z` to indicate that a streams is currently operational.

> ### :information_source: General requirements
> 
> * The recording station code.
> * The recording site location code.
> * The expected recording sample rates.
>
> ### :information_source: Field requirements
>
> * Whether the sensor is installed axially (_not aligned north south_)?
> * Whether, for some reason or another, the sensor signal has been identified as being reversed over the given time window?
> * The time window that this information is valid for.

> ### :heavy_check_mark: Delta prerequisites
> 
> * The recording __station__ need to be present in the __network/stations.csv__ file.
> * The recording __site__ need to be present in the __network/sites.csv__ file.
> 
> ### :small_orange_diamond: Delta checks
> 
> * The sampling rate needs to be within the expected list: 0.1, 1, 10, 50, 100, or 200 Hz.
> * Only one stream with a given sampling rate can be defined for a recording __site__ at any given time.
>
