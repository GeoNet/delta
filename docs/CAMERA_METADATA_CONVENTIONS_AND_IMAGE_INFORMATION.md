# Camera Metadata Conventions

A camera site name consists of two components, a mount code, and a view code.

## Mount Codes

Mount codes are the equivalent of site codes for seismic and acoustic sensors. Mount codes are assigned at the first installation of a camera at a location and do not change.

Mount codes are unique across all sites, including seismic, GNSS, and other sites. They are found in the [network/mounts.csv](https://github.com/GeoNet/delta/blob/main/network/mounts.csv) file.

Mount codes are named for ease of distinction and reflect either (or both) camera location or monitoring subject cues.

Mount codes are 4 or more letters.

For some older camera sites, the first two letters of the mount code give an indication of where the site is (i.e. they will be an abbreviation of a significant local geographic feature), and the last two letters give a sub-location, particularly where there are multiple mounts in the same area, e.g. WINR is at White Island (WI) on the north rim of the crater (NR). In other old camera sites, the mount codes do not follow this pattern, typically for historical or other reasons, e.g. TEMO is Taranaki Emergency Management Office, and KAKA is Kakaramea.

Newer mount codes either abbreviate the location prefix to one letter or discard it completely, e.g., WSMC (White Island Summit Camera).

As of 2026, permanent camera sites are named with a C suffix to indicate their nature. This allows distinction of camera sites from GNSS at the same location where they would otherwise have equivalent station codes.

## View Codes

Cameras are installed under view codes at the mount. View codes are the equivalent of location codes for seismic and acoustic sensors.

View codes have two digits. They are unique for a given mount code. The initial view code assigned is usually 01. Later view codes are 02, 03, etc. They are found in the [network/views.csv](https://github.com/GeoNet/delta/blob/main/network/views.csv) file.

The view code is used for one of two purposes:
- To distinguish between multiple cameras installed at a single observation site where the same mount code is used.
- To show that the view of a particular camera has changed substantially so that the subject of the view does not have the same prominence or position in the image. This is intended to reflect an intentional change in view, but is also used when there is an unintentional change if that change is not subsequently corrected.

## Naming Conventions When Moving or Changing Cameras

There are situations when a mount code and view code should be changed following the movement, changing, or addition of a camera. In part this depends on what is considered the subject of the camera - the part of the image view where visual changes of interest are expected to occur. This implicitly requires the subject for each camera to be defined. Mount or view codes change will be made:
1. If the place where a camera is mounted (e.g. the building or pole on which the camera is attached) is changed, a new mount code will be made.
1. If the camera's view is changed intentionally (e.g. rotated to focus on volcano B rather than volcano A), but the mount code does not change, then a new view code should be made.
1. If the camera's view is changed unintentionally, but the new view persists and is not reverted to it's former view, then a new view code should be made.
1. If another camera is added to an existing mount, then that camera has a different view code from the existing camera.

These conventions are ultimately used at the discretion of the person(s) responsible for the equipment change, who should use expert judgement and consider the purpose and use of the data. Consistency in applying these conventions is important.

### Camera and Lens Changes

If there is a change in the camera installed at a mount, either a like-for-like swap, or a change of camera make or model, this does not result in a change in view code, unless one of the conditions for a view code change is also satisfied.

Some cameras have a different lens for daytime and nightime views. Images from both lenses use the same view code.

# Image File Naming Conventions

The camera image file naming conventions describe how the camera details and time an image is taken are used to generate a unique name for each image.

Each image file is identified by set of codes:

__YEAR.DOY.HOURMIN.SEC.MOUNTCODE.VIEWCODE.jpg__

e.g. 2021.090.1915.44.WHOH.01.jpg

- The image time is measured in Universal Time (UTC), which is 12 hours behind New Zealand Standard Time (NZST) and 13 hours behind New Zealand Daylight Time (NZDT).
- YEAR is a four digit number, e.g. 2021.
- DOY represents the day of the year, as a three digit number. This starts at 001 for January 1, increments by one for every day, ending at 365 for December 31, e.g. 090 (for March 30).
- HOURMIN is a four digit representation of the time the image was taken, e.g. 1915. HOUR has two digits, ranging from 00 to 23, and MIN has two digits, ranging from 00 to 59.
- SEC is a two digit representation of the seconds when the image was taken. e.g. 44.
- MOUNTCODE is a four letter code for the mount position, e.g. WHOH.
- VIEWCODE is a two digit code for the view, e.g. 01.
- jpg is the standard representation of the image format [JPEG](https://en.wikipedia.org/wiki/JPEG).

# Detailed Image Information

To perform analysis of an image, a user sometimes requires knowledge of details such as the make and model of the camera taking the image, and technical details such as the angle of view of the lens, the lens focal length, camera aperture f-stop, shutter speed, etc.

## Camera Make and Model

The make and model of the camera at a particular mount and view is listed in the [install/cameras.csv](https://github.com/GeoNet/delta/blob/main/install/cameras.csv) file. This includes the azimuth (direction of view) and the dip (camera angle relative to the horizontal). Notes provide some additional information about the view of a camera.

## EXIF Information

The only way to access technical details for an image, such as the angle of view of the lens, the lens focal length, camera aperture f-stop, shutter speed, etc, is to view an image's [EXIF information](https://en.wikipedia.rg/wiki/Exif). The extent of EXIF information available for an image varies with the make and model of the camera that took the image.

### Image Processing Software

Most commonly used image processing software can read and show EXIF information.

### Other Tools

There are also tools to display EXIF information that do not require a user to use image processing software.

- The file manager GUI software used in most computer operating systems allows a user to click on an image and select `Properties`, or similar, and see the main EXIF information, though typically not everyting.
- Most common web browsers permit a user to install an extension to view EXIF information for an image viewed from within the browser.
- There are also command line tools to examine EXIF information without using a file manager GUI or a browser.
