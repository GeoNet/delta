# COMMON ENVIROSENSOR TASKS

* :hammer: [Combined editing steps](#steps)
* :scroll: :one: [Creating, updating or closing an __Envirosensor__ site](#openclose) :dragon: 
* :scroll: :two: [Specific change examples](#specifcexamples) :monkey:

Envirosensors are installed at monitoring __stations__ and measure data from one or more physical __features__. A __station__ can 
have multiple __sites__ (in this case, each site equates to a sensor), with location and subsensor codes that are used to distinguish the type of quantity that is measured.

Envirosensors are composed of one or more __sensors__ connected to a __datalogger__. The recorded quantity might need to be 
converted to the International System of Units (SI) via  a __scale factor__ and __scale bias__, if the sensor does not natively output SI units. The __scale factor__ and __scale bias__ can also be used to adjust observed values relative to reference or initial values.


A complete list of files that cover Envirosensor metadata is below.

folders | tables
-------|-------
assets | dataloggers.csv, sensors.csv
environment | features.csv 
install | connections.csv, dataloggers.csv, gains.csv, sensors.csv, streams.csv, firmware.csv
network | sites.csv, stations.csv


------
## <a name="steps"></a>_Overall steps_

> * :file_folder: Using a suitable mechanism create a new _git_ branch for the changes.
> * :pencil2: Update the csv tables in __network__, __install__, __environment__ and __assets__ folders, depending on the equipment change.
> * :computer: Locally run the test to ensure updates are consistent.
> * :open_file_folder: Build the pull request with a meaningful title.
> * :link: Assign suitable tags and set reviewers.
> * :repeat: If the tests fail, the above changes may need some iteration until they pass.
> * :sos: If the tests are still failing, escalate as this may indicate some inconsistency within the network configuration.
> * :ok: Once the tests have passed and the pull request reviewed, depending on policy, the pull request can be merged.

------

## :one: <a name="openclose"></a>_Creating, updating or closing an Envirosensor site_

For new stations a __site name__ and  __site code__ will need to be assigned by the appropriate mechanism.

> ### :page_with_curl: Files to update
> * __network/stations.csv__
> * __network/sites.csv__
> * __install/connections.csv__
> * __install/dataloggers.csv__
> * __install/firmware.csv__
> * __install/gains.csv__
> * __install/sensors.csv__
> * __install/streams.csv__
> * __environment/features.csv__


> ### :page_with_curl: Reference files
> * __network/networks.csv__
> * __assets/dataloggers.csv__
> * __assets/sensors.csv__


> ### :information_source: General requirements
> * The Envirosensor __station__ and __site__ codes.
> * The feature and quantity that is measured. In some some cases, a single datalogger may collect data from two or more closely located, but distinct features.
> * ....
>
> ### :information_source: Field requirements
> * The location code and subsensor codes and their relation with the observed quantity
> * The sensor make, model and serial numbers
> * The association between location and sublocation code and measured quantity
> * Coordinates of the __station__ and __site__
> * Any description of the sensor location that may be required to ensure data users know exactly what feature the observations are from.


> ### :heavy_check_mark: Delta prerequisites
> * The recording __station__ need to be present in the __network/stations.csv__ file.
> * The recording __site__ need to be present in the __network/sites.csv__ file.

>
> ### :small_orange_diamond: Delta checks
> * For this platform, the only (so far) used sampling rate is one sample per 10 minutes. This is indicated by a sampling rate of -600. The negative sign indicates seconds per sample, rather than the usual samples per second.
> * Only one data stream with a given sampling rate can be defined for a recording __site__ at any given time.

## :two: <a name="specifcexamples"></a>_Specific change examples 

> ### :information_source: Replacing a datalogger
> * The datalogger is replaced, but sensors remain unchanged. This situation might apply if a datalogger has a fault, or a firmware or other change is required that is more easily achieved by pre-preparing a new datalogger in the workshop ahead of visiting the site. The field technician can then simply swap dataloggers while onsite.
> * The new datalogger needs to be added to the __assets/dataloggers.csv__ file.
> * The new datalogger needs to be installed in the __install/dataloggers.csv__ file. This is a two step process.
>   * Uninstall the old datalogger. This involves setting the __End Date__ to a time after the last data recorded by the datalogger.
>   * Install the new datalogger. This involves adding an entry to the file for the new datalogger, and setting the __Start Date__ to a time before the first data recorded by the datalogger. The __End Date__ should be set to _9999-01-01T00:00:00Z_.
>   * Ensure that the __End Date__ of the old datalogger and the __Start Date__ of the new datalogger do not overlap.
> * The new datalogger's firmware version needs to be added to __install/firmware.csv__ file.

