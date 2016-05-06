package main

var sitelogTemplate string = `     {{.SiteIdentification.FourCharacterID}} Site Information Form (site log)
     International GNSS Service
     See Instructions at:
       ftp://igscb.jpl.nasa.gov/pub/station/general/sitelog_instr.txt 

0.   Form

     Prepared by (full name)  : {{.FormInformation.PreparedBy}}
     Date Prepared            : {{.FormInformation.DatePrepared}}
     Report Type              : {{.FormInformation.ReportType}}
     If Update:
      Previous Site Log       : 
      Modified/Added Sections : 


1.   Site Identification of the GNSS Monument

     Site Name                : {{.SiteIdentification.SiteName}}
     Four Character ID        : {{.SiteIdentification.FourCharacterID}}
     Monument Inscription     : {{.SiteIdentification.MonumentInscription|empty "none"}}
     IERS DOMES Number        : {{.SiteIdentification.IersDOMESNumber|empty "none"}}
     CDP Number               : {{.SiteIdentification.CdpNumber|empty "none"}}
     Monument Description     : {{.SiteIdentification.MonumentDescription|tolower}}
       Height of the Monument : {{.SiteIdentification.HeightOfTheMonument}}
       Monument Foundation    : {{.SiteIdentification.MonumentFoundation|tolower}}
       Foundation Depth       : {{.SiteIdentification.FoundationDepth}}
     Marker Description       : {{.SiteIdentification.MarkerDescription}}
     Date Installed           : {{.SiteIdentification.DateInstalled}}
     Geologic Characteristic  : {{.SiteIdentification.GeologicCharacteristic}}
       Bedrock Type           : {{.SiteIdentification.BedrockType}}
       Bedrock Condition      : {{.SiteIdentification.BedrockCondition}}
       Fracture Spacing       : {{.SiteIdentification.FractureSpacing}}
       Fault zones nearby     : {{.SiteIdentification.FaultZonesNearby}}
         Distance/activity    : {{.SiteIdentification.DistanceActivity}}
     Additional Information   : {{.SiteIdentification.Notes}}


2.   Site Location Information

     City or Town             : {{.SiteLocation.City}}
     State or Province        : {{.SiteLocation.State}}
     Country                  : {{.SiteLocation.Country}}
     Tectonic Plate           : {{.SiteLocation.TectonicPlate}}
     Approximate Position (ITRF)
       X coordinate (m)       : {{.SiteLocation.ApproximatePositionITRF.XCoordinateInMeters}}
       Y coordinate (m)       : {{.SiteLocation.ApproximatePositionITRF.YCoordinateInMeters}}
       Z coordinate (m)       : {{.SiteLocation.ApproximatePositionITRF.ZCoordinateInMeters}}
       Latitude (N is +)      : {{.SiteLocation.ApproximatePositionITRF.LatitudeNorth|lat}}
       Longitude (E is +)     : {{.SiteLocation.ApproximatePositionITRF.LongitudeEast|lon}}
       Elevation (m,ellips.)  : {{.SiteLocation.ApproximatePositionITRF.ElevationMEllips}}
     Additional Information   : {{.SiteLocation.Notes}}


3.   GNSS Receiver Information

{{ range $n, $r := .GnssReceivers}}3.{{plus $n}} Receiver Type            : {{$r.ReceiverType}}
     Satellite System         : {{$r.SatelliteSystem}}
     Serial Number            : {{$r.SerialNumber}}
     Firmware Version         : {{$r.FirmwareVersion}}
     Elevation Cutoff Setting : {{$r.ElevationCutoffSetting}} deg
     Date Installed           : {{$r.DateInstalled}}
     Date Removed             : {{$r.DateRemoved|empty "CCYY-MM-DDThh:mmZ"}}
     Temperature Stabiliz.    : {{$r.TemperatureStabilization|empty "none"}}
     Additional Information   :  {{$r.Notes}}

{{end}}3.x  Receiver Type            : (A20, from rcvr_ant.tab; see instructions)
     Satellite System         : (GPS+GLO+GAL+BDS+QZSS+SBAS)
     Serial Number            : (A20, but note the first A5 is used in SINEX)
     Firmware Version         : (A11)
     Elevation Cutoff Setting : (deg)
     Date Installed           : (CCYY-MM-DDThh:mmZ)
     Date Removed             : (CCYY-MM-DDThh:mmZ)
     Temperature Stabiliz.    : (none or tolerance in degrees C)
     Additional Information   : (multiple lines)


4.   GNSS Antenna Information

{{ range $n, $a := .GnssAntennas}}4.{{plus $n}} Antenna Type             : {{$a.AntennaType|printf "%-16s"}}{{$a.AntennaRadomeType|empty "NONE"}}
     Serial Number            : {{$a.SerialNumber}}
     Antenna Reference Point  : {{$a.AntennaReferencePoint}}
     Marker->ARP Up Ecc. (m)  : {{$a.MarkerArpUpEcc}}  
     Marker->ARP North Ecc(m) : {{$a.MarkerArpNorthEcc}}  
     Marker->ARP East Ecc(m)  : {{$a.MarkerArpEastEcc}}  
     Alignment from True N    : {{$a.AlignmentFromTrueNorth}}
     Antenna Radome Type      : {{$a.AntennaRadomeType}}
     Radome Serial Number     : {{$a.RadomeSerialNumber}}
     Antenna Cable Type       : {{$a.AntennaCableType}}
     Antenna Cable Length     :{{$a.AntennaCableLength}}
     Date Installed           : {{$a.DateInstalled}}
     Date Removed             : {{$a.DateRemoved|empty "CCYY-MM-DDThh:mmZ"}}
     Additional Information   : {{$a.Notes}}

{{end}}4.x  Antenna Type             : (A20, from rcvr_ant.tab; see instructions)
     Serial Number            : (A*, but note the first A5 is used in SINEX)
     Antenna Reference Point  : (BPA/BCR/XXX from "antenna.gra"; see instr.)
     Marker->ARP Up Ecc. (m)  : (F8.4)
     Marker->ARP North Ecc(m) : (F8.4)
     Marker->ARP East Ecc(m)  : (F8.4)
     Alignment from True N    : (deg; + is clockwise/east)
     Antenna Radome Type      : (A4 from rcvr_ant.tab; see instructions)
     Radome Serial Number     : 
     Antenna Cable Type       : (vendor & type number)
     Antenna Cable Length     : (m)
     Date Installed           : (CCYY-MM-DDThh:mmZ)
     Date Removed             : (CCYY-MM-DDThh:mmZ)
     Additional Information   : (multiple lines)

5.   Surveyed Local Ties

5.x  Tied Marker Name         : 
     Tied Marker Usage        : (SLR/VLBI/LOCAL CONTROL/FOOTPRINT/etc)
     Tied Marker CDP Number   : (A4)
     Tied Marker DOMES Number : (A9)
     Differential Components from GNSS Marker to the tied monument (ITRS)
       dx (m)                 : (m)
       dy (m)                 : (m)
       dz (m)                 : (m)
     Accuracy (mm)            : (mm)
     Survey method            : (GPS CAMPAIGN/TRILATERATION/TRIANGULATION/etc)
     Date Measured            : (CCYY-MM-DDThh:mmZ)
     Additional Information   : (multiple lines)


6.   Frequency Standard

6.x  Standard Type            : (INTERNAL or EXTERNAL H-MASER/CESIUM/etc)
       Input Frequency        : (if external)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Notes                  : (multiple lines)


7.   Collocation Information

7.x  Instrumentation Type     : (GPS/GLONASS/DORIS/PRARE/SLR/VLBI/TIME/etc)
       Status                 : (PERMANENT/MOBILE)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Notes                  : (multiple lines)


8.   Meteorological Instrumentation

{{ range $n, $m := .GnssMetSensors}}8.1.{{plus $n}}Humidity Sensor Model   : 
       Manufacturer           : 
       Serial Number          : {{$m.SerialNumber}}
       Data Sampling Interval : (sec)
       Accuracy (% rel h)     : (% rel h)
       Aspiration             : (UNASPIRATED/NATURAL/FAN/etc)
       Height Diff to Ant     : (m)
       Calibration date       : (CCYY-MM-DD)
       Effective Dates        : {{$m.EffectiveDates}}
       Notes                  : {{$m.Notes}}
{{end}}

8.1.x Humidity Sensor Model   : 
       Manufacturer           : 
       Serial Number          : 
       Data Sampling Interval : (sec)
       Accuracy (% rel h)     : (% rel h)
       Aspiration             : (UNASPIRATED/NATURAL/FAN/etc)
       Height Diff to Ant     : (m)
       Calibration date       : (CCYY-MM-DD)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Notes                  : (multiple lines)

{{ range $n, $m := .GnssMetSensors}}8.2.{{plus $n}}Pressure Sensor Model   : 
       Manufacturer           : {{$m.Manufacturer}}
       Serial Number          : {{$m.SerialNumber}}
       Data Sampling Interval : (sec)
       Accuracy               : (hPa)
       Height Diff to Ant     : (m)
       Calibration date       : (CCYY-MM-DD)
       Effective Dates        : {{$m.EffectiveDates}}
       Notes                  : {{$m.Notes}}
{{end}}

8.2.x Pressure Sensor Model   : 
       Manufacturer           : 
       Serial Number          : 
       Data Sampling Interval : (sec)
       Accuracy               : (hPa)
       Height Diff to Ant     : (m)
       Calibration date       : (CCYY-MM-DD)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Notes                  : (multiple lines)

{{ range $n, $m := .GnssMetSensors}}8.3.{{plus $n}}Temp. Sensor Model      : 
       Manufacturer           : {{$m.Manufacturer}}
       Serial Number          : {{$m.SerialNumber}}
       Data Sampling Interval : (sec)
       Accuracy               : (deg C)
       Aspiration             : (UNASPIRATED/NATURAL/FAN/etc)
       Height Diff to Ant     : (m)
       Calibration date       : (CCYY-MM-DD)
       Effective Dates        : {{$m.EffectiveDates}}
       Notes                  : {{$m.Notes}}
{{end}}

8.3.x Temp. Sensor Model      : 
       Manufacturer           : 
       Serial Number          : 
       Data Sampling Interval : (sec)
       Accuracy               : (deg C)
       Aspiration             : (UNASPIRATED/NATURAL/FAN/etc)
       Height Diff to Ant     : (m)
       Calibration date       : (CCYY-MM-DD)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Notes                  : (multiple lines)

8.4.x Water Vapor Radiometer  : 
       Manufacturer           : 
       Serial Number          : 
       Distance to Antenna    : (m)
       Height Diff to Ant     : (m)
       Calibration date       : (CCYY-MM-DD)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Notes                  : (multiple lines)

8.5.x Other Instrumentation   : (multiple lines)


9.  Local Ongoing Conditions Possibly Affecting Computed Position

9.1.x Radio Interferences     : (TV/CELL PHONE ANTENNA/RADAR/etc)
       Observed Degradations  : (SN RATIO/DATA GAPS/etc)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Additional Information : (multiple lines)

9.2.x Multipath Sources       : (METAL ROOF/DOME/VLBI ANTENNA/etc)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Additional Information : (multiple lines)

9.3.x Signal Obstructions     : (TREES/BUILDINGS/etc)
       Effective Dates        : (CCYY-MM-DD/CCYY-MM-DD)
       Additional Information : (multiple lines)

10.  Local Episodic Effects Possibly Affecting Data Quality

10.x Date                     : (CCYY-MM-DD/CCYY-MM-DD)
     Event                    : (TREE CLEARING/CONSTRUCTION/etc)

11.   On-Site, Point of Contact Agency Information

     Agency                   : {{.ContactAgency.Agency}}
     Preferred Abbreviation   : {{.ContactAgency.PreferredAbbreviation}}
     Mailing Address          : {{.ContactAgency.MailingAddress|lines "                              : "}}
     Primary Contact
       Contact Name           : {{.ContactAgency.PrimaryContact.Name}}
       Telephone (primary)    : {{.ContactAgency.PrimaryContact.TelephonePrimary}}
       Telephone (secondary)  : {{.ContactAgency.PrimaryContact.TelephoneSecondary}}
       Fax                    : {{.ContactAgency.PrimaryContact.Fax}}
       E-mail                 : {{.ContactAgency.PrimaryContact.Email}}
     Secondary Contact
       Contact Name           : {{.ContactAgency.SecondaryContact.Name}}
       Telephone (primary)    : {{.ContactAgency.SecondaryContact.TelephonePrimary}}
       Telephone (secondary)  : {{.ContactAgency.SecondaryContact.TelephoneSecondary}}
       Fax                    : {{.ContactAgency.SecondaryContact.Fax}}
       E-mail                 : {{.ContactAgency.SecondaryContact.Email}}
     Additional Information   : {{.ContactAgency.Notes}}


12.  Responsible Agency (if different from 11.)

     Agency                   : {{.ResponsibleAgency.Agency}}
     Preferred Abbreviation   : {{.ResponsibleAgency.PreferredAbbreviation}}
     Mailing Address          : {{.ResponsibleAgency.MailingAddress|lines "                              : "}}
     Primary Contact
       Contact Name           : {{.ResponsibleAgency.PrimaryContact.Name}}
       Telephone (primary)    : {{.ResponsibleAgency.PrimaryContact.TelephonePrimary}}
       Telephone (secondary)  : {{.ResponsibleAgency.PrimaryContact.TelephoneSecondary}}
       Fax                    : {{.ResponsibleAgency.PrimaryContact.Fax}}
       E-mail                 : {{.ResponsibleAgency.PrimaryContact.Email}}
     Secondary Contact
       Contact Name           : {{.ResponsibleAgency.SecondaryContact.Name}}
       Telephone (primary)    : {{.ResponsibleAgency.SecondaryContact.TelephonePrimary}}
       Telephone (secondary)  : {{.ResponsibleAgency.SecondaryContact.TelephoneSecondary}}
       Fax                    : {{.ResponsibleAgency.SecondaryContact.Fax}}
       E-mail                 : {{.ResponsibleAgency.SecondaryContact.Email}}
     Additional Information   : {{.ResponsibleAgency.Notes}}

13.  More Information

     Primary Data Center      : {{.MoreInformation.PrimaryDataCenter}}
     Secondary Data Center    :{{.MoreInformation.SecondaryDataCenter}}
     URL for More Information : {{.MoreInformation.UrlForMoreInformation}}
     Hardcopy on File
       Site Map               : {{.MoreInformation.HardCopyOnFile}}
       Site Diagram           : {{.MoreInformation.SiteMap}}
       Horizon Mask           : {{.MoreInformation.HorizonMask}}
       Monument Description   : {{.MoreInformation.MonumentDescription}}
       Site Pictures          : {{.MoreInformation.SitePictures}}
     Additional Information   : {{.MoreInformation.Notes|lines "                              : "}}
     Antenna Graphics with Dimensions

{{ with .MoreInformation.AntennaGraphicsWithDimensions}}{{.}}{{end}}`
