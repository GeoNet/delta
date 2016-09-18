package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	siteStation = iota
	siteLocation
	siteLatitude
	siteLongitude
	siteElevation
	siteDatum
	siteSurvey
	siteStart
	siteEnd
	siteLast
)

type Site struct {
	Point
	Span

	Station  string
	Location string
	Survey   string
}

func (s Site) Less(site Site) bool {
	switch {
	case s.Station < site.Station:
		return true
	case s.Station > site.Station:
		return false
	case s.Location < site.Location:
		return true
	case s.Location > site.Location:
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
		"Station",
		"Location",
		"Latitude",
		"Longitude",
		"Elevation",
		"Datum",
		"Survey",
		"Start Date",
		"End Date",
	}}
	for _, v := range s {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
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
			if elev, err = strconv.ParseFloat(d[siteElevation], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[siteStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[siteEnd]); err != nil {
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
				Station:  strings.TrimSpace(d[siteStation]),
				Location: strings.TrimSpace(d[siteLocation]),
				Survey:   strings.TrimSpace(d[siteSurvey]),
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
