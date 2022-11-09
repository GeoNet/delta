package meta

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

const (
	RadiansToDegrees = 57.2957795
	DegreesToKm      = math.Pi * 6371.0 / 180.0
	RadiansToKm      = RadiansToDegrees * DegreesToKm
)

const (
	placenameName = iota
	placenameLatitude
	placenameLongitude
	placenameLevel
	placenameLast
)

// Placename is used to describe distances and azimuths to known places.
type Placename struct {
	Name      string
	Latitude  float64
	Longitude float64
	Level     int

	latitude  string
	longitude string
}

// Distance returns the distance in kilometres from the given latitude and longitude to the Placename.
func (p Placename) Distance(lat, lon float64) float64 {

	if (p.Latitude == lat) && (p.Longitude == lon) {
		return 0.0
	}

	esq := (1.0 - 1.0/298.25) * (1.0 - 1.0/298.25)
	alat3 := math.Atan(math.Tan(p.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees
	alat4 := math.Atan(math.Tan(lat/RadiansToDegrees)*esq) * RadiansToDegrees

	rlat1 := alat3 / RadiansToDegrees
	rlat2 := alat4 / RadiansToDegrees
	rdlon := (lon - p.Longitude) / RadiansToDegrees

	clat1 := math.Cos(rlat1)
	clat2 := math.Cos(rlat2)
	slat1 := math.Sin(rlat1)
	slat2 := math.Sin(rlat2)
	cdlon := math.Cos(rdlon)

	cdel := slat1*slat2 + clat1*clat2*cdlon
	switch {
	case cdel > 1.0:
		cdel = 1.0
	case cdel < -1.0:
		cdel = -1.0
	}

	return RadiansToKm * math.Acos(cdel)
}

// Azimuth returns the azimuth in degrees from the given latitude and longitude to the Placename.
func (p Placename) Azimuth(lat, lon float64) float64 {

	if (p.Latitude == lat) && (p.Longitude == lon) {
		return 0.0
	}

	esq := (1.0 - 1.0/298.25) * (1.0 - 1.0/298.25)
	alat3 := math.Atan(math.Tan(p.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees
	alat4 := math.Atan(math.Tan(lat/RadiansToDegrees)*esq) * RadiansToDegrees

	rlat1 := alat3 / RadiansToDegrees
	rlat2 := alat4 / RadiansToDegrees
	rdlon := (lon - p.Longitude) / RadiansToDegrees

	clat1 := math.Cos(rlat1)
	clat2 := math.Cos(rlat2)
	slat1 := math.Sin(rlat1)
	slat2 := math.Sin(rlat2)
	cdlon := math.Cos(rdlon)
	sdlon := math.Sin(rdlon)

	yazi := sdlon * clat2
	xazi := clat1*slat2 - slat1*clat2*cdlon

	azi := RadiansToDegrees * math.Atan2(yazi, xazi)

	if azi < 0.0 {
		azi += 360.0
	}

	return azi
}

// BackAzimuth returns the back-azimuth in degrees from the given latitude and longitude to the Placename.
func (p Placename) BackAzimuth(lat, lon float64) float64 {

	if (p.Latitude == lat) && (p.Longitude == lon) {
		return 0.0
	}

	esq := (1.0 - 1.0/298.25) * (1.0 - 1.0/298.25)
	alat3 := math.Atan(math.Tan(p.Latitude/RadiansToDegrees)*esq) * RadiansToDegrees
	alat4 := math.Atan(math.Tan(lat/RadiansToDegrees)*esq) * RadiansToDegrees

	rlat1 := alat3 / RadiansToDegrees
	rlat2 := alat4 / RadiansToDegrees
	rdlon := (lon - p.Longitude) / RadiansToDegrees

	clat1 := math.Cos(rlat1)
	clat2 := math.Cos(rlat2)
	slat1 := math.Sin(rlat1)
	slat2 := math.Sin(rlat2)
	cdlon := math.Cos(rdlon)
	sdlon := math.Sin(rdlon)

	ybaz := -sdlon * clat1
	xbaz := clat2*slat1 - slat2*clat1*cdlon

	baz := RadiansToDegrees * math.Atan2(ybaz, xbaz)

	if baz < 0.0 {
		baz += 360.0
	}

	return baz
}

// Compass returns a text representation of the azimuth from the given latitude and longitude to the Placename.
func (p Placename) Compass(lat, lon float64) string {
	azimuth := p.Azimuth(lat, lon) + 22.5

	for azimuth < 0.0 {
		azimuth += 360.0
	}
	for azimuth >= 360.0 {
		azimuth -= 360.0
	}

	switch int(math.Floor(azimuth / 45.0)) {
	case 0:
		return "north"
	case 1:
		return "north-east"
	case 2:
		return "east"
	case 3:
		return "south-east"
	case 4:
		return "south"
	case 5:
		return "south-west"
	case 6:
		return "west"
	default:
		return "north-west"
	}
}

type PlacenameList []Placename

func (p PlacenameList) Len() int      { return len(p) }
func (p PlacenameList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PlacenameList) Less(i, j int) bool {
	return strings.ToLower(p[i].Name) < strings.ToLower(p[j].Name)
}

func (p PlacenameList) encode() [][]string {
	data := [][]string{{
		"Name",
		"Latitude",
		"Longitude",
		"Level",
	}}

	for _, v := range p {
		data = append(data, []string{
			strings.TrimSpace(v.Name),
			strings.TrimSpace(v.latitude),
			strings.TrimSpace(v.longitude),
			strconv.Itoa(v.Level),
		})
	}
	return data
}

func (p *PlacenameList) decode(data [][]string) error {
	var placenames []Placename
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != placenameLast {
				return fmt.Errorf("incorrect number of placename fields")
			}

			latitude, err := strconv.ParseFloat(d[placenameLatitude], 64)
			if err != nil {
				return err
			}

			longitude, err := strconv.ParseFloat(d[placenameLongitude], 64)
			if err != nil {
				return err
			}

			level, err := ParseInt(d[placenameLevel])
			if err != nil {
				return err
			}

			placenames = append(placenames, Placename{
				Name:      strings.TrimSpace(d[placenameName]),
				Latitude:  latitude,
				Longitude: longitude,
				Level:     level,

				latitude:  strings.TrimSpace(d[placenameLatitude]),
				longitude: strings.TrimSpace(d[placenameLongitude]),
			})
		}

		*p = PlacenameList(placenames)
	}
	return nil
}

// Closest returns the Placename which is the closest to the given latitude and longitude taking into
// account the Placename level. The level is used to avoid small places taking precident over larger
// places at longer distances. Currently level three addresses will be used for distances within 20 km,
// level two within 100 km, level one within 500km, and level zero has no distance threshold.
func (p PlacenameList) Closest(lat, lon float64) (Placename, bool) {
	var res Placename

	sort.Sort(p)

	var found bool
	var distance float64

	for _, placename := range p {
		dist := placename.Distance(lat, lon)
		if dist > 20.0 && placename.Level > 2 {
			continue
		}
		if dist > 100.0 && placename.Level > 1 {
			continue
		}
		if dist > 500.0 && placename.Level > 0 {
			continue
		}
		if !found || dist < distance {
			distance, res, found = dist, placename, true
		}
	}

	return res, found
}

// Description returns a string representation of where a point location is relative to the nearest place.
func (p PlacenameList) Description(lat, lon float64) string {

	loc, ok := p.Closest(lat, lon)
	if !ok {
		return ""
	}

	switch dist := loc.Distance(lat, lon); {
	case dist < 5.0:
		return fmt.Sprintf("within 5 km of %s", loc.Name)
	default:
		return fmt.Sprintf("%.0f km %s of %s", dist, loc.Compass(lat, lon), loc.Name)
	}
}

func LoadPlacenames(path string) ([]Placename, error) {
	var s []Placename

	if err := LoadList(path, (*PlacenameList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(PlacenameList(s))

	return s, nil
}
