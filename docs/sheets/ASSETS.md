# COMMON ASSET MANAGEMENT TASKS

Although not strictly needed, the __asset__ files are used as a cross-check for when equipment
is installed. This is to try and catch errors related to serial numbers, and also to allow
using the asset numbers to further identify installed equipment. 

The files in the __assets__ directory all have the same format. Each file represents an equipment
type which also constrains where the equipment can be installed.
For instance, only equipment in the __assets/sensors.csv__ file can be installed in the __install/sensors.csv__ file.

Generally __assets__ should be entered in before they are added to files in the __installs__ directory.

The presence of an asset "number" is an indication that the equipment is managed directly by __GeoNet__.

[_The files were not built for asset tracking or inventory management, however this may change with further requirements_]

## Files

---
* :page_with_curl: __assets/antennas.csv__
* :page_with_curl: __assets/metsensors.csv__
* :page_with_curl: __assets/radomes.csv__
* :page_with_curl: __assets/receivers.csv__
---
* :page_with_curl: __assets/sensors.csv__
* :page_with_curl: __assets/recorders.csv__
* :page_with_curl: __assets/dataloggers.csv__
---
* :page_with_curl: __assets/cameras.csv__
---

## :one: _Updating an existing asset entry_

Find the appropriate __asset__ file for the type of equipment.

> ### :information_source: General requirements
> 
> * If changed, the make, model, or serial number of the assets to update.
> * If changed, the equipment asset "numbers".
> * If changed, any equipment notes for the updated assets.

> ### :small_orange_diamond: Delta checks
> 
> * There can only be one asset with the same make, model, and serial number.
> * There can only be one asset with the same asset "number".
> * Any references to these assets in the __install__ directory files will also need to be consistent.

## :two: _Adding a new asset entry_

Find the appropriate __asset__ file for the type of equipment.

> ### :information_source: General requirements
> 
> * The make, model, and serial number of the new `assets.
> * Any optional asset "numbers" that may be available.
> * Any extra equipment notes that may be appropriate.
>
> ### :small_orange_diamond: Delta checks
> 
> * There can only be one asset with the same make, model, and serial number.
> * There can only be one asset with the same asset "number".

## _Overall steps_

> * :file_folder: Using a suitable mechanism create a new _git_ branch for the changes.
> * :one: :pencil2: Update the __assets/&lt;type&gt;.csv__ files appropriate to the equipment model types.
> * :one: :pencil2: Update the __install/&lt;type&gt;.csv__ files where the equipment may be referenced.
> * :two: :pencil2: Add the new equipment into the __assets/&lt;type&gt;.csv__ files appropriate to the model types.
> * :open_file_folder: Build the pull request with a meaningful title.
> * :link: Assign suitable tags and set reviewers.
> * :repeat: If the tests fail, the above changes may need some iteration until they pass.
> * :sos: If the tests are still failing, escalate as this may indicate some inconsistency within the network configuration.
> * :ok: Once the tests have passed and the pull request reviewed, depending on policy, the pull request can be merged
