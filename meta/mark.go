package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	markCode = iota
	markNetwork
	markIgs
	markName
	markLatitude
	markLongitude
	markElevation
	markDatum
	markStartTime
	markEndTime
	markLast
)

var markHeaders Header = map[string]int{
	"Mark":       markCode,
	"Network":    markNetwork,
	"Igs":        markIgs,
	"Name":       markName,
	"Latitude":   markLatitude,
	"Longitude":  markLongitude,
	"Elevation":  markElevation,
	"Datum":      markDatum,
	"Start Date": markStartTime,
	"End Date":   markEndTime,
}

var MarkTable Table = Table{
	name:    "Mark",
	headers: markHeaders,
	primary: []string{"Mark", "Start Date"},
	native:  []string{"Latitude", "Longitude", "Elevation"},
	foreign: map[string][]string{
		"Network": {"Network"},
	},
	remap: map[string]string{
		"Mark":       "Code",
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type Mark struct {
	Reference
	Position
	Span

	Igs bool `json:"igs,omitempty"`
}

type MarkList []Mark

func (m MarkList) Len() int           { return len(m) }
func (m MarkList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MarkList) Less(i, j int) bool { return m[i].Code < m[j].Code }

func (m MarkList) encode() [][]string {
	var data [][]string

	data = append(data, markHeaders.Columns())
	for _, row := range m {
		data = append(data, []string{
			strings.TrimSpace(row.Code),
			strings.TrimSpace(row.Network),
			func() string {
				if row.Igs {
					return "yes"
				}
				return "no"
			}(),
			strings.TrimSpace(row.Name),
			strings.TrimSpace(row.latitude),
			strings.TrimSpace(row.longitude),
			strings.TrimSpace(row.elevation),
			strings.TrimSpace(row.Datum),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (m *MarkList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var marks []Mark

	fields := markHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		igs, err := strconv.ParseBool(d[markIgs])
		if err != nil {
			switch d[markIgs] {
			case "y", "Y", "yes", "YES":
				igs = true
			case "n", "N", "no", "NO":
				igs = false
			default:
				return err
			}
		}

		lat, err := strconv.ParseFloat(d[markLatitude], 64)
		if err != nil {
			return err
		}
		lon, err := strconv.ParseFloat(d[markLongitude], 64)
		if err != nil {
			return err
		}
		elev, err := strconv.ParseFloat(d[markElevation], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[markStartTime])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[markEndTime])
		if err != nil {
			return err
		}

		marks = append(marks, Mark{
			Reference: Reference{
				Code:    strings.TrimSpace(d[markCode]),
				Network: strings.TrimSpace(d[markNetwork]),
				Name:    strings.TrimSpace(d[markName]),
			},
			Span: Span{
				Start: start,
				End:   end,
			},
			Position: Position{
				Latitude:  lat,
				Longitude: lon,
				Elevation: elev,
				Datum:     strings.TrimSpace(d[markDatum]),

				latitude:  strings.TrimSpace(d[markLatitude]),
				longitude: strings.TrimSpace(d[markLongitude]),
				elevation: strings.TrimSpace(d[markElevation]),
			},
			Igs: igs,
		})
	}

	*m = MarkList(marks)

	return nil
}

func LoadMarks(path string) ([]Mark, error) {
	var m []Mark

	if err := LoadList(path, (*MarkList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(MarkList(m))

	return m, nil
}
