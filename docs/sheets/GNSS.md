# COMMON GNSS TASKS

* :scroll: :one: [Creating or updating a GNSS __mark__](#mark)
* :scroll: :two: [Creating or updating a GNSS __monument__](#monument)
* :scroll: :three: [Installing or updating a GNSS __antenna__](#antenna)
* :scroll: :four: [Installing or updating a GNSS __receiver__](#receiver)
* :scroll: :five: [Installing or updating a GNSS __receiver__ firmware version](#firmware)
* :scroll: :six: [Installing or updating a GNSS __receiver__ session configuration](#session) :dragon:
* :hammer: [Combined editing steps](#steps)

Global Navigation Satellite System (GNSS) equipment are installed at GNSS __marks__. Each __mark__ is
attached to a __monument__. The details of the __mark__, as found in the __network/marks.csv__ file,
tend to relate more to documentation, such as name and location. Whereas the __network/monuments.csv__ file
tends to list physical aspects of the actual monument as constructed.

Equipment installed should be found in the the relevant __assets__ directory files. Actual times when
__antennas__ and __receivers__ are deployed is managed via the __install/antennas.csv__ and __install/receiver.csv__
files respectively. As the firmware of the __receiver__ is important for data post-processing, this is maintained
in the __install/firmware.csv__ file. This is independent of installation and can also be maintained for instruments
not actively deployed.

Describing how the data is being recorded is managed through the __install/sessions__ file. This correlates with the general model
of the receiver being deployed and the __session__ sample "intervals". Other meta-data information required for post-processing
or analysis is also stored in this file. These settings can be delimited by time ranges and are not dependent on
any particular installed receiver but only on the receiver model type.

## :one: <a name="mark"></a>_Creating or updating a GNSS mark_

For new marks a __code__ will need to be assigned by the appropriate mechanism.

> ### :page_with_curl: Files to update
>
> * __network/marks.csv__
>
> ### :page_with_curl: Reference files
>
> * __network/networks.csv__
>
> ### :information_source: General requirements
> 
> * The GNSS __mark__ "code".
> * The "name" of the GNSS __mark__.
> * The "code" for the __network__ that the __mark__ is a member of.
> * Whether the __mark__ is an "_IGS_" site.
>
> ### :information_source: Field requirements
>
> * The location of the __mark__ to create or update.
> * The operational time range of the __mark__.
> 
> ### :heavy_check_mark: Delta prerequisites
> 
> * There can only be one __mark__ for any given code.
> * There must be an entry in the __networks.csv__ file for the __mark's__ __network__ code.
>
> ### :small_orange_diamond: Delta checks
>
> * The "_IGS_" entry is either "__yes__" or "__no__".

## :two: <a name="monument"></a>_Creating or updating a GNSS monument_

> ### :page_with_curl: Files to update
>
> * __network/monuments.csv__
>
> ### :page_with_curl: Reference files
>
> * __network/marks.csv__
>
> ### :information_source: General requirements
> 
> * The "code" of the __mark__ that the monument is associated with.
> * The "domes number" of the __monument__.
>
> ### :information_source: Field requirements
>
> * The "mark type" of the __monument__.
> * The ground relationship of the __monument__.
> * The __monument__ "foundation type".
> * The __monument__ "foundation depth".
> * The operational time range of the __monument__.
> 
> ### :heavy_check_mark: Delta prerequisites
> 
> * There can be only one __mark__ attached to a __monument__.
> * Only one __monument__ can be associated with any given __mark__.
>
> ### :small_orange_diamond: Delta checks
>
> * The "mark type" of the __monument__ must be one of:
> * * "Shallow Rod / Braced Antenna Mount"
> * * "Wyatt/Agnew Drilled-Braced"
> * * "Pillar"
> * * "Steel Mast"
> * * "Unknown"
> * The "ground relationship" must be equal to or less than zero.

## :three: <a name="antenna"></a>_Installing or replacing a GNSS antenna_

> ### :page_with_curl: Files to update
>
> * __install/antennas.csv__
>
> ### :page_with_curl: Reference files
>
> * __assets/antennas.csv__
> * __network/marks.csv__
>
> ### :information_source: General requirements
>
> * GNSS __mark__ code (_where the antenna is mounted_)
>
> ### :information_source: Field requirements
> 
> * Antenna make, model, and serial number.
> * Installation details, such as height, offset from the mark, and azimuth.
> * Antenna installation or removal times.
> 
> ### :heavy_check_mark: Delta prerequisites
> 
> * The antennas need to be listed in the __assets/antennas.csv__ file.
> * The associated __mark__ needs to be listed in the __network/marks.csv__ file.
> 
> ### :small_orange_diamond: Delta checks
> 
> * A given __antenna__ can only be installed at a single __mark__ at any given time.
> * A __mark__ can only have one __antenna__ mounted at any given time.

## :four: <a name="receiver"></a>_Installing or replacing a GNSS receiver_

> ### :page_with_curl: Files to update
>
> * __install/receivers.csv__
>
> ### :page_with_curl: Reference files
>
> * __assets/receivers.csv__
> * __network/marks.csv__
>
> ### :information_source: General requirements
>
> * GNSS __mark__ code (_where the antenna is mounted_)
>
> ### :information_source: Field requirements
> 
> * Receiver make, model, and serial number.
> * Receiver deployment or removal times.
> 
> ### :heavy_check_mark: Delta prerequisites
> 
> * Receivers need to be listed in the __assets/receivers.csv__ file.
> * The associated __mark__ needs to be listed in the __network/marks.csv__ file.
> 
> ### :small_orange_diamond: Delta checks
> 
> * A given __receiver__ can only be deployed to a single __mark__ at any given time.
> * A __mark__ can only have one __receiver__ mounted at any given time.
>

## :five: <a name="firmware"></a>_Updating or adding a GNSS receiver firmware version_

> ### :page_with_curl: Files to update
>
> * __install/firmware.csv__
>
> ### :page_with_curl: Reference files
>
> * __assets/receivers.csv__
>
> ### :information_source: Field requirements
> 
> * Receiver make, model, and serial number.
> * Receiver "firmware" version number.
> * Receiver "firmware" version installation times.
> * Any extra public "firmware" version notes.
> 
> ### :heavy_check_mark: Delta prerequisites
>
> * The __receiver__ needs to be matched in the __assets/receiver.csv__ file.
>
> ### :small_orange_diamond: Delta checks
>
> * A given __receiver__ can only have one __firmware__ version at any given time.

## :six: <a name="session"></a>_Updating or adding a GNSS receiver session configuration_

> ### :page_with_curl: Files to update
>
> * __install/sessions.csv__
>
> ### :page_with_curl: Reference files
>
> * __network/marks.csv__
>
> ### :information_source: General requirements
>
> * GNSS session __mark__ "code".
> * GNSS meta-data __mark__ "operator" and "agency".
> * GNSS meta-data "header comment" reference.
> * GNSS meta-data "header format" look up.
>
> ### :information_source: Field requirements
>
> * GNSS __receiver__ "satellite system" configured.
> * GNSS __receiver__ session "interval".
> * GNSS __receiver__ configured "elevation mask".
> * The time range over which the session is active.
> 
> ### :heavy_check_mark: Delta prerequisites
>
> * The associated __mark__ needs to be listed in the __network/marks.csv__ file.
>
> ### :small_orange_diamond: Delta checks
>
> * There can only be one __session__ per __mark__ with the same interval at any given time.
> * The __session__ start times need to be before the associated end times.
> * Session "satellite systems" need to be one of:
> * * "GPS"
> * * "GPS+GLO"
> * * "GPS+GLO+GAL+BDS+QZSS"
>

## <a name="steps"></a>_Overall steps_

> * :file_folder: Using a suitable mechanism create a new _git_ branch for the changes.
> * :one: :pencil2: Update the __network/marks.csv__ file to add or update any GNSS __mark__ information, the entries will need to be in order of __mark__ "code".
> * :two: :pencil2: Update the __network/monuments.csv__ file to add or update any GNSS __monument__ information, the entries will need to be in order of associated __mark__ "code".
> * :three: :pencil2: Update the __install/antenna.csv__ file to add or update any installed antennas, the entries will need to be in order of make, model, serial number, and installation time.
> * :four: :pencil2: Edit the __install/receivers.csv__ file to add or update any deployed receivers, the entries will need to be in order of make, model, serial number, and deployment time.
> * :five: :pencil2: Update the __install/firmware.csv__ file to add or adjust the deployed __receiver__ "firmware" versions, the entries will need to be in order of make, model, serial number, and firmware start times.
> * :six: :pencil2: Update the __install/sessions.csv__ file to add or adjust the deployed __receiver__ __session__ settings, the entries will need to be in order of __mark__ "code", session "interval", and session start times.
> * :open_file_folder: Build the pull request with a meaningful title.
> * :link: Assign suitable tags and set reviewers.
> * :repeat: If the tests fail, the above changes may need some iteration until they pass.
> * :sos: If the tests are still failing, escalate as this may indicate some inconsistency within the network configuration.
> * :ok: Once the tests have passed and the pull request reviewed, depending on policy, the pull request can be merged
