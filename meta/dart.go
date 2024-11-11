package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	dartStation = iota
	dartPid
	dartWmoIdentifier
	dartStartTime
	dartEndTime
	dartLast
)

var dartHeaders Header = map[string]int{
	"Station":        dartStation,
	"Pid":            dartPid,
	"WMO Identifier": dartWmoIdentifier,
	"Start Date":     dartStartTime,
	"End Date":       dartEndTime,
}

var DartTable Table = Table{
	name:    "Dart",
	headers: dartHeaders,
	primary: []string{"Station", "Start Date"},
	native:  []string{},
	foreign: map[string][]string{},
	remap: map[string]string{
		"WMO Identifier": "WmoIdentifier",
		"Start Date":     "Start",
		"End Date":       "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Dart struct {
	Span

	Station       string `json:"station"`
	Pid           string `json:"pid"`
	WmoIdentifier string `json:"wmo-identifier"`
}

type DartList []Dart

func (d DartList) Len() int           { return len(d) }
func (d DartList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DartList) Less(i, j int) bool { return d[i].Station < d[j].Station }

func (d DartList) encode() [][]string {
	var data [][]string

	data = append(data, dartHeaders.Columns())
	for _, row := range d {
		data = append(data, []string{
			strings.TrimSpace(row.Station),
			strings.TrimSpace(row.Pid),
			strings.TrimSpace(row.WmoIdentifier),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (d *DartList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var darts []Dart

	fields := dartHeaders.Fields(data[0])
	for _, row := range data[1:] {
		r := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, r[dartStartTime])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, r[dartEndTime])
		if err != nil {
			return err
		}

		darts = append(darts, Dart{
			Span: Span{
				Start: start,
				End:   end,
			},

			Station:       strings.TrimSpace(r[dartStation]),
			Pid:           strings.TrimSpace(r[dartPid]),
			WmoIdentifier: strings.TrimSpace(r[dartWmoIdentifier]),
		})
	}

	*d = DartList(darts)

	return nil
}

func LoadDarts(path string) ([]Dart, error) {
	var d []Dart

	if err := LoadList(path, (*DartList)(&d)); err != nil {
		return nil, err
	}

	sort.Sort(DartList(d))

	return d, nil
}
