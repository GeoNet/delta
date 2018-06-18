# COMMON CAMERA TASKS

* :scroll: :one: [Creating or updating a camera __mount__](#mount) :dragon:
* :scroll: :two: [Installing or updating a __camera__](#camera)
* :hammer: [Combined editing steps](#steps)

The mounting of cameras used for monitoring operations are described in the __network/mounts.csv__ file.
Attached to each __mount__ are a camera "code", the physical location, a "name" (or caption), and a description of the overall "view".

Physical camera installs are managed through the __install/cameras.csv__ file. This includes
camera hardware details and the installation times.

## :one: <a name="mark"></a>_Creating or updating an existing camera mount_

> ### :page_with_curl: Files to update
>
> * __network/mounts.csv__
>
> ### :page_with_curl: Reference files
>
> * __network/networks.csv__ [_currently not populated_]
>

New camera __mounts__ will need to have a "code" assigned by the appropriate mechanism.

It is intended that camera mounts are fixed in place, any adjustments will be applied retrospectively.
Changes over time of a camera's location, or view point, should be handled by building new camera __mounts__.

:dragon: **The image "name", and general "view" description are prominently used on the GeoNet web pages.**

> ### :information_source: General requirements
> 
> * The camera __mount__ "code".
> * The expected image "name" or caption.
> * The overall "view" description.
> 
> ### :information_source: Field requirements
> 
> * The camera mount location.
>
> ### :heavy_check_mark: Delta prerequisites
>
> * There can only be one camera __mount__ for any given code.

## :two: <a name="camera"></a>_Installing or replacing cameras_

> ### :page_with_curl: Files to update
>
> * __install/cameras.csv__
>
> ### :page_with_curl: Reference files
>
> * __assets/cameras.csv__
> * __network/mounts.csv__
>
> ### :information_source: General requirements
>
> * Camera __mount__ "code" (_where the camera is mounted_).
> * Any publicly readable notes relating to the installation.
>
> ### :information_source: Field requirements
> 
> * Camera "make", "model", and "serial" numbers (_what is being installed or removed_)
> * Physical installation details, such as dip, azimuth (_direction the camera is pointing_), and height.
> * Camera installation times (_when the camera is installed, or removed_)
> 
> If replacing an existing camera the installation details are likely to be the same,
> although a cross-check may be worthwhile.
>
> Use an end date of `9999-01-01T00:00:00Z` to indicate that a camera is currently installed.
> 
> ### :heavy_check_mark: Delta prerequisites
> 
> * The cameras need to be listed in the __assets/cameras.csv__ file.
> * The camera mounts need to be present in the __network/mounts.csv__ file.
> 
> ### :small_orange_diamond: Delta checks
> 
> * A given camera can only be mounted to a single __mount__ at any given time.
> * A camera __mount__ can only have one camera mounted at any given time.

## <a name="steps"></a>_Overall steps_

> * :file_folder: Using a suitable mechanism create a new _git_ branch for the changes.
> * :one: :pencil2: Update the __network/mounts.csv__ file to add or update any camera mount information, the entries will need to be in order of mount code.
> * :two: :pencil2: Update in the __install/cameras.csv__ file any equipment that has been removed or altered.
> * :open_file_folder: Build the pull request with a meaningful title.
> * :link: Assign suitable tags and set reviewers.
> * :repeat: If the tests fail, the above changes may need some iteration until they pass.
> * :sos: If the tests are still failing, escalate as this may indicate some inconsistency within the network configuration.
> * :ok: Once the tests have passed and the pull request reviewed, depending on policy, the pull request can be merged
