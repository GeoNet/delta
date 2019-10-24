## ANTENNA CALIBRATIONS ##

The `calibrations/antennas` contains individual GNSS receiver antenna calibration files in ANTEX (ANTenna EXchange) format.

An antenna calibration file provides an accurate measurement of the Phase Centre Variations (PCV) position with respect to the GNSS Antenna Reference Point (ARP). 

Individual GNSS antenna calibration files have been provided by [Geoscience Australia GNSS Antenna Calibration Facility](https://www.ga.gov.au/scientific-topics/positioning-navigation/geodesy/gnss-acf) and are also available from Geoscience Australia [ftp](ftp://ftp.ga.gov.au/geodesy-outgoing/gnss/products/antenna).

A full description of the ANTEX format is provided by the International GNSS Service (IGS) at [this link](https://kb.igs.org/hc/en-us/articles/216104678-ANTEX-format-description).

Individual antenna calibrations can be used while processing raw GNSS data to improve the position accuracy of continuous GNSS marks.

Files in this directory are named based on antenna code, radome code, antenna serial number, as prescribed by IGS guidelines.
