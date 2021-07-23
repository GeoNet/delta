package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	locationCode = iota
	locationDescription
	locationLatitude
	locationLongitude
	locationElevation
	locationDatum
	locationHeight
	locationLast
)

type Location struct {
	Point

	Code        string
	Description string
	Height      float64

	height string
}

type LocationList []Location

func (s LocationList) Len() int      { return len(s) }
func (s LocationList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s LocationList) Less(i, j int) bool {
	switch {
	case s[i].Code < s[j].Code:
		return true
	case s[i].Code > s[j].Code:
		return false
	default:
		return false
	}
}

func (s LocationList) encode() [][]string {
	data := [][]string{{
		"Code",
		"Description",
		"Latitude",
		"Longitude",
		"Elevation",
		"Datum",
		"Height",
	}}

	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Code),
			strings.TrimSpace(v.Description),
			strings.TrimSpace(v.Point.latitude),
			strings.TrimSpace(v.Point.longitude),
			strings.TrimSpace(v.Point.elevation),
			strings.TrimSpace(v.Datum),
			strings.TrimSpace(v.height),
		})
	}
	return data
}
func (s *LocationList) decode(data [][]string) error {
	var locations []Location
	if len(data) > 1 {
		for _, d := range data[1:] {
			if n := len(d); n != locationLast {
				return fmt.Errorf("incorrect number of location fields (got %d, expected %d)", n, locationLast)
			}
			var err error

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[locationLatitude], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[locationLongitude], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[locationElevation], 64); err != nil {
				return err
			}

			var height float64
			if v := d[locationHeight]; v != "" {
				if height, err = strconv.ParseFloat(v, 64); err != nil {
					return err
				}
			}

			locations = append(locations, Location{
				Point: Point{
					// geographic details
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[locationDatum]),

					// shadow variables
					latitude:  d[locationLatitude],
					longitude: d[locationLongitude],
					elevation: d[locationElevation],
				},

				Code:        strings.TrimSpace(d[locationCode]),
				Description: strings.TrimSpace(d[locationDescription]),
				Height:      height,

				height: d[locationHeight],
			})
		}

		*s = LocationList(locations)
	}
	return nil
}

func LoadLocations(path string) ([]Location, error) {
	var s []Location

	if err := LoadList(path, (*LocationList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(LocationList(s))

	return s, nil
}
