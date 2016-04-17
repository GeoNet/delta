package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	siteStationCode = iota
	siteLocationCode
	siteLatitude
	siteLongitude
	siteHeight
	siteDatum
	siteSurvey
	siteStartTime
	siteEndTime
	siteLast
)

type Site struct {
	Point
	Span

	StationCode  string
	LocationCode string
	Survey       string
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

type SiteList []Site

func (s SiteList) Len() int           { return len(s) }
func (s SiteList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SiteList) Less(i, j int) bool { return s[i].Less(s[j]) }

func (s SiteList) encode() [][]string {
	data := [][]string{{
		"Station Code",
		"Location Code",
		"Latitude",
		"Longitude",
		"Height",
		"Datum",
		"Survey",
		"Start Time",
		"End Time",
	}}
	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.StationCode),
			strings.TrimSpace(v.LocationCode),
			strconv.FormatFloat(v.Latitude, 'g', -1, 64),
			strconv.FormatFloat(v.Longitude, 'g', -1, 64),
			strconv.FormatFloat(v.Elevation, 'g', -1, 64),
			strings.TrimSpace(v.Datum),
			strings.TrimSpace(v.Survey),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (s *SiteList) decode(data [][]string) error {
	var sites []Site
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != siteLast {
				return fmt.Errorf("incorrect number of installed site fields")
			}
			var err error

			var lat, lon, elev float64
			if lat, err = strconv.ParseFloat(d[siteLatitude], 64); err != nil {
				return err
			}
			if lon, err = strconv.ParseFloat(d[siteLongitude], 64); err != nil {
				return err
			}
			if elev, err = strconv.ParseFloat(d[siteHeight], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[siteStartTime]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[siteEndTime]); err != nil {
				return err
			}

			sites = append(sites, Site{
				Point: Point{
					Latitude:  lat,
					Longitude: lon,
					Elevation: elev,
					Datum:     strings.TrimSpace(d[siteDatum]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
				StationCode:  strings.TrimSpace(d[siteStationCode]),
				LocationCode: strings.TrimSpace(d[siteLocationCode]),
				Survey:       strings.TrimSpace(d[siteSurvey]),
			})
		}

		*s = SiteList(sites)
	}
	return nil
}

func LoadSites(path string) ([]Site, error) {
	var s []Site

	if err := LoadList(path, (*SiteList)(&s)); err != nil {
		return nil, err
	}

	sort.Sort(SiteList(s))

	return s, nil
}
