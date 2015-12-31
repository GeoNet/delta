package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type InstalledGauge struct {
	Install
	Offset
	Orientation

	StationCode  string
	LocationCode string
	CableLength  float64
}

type InstalledGaugeList []InstalledGauge

func (g InstalledGaugeList) Len() int           { return len(g) }
func (g InstalledGaugeList) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g InstalledGaugeList) Less(i, j int) bool { return g[i].Install.less(g[j].Install) }

func (g InstalledGaugeList) encode() [][]string {
	data := [][]string{{
		"Gauge Make",
		"Gauge Model",
		"Serial Number",
		"Station Code",
		"Location Code",
		"Installation Dip",
		"Vertical Offset",
		"Offset North",
		"Offset East",
		"Cable Length",
		"Installation Date",
		"Removal Date",
	}}
	for _, v := range g {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.StationCode),
			strings.TrimSpace(v.LocationCode),
			strconv.FormatFloat(v.Dip, 'g', -1, 64),
			strconv.FormatFloat(v.Height, 'g', -1, 64),
			strconv.FormatFloat(v.North, 'g', -1, 64),
			strconv.FormatFloat(v.East, 'g', -1, 64),
			strconv.FormatFloat(v.CableLength, 'g', -1, 64),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (g *InstalledGaugeList) decode(data [][]string) error {
	var gauges []InstalledGauge
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 12 {
				return fmt.Errorf("incorrect number of installed gauge fields")
			}
			var err error

			var dip float64
			if dip, err = strconv.ParseFloat(d[5], 64); err != nil {
				return err
			}

			var height, north, east float64
			if height, err = strconv.ParseFloat(d[6], 64); err != nil {
				return err
			}
			if north, err = strconv.ParseFloat(d[7], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[8], 64); err != nil {
				return err
			}

			var length float64
			if length, err = strconv.ParseFloat(d[9], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[10]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[11]); err != nil {
				return err
			}

			gauges = append(gauges, InstalledGauge{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[0]),
						Model:  strings.TrimSpace(d[1]),
						Serial: strings.TrimSpace(d[2]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Orientation: Orientation{
					Dip: dip,
				},
				Offset: Offset{
					Height: height,
					North:  north,
					East:   east,
				},
				StationCode:  strings.TrimSpace(d[3]),
				LocationCode: strings.TrimSpace(d[4]),
				CableLength:  length,
			})
		}

		*g = InstalledGaugeList(gauges)
	}
	return nil
}

func LoadInstalledGauges(path string) ([]InstalledGauge, error) {
	var g []InstalledGauge

	if err := LoadList(path, (*InstalledGaugeList)(&g)); err != nil {
		return nil, err
	}

	sort.Sort(InstalledGaugeList(g))

	return g, nil
}
