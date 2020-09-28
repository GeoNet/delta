# Site Codes

Site codes are assigned at the first installation of a seismic sensor at a location and do not change. Sensors are installed under location codes at the site following the conventions detailed in this document.

Site codes are unique within a given reference set. Weak motion site codes are unique for the global set of site codes managed by the ISC and must be registered with the ISC to safeguard uniqueness. Strong motion site codes must only be unique within the set of site codes used by GeoNet.     


## Weak Motion

Weak motion site codes are assigned according to the network it comes under. There are two types of network: the national network and the regional networks, with the national network comprising broadband seismometers and the regional networks being mostly composed of short period seismometers. National network site codes are 3 letters and end in 'Z'. Regional network sites codes are 4 letters with the last two letters coming from that regional network's code, e.g. RIGZ is at Rimuhau Hill in the Gisborne Network (GZ), and RAHZ is at Arahi in the Hawke's Bay Network (HZ). For both national and regional network site codes, the first two letters try to give an indication of where the site is (i.e. they will be an abbreviation of a close town or farm station name).

## Strong Motion

Strong motion site codes are not restricted by network codes and instead are composed of 4 letters describing where the site is, e.g. LPLS is Lake Paringa, PNGA is Paringa.

There is a legacy of loose site code conventions in strong motion site codes, such as in the Canterbury strong motion network and more broadly in all strong motion site codes with the last letter being 'S'. As of 2020, such loose conventions are no longer in use. 

# Station and Location Code Changes

The conventions GeoNet use for site and location code changes are outlined in this document.

The ISC guidelines state "because of the need for accurate station positions for hypocenter location programs, a new international code is assigned if a station is moved more than one (1) kilometer from the previous location. If the move is less than one (1) kilometer, a new code will be assigned if requested by the operator."

We note that these guidelines are geared towards global seismic networks, and that our network is regional and local in its nature. Where in the global case a sensor change within 1 km does not change derived earthquake solutions much, in the case of networks designed for regional and local-scale earthquake studies, such a change in position without appropriate metadata changes would cause gross error in derived earthquake solutions.

## Stations, Sites, Location Codes

The designation practices of station, site, and location codes are described in the sister document `SEISMIC_CHANNEL_NAMING_CONVENTIONS.md`.  In short, station codes are designed to convey an idea of where the station is located, site codes denote different sensor installation locations at a station, and location codes describe the type(s) of sensor installed at a site.

A station can have many sites, and a site can have many location codes. Ideally, each site will only have one location code associated with it, but this is often not the case. There may be multiple location codes at a site when:
1. There are sensors of multiple types installed at a site, e.g. a weak motion sensor and a strong motion sensor.
1. There are multiple sensors of the same type installed at a site, e.g. two weak motion sensors.

It should be noted that the choice of having multiple location codes at a site is arbitrary: it is equally as possible to have one site per location code as it is to have one site for all location codes. The former case has the advantage of being able to provide the coordinates of each site as those of the given sensor rather than those of one of (or an average of) the sensors at the site and correspondingly is easier to produce high quality metadata for. As such this is the preferred of the two options. However, the legacy of decisions made prior to the creation of clear conventions muddies our metadata in this respect.

## Naming Conventions When Moving or Changing Sensors

There are 3 cases of how site code and location code change at a station following the movement or changing of a sensor:
1. If the sensor changes to one that is not of the same type, a new location code will be made for the sensor installation, e.g. if a short period seismometer was swapped for a strong motion sensor.
1. If the sensor location moves between than 1-200 m, a new location code will be made for the sensor installation.
1. If the sensor location moves more than 200 m, a new station will be established containing the sensor installation in a site and location code at that station.

In all cases, the caveat of whether the new site/location code(s) have a one-to-one or a one-to-many relationship holds. 
 
These conventions are ultimately used at the discretion of the person(s) responsible for the equipment change. If, for example, a sensor moved less than 200 m but the geologic or local site conditions changed substantially, a new station may be established to reflect this change. 