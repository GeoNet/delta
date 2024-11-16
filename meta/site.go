package meta

import (
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
	siteDepth
	siteDatum
	siteSurvey
	siteStart
	siteEnd
	siteLast
)

var siteHeaders Header = map[string]int{
	"Station":    siteStation,
	"Location":   siteLocation,
	"Latitude":   siteLatitude,
	"Longitude":  siteLongitude,
	"Elevation":  siteElevation,
	"Depth":      siteDepth,
	"Datum":      siteDatum,
	"Survey":     siteSurvey,
	"Start Date": siteStart,
	"End Date":   siteEnd,
}

var SiteTable Table = Table{
	name:    "Site",
	headers: siteHeaders,
	primary: []string{"Station", "Location", "Start Date"},
	native:  []string{"Latitude", "Longitude", "Elevation", "Depth"},
	foreign: map[string]map[string]string{
		"Station": {"Station": "Station"},
	},
	nullable: []string{"Depth", "Location"}, //TODO: the Location shouldn't be empty
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Site struct {
	Position
	Span

	Station  string `json:"station"`
	Location string `json:"location,omitempty"`
	Survey   string `json:"survey,omitempty"`
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
	var data [][]string

	data = append(data, siteHeaders.Columns())
	for _, row := range s {
		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Location),
			strings.TrimSpace(row.latitude),
			strings.TrimSpace(row.longitude),
			strings.TrimSpace(row.elevation),
			strings.TrimSpace(row.depth),
			strings.TrimSpace(row.Datum),
			strings.TrimSpace(row.Survey),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (s *SiteList) decode(data [][]string) error {
	var sites []Site

	if !(len(data) > 1) {
		return nil
	}

	fields := siteHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		lat, err := strconv.ParseFloat(d[siteLatitude], 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(d[siteLongitude], 64)
		if err != nil {
			return err
		}
		elev, err := ParseFloat64(d[siteElevation])
		if err != nil {
			return err
		}
		depth, err := ParseFloat64(d[siteDepth])
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[siteStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[siteEnd])
		if err != nil {
			return err
		}

		sites = append(sites, Site{
			Position: Position{
				// geographic details
				Latitude:  lat,
				Longitude: lon,
				Elevation: elev,
				Depth:     depth,
				Datum:     strings.TrimSpace(d[siteDatum]),

				// shadow variables
				latitude:  d[siteLatitude],
				longitude: d[siteLongitude],
				elevation: d[siteElevation],
				depth:     d[siteDepth],
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
