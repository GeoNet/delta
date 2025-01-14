package main

import (
	"fmt"

	"github.com/GeoNet/Golang-Ellipsoid/ellipsoid"
	"github.com/GeoNet/delta/meta"
)

var geo = ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)

// genetrate skeleton file for a site,
// based on the receiver, firmware, and antenna metadata at a given time (e.g. now)
func skeleton(code string, country string, set *meta.Set, ts int64) (content string, err error) {
	var mark meta.Mark

	defer func() {
		if err != nil {
			// if any error occurs, compose generic header + obstypes + distribution comment as the result
			// however, still keep the error
			content = fmt.Sprintf(genericHeaderFormat,
				fmt.Sprintf("%s00%s", code, country), // MARKER NAME
			)
			content += obsTypesFour // generic headers use full obsTypesBody
			if mark.Network == "LI" || mark.Network == "GT" {
				content += linzComment
			} else {
				content += geonetComment
			}
		}
	}()

	mark, ok := set.Mark(code)
	if !ok {
		err = fmt.Errorf("no mark found for %s", code)
		return
	}
	if !inWindow(ts, mark.Span) {
		err = fmt.Errorf("no valid mark found for this time period for %s", code)
		return
	}

	receivers := set.DeployedReceivers()
	var dr meta.DeployedReceiver
	for _, r := range receivers {
		if r.Mark == code && inWindow(ts, r.Span) {
			dr = r
			break
		}
	}

	if dr.Span.Start.IsZero() { // empty value
		err = fmt.Errorf("no effective deployed receiver version found for this time for %s", mark.Code)
		return
	}

	x, y, z := geo.ToECEF(mark.Latitude, mark.Longitude, mark.Elevation)
	firmwares := set.FirmwareHistory()
	var ifirm meta.FirmwareHistory
	for _, f := range firmwares {
		if dr.Model == f.Model && dr.Serial == f.Serial && inWindow(ts, f.Span) {
			ifirm = f
			break
		}
	}

	if ifirm.Span.Start.IsZero() {
		err = fmt.Errorf("no effective firmware version found for this time for %s", mark.Code)
		return
	}

	antennas := set.InstalledAntennas()
	var ia meta.InstalledAntenna
	for _, a := range antennas {
		if a.Mark == code && inWindow(ts, a.Span) {
			ia = a
			break
		}
	}

	if ia.Span.Start.IsZero() {
		err = fmt.Errorf("no effective installed antenna found for this time for %s", mark.Code)
		return
	}

	radomes := set.InstalledRadomes()
	var rad meta.InstalledRadome
	for _, r := range radomes {
		if r.Mark == code && inWindow(ts, r.Span) {
			rad = r
			break
		}
	}
	// radome is optional, can be nil
	var radome string
	if rad.Span.Start.IsZero() {
		radome = "NONE"
	} else {
		radome = rad.Model
	}

	var monument meta.Monument
	for _, n := range set.Monuments() {
		if n.Mark == code && inWindow(ts, n.Span) {
			monument = n
			break
		}
	}

	var session meta.Session
	for _, s := range set.Sessions() {
		if s.Mark == code && inWindow(ts, s.Span) {
			session = s
			break
		}
	}

	if session.Span.Start.IsZero() {
		err = fmt.Errorf("no effective session found for this time for %s", mark.Code)
		return
	}

	if !monument.Span.Start.IsZero() && monument.DomesNumber != "" {
		content = fmt.Sprintf(headerFormat,
			fmt.Sprintf("%s00%s", code, country), // MARKER NAME
			monument.DomesNumber,                 // MARKER NUMBER
			dr.Serial, dr.Model, ifirm.Version,   // REC # / TYPE / VERS
			ia.Serial, ia.Model, radome, //ANT # / TYPE
			x, y, z, //APPROX POSITION XYZ
			ia.Offset.Vertical, ia.Offset.East, ia.Offset.North) // ANTENNA: DELTA H/E/N
	} else {
		content = fmt.Sprintf(noMarkerNumberHeaderForamt,
			fmt.Sprintf("%s00%s", code, country), // MARKER NAME
			dr.Serial, dr.Model, ifirm.Version,   // REC # / TYPE / VERS
			ia.Serial, ia.Model, radome, //ANT # / TYPE
			x, y, z, //APPROX POSITION XYZ
			ia.Offset.Vertical, ia.Offset.East, ia.Offset.North) // ANTENNA: DELTA H/E/N
	}

	switch session.SatelliteSystem {
	case "GPS+GLO":
		content += obsTypesTwo
	case "GPS+GLO+GAL+BDS+QZSS":
		content += obsTypesFour
	default:
		err = fmt.Errorf("not valid a SatelliteSystem (%s) for %s", session.SatelliteSystem, mark.Code)
		return
	}

	if mark.Network == "LI" || mark.Network == "GT" {
		content += linzComment
	} else {
		content += geonetComment
	}

	return
}

// Note: A skeleton file is consiste of :
//   - header (3 different kinds)
//   - obsTypesBody (all the same)
//   - comment (LINZ or GeoNet)

const headerFormat = `                    OBSERVATION DATA    M                   RINEX VERSION / TYPE
                                                            PGM / RUN BY / DATE
%-60sMARKER NAME
%-60sMARKER NUMBER
%-20s%-20s%-20sREC # / TYPE / VERS
%-20s%-16s%-4s                    ANT # / TYPE
%14.4f%14.4f%14.4f                  APPROX POSITION XYZ
%14.4f%14.4f%14.4f                  ANTENNA: DELTA H/E/N
`

const noMarkerNumberHeaderForamt = `                    OBSERVATION DATA    M                   RINEX VERSION / TYPE
                                                            PGM / RUN BY / DATE
%-60sMARKER NAME
%-20s%-20s%-20sREC # / TYPE / VERS
%-20s%-16s%-4s                    ANT # / TYPE
%14.4f%14.4f%14.4f                  APPROX POSITION XYZ
%14.4f%14.4f%14.4f                  ANTENNA: DELTA H/E/N
`

const genericHeaderFormat = `                    OBSERVATION DATA    M                   RINEX VERSION / TYPE
                                                            PGM / RUN BY / DATE
%-60sMARKER NAME
`
const obsTypesTwo = `GEODETIC                                                    MARKER TYPE
GeoNet              GNS                                     OBSERVER / AGENCY
G    9 C1C C2W C5X L1C L2W L5X S1C S2W S5X                  SYS / # / OBS TYPES
R   12 C1C C1P C2C C2P L1C L1P L2C L2P S1C S1P S2C S2P      SYS / # / OBS TYPES
G                                                           SYS / PHASE SHIFT
R                                                           SYS / PHASE SHIFT
`
const obsTypesFour = `GEODETIC                                                    MARKER TYPE
GeoNet              GNS                                     OBSERVER / AGENCY
G    9 C1C C2W C5X L1C L2W L5X S1C S2W S5X                  SYS / # / OBS TYPES
R   12 C1C C1P C2C C2P L1C L1P L2C L2P S1C S1P S2C S2P      SYS / # / OBS TYPES
E   15 C1X C5X C6X C7X C8X L1X L5X L6X L7X L8X S1X S5X S6X  SYS / # / OBS TYPES
       S7X S8X                                              SYS / # / OBS TYPES
C    6 C2I C7I L2I L7I S2I S7I                              SYS / # / OBS TYPES
J    9 C1C C2X C5X L1C L2X L5X S1C S2X S5X                  SYS / # / OBS TYPES
G                                                           SYS / PHASE SHIFT
R                                                           SYS / PHASE SHIFT
E                                                           SYS / PHASE SHIFT
C                                                           SYS / PHASE SHIFT
J                                                           SYS / PHASE SHIFT
`

const geonetComment = `These data are supplied by GeoNet. GeoNet is core           COMMENT
funded by NHC, LINZ and MBIE and is operated by             COMMENT
GNS Science on behalf of stakeholders and all New           COMMENT
Zealanders. The data policy, disclaimer, licence and        COMMENT
contact information can be found at www.geonet.org.nz       COMMENT
                                                            END OF HEADER
`

const linzComment = `This station is part of the LINZ and GeoNet cGNSS networks. COMMENT
These networks are operated in partnership between Land     COMMENT
Information New Zealand and GNS Science.                    COMMENT
This data is licensed for re-use under the Creative Commons COMMENT
Attribution 4.0 International licence. For more detail      COMMENT
please refer to https://www.linz.govt.nz/linz-copyright     COMMENT
                                                            END OF HEADER
`

func inWindow(t int64, s meta.Span) bool {
	return t >= s.Start.Unix() && t <= s.End.Unix()
}
