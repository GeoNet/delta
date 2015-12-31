package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Site struct {
	Point

	StationCode  string
	LocationCode string
}

func (s Site) Less(site Site) bool {
	switch {
	case s.StationCode < site.StationCode:
		return true
	case s.StationCode > site.StationCode:
		return false
	case s.LocationCode < site.LocationCode:
		return true
	case s.LocationCode > site.LocationCode:
		return false
	default:
		return false
	}
}

type Sites []Site

func (s Sites) Len() int           { return len(s) }
func (s Sites) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Sites) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s Sites) encode() [][]string {
	data := [][]string{{
		"Station Code",
		"Location Code",
		"Latitude",
		"Longitude",
		"Height",
		"Datum",
	}}
	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.StationCode),
			strings.TrimSpace(v.LocationCode),
			strconv.FormatFloat(v.Latitude, 'g', -1, 64),
			strconv.FormatFloat(v.Longitude, 'g', -1, 64),
			strconv.FormatFloat(v.Elevation, 'g', -1, 64),
			strings.TrimSpace(v.Datum),
		})
	}
	return data
}

func (s *Sites) decode(data [][]string) error {
	var sites []Site
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 6 {
				return fmt.Errorf("incorrect number of installed site fields")
			}
			var err error

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[2], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[3], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[4], 64); err != nil {
				return err
			}

			sites = append(sites, Site{
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[5]),
				},
				StationCode:  strings.TrimSpace(d[0]),
				LocationCode: strings.TrimSpace(d[1]),
			})
		}

		*s = Sites(sites)
	}
	return nil
}

func LoadSites(path string) ([]Site, error) {
	var s []Site

	if err := LoadList(path, (*Sites)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(Sites(s))

	return s, nil
}
