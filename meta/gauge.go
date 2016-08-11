package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	gaugeGaugeMake = iota
	gaugeGaugeModel
	gaugeSerialNumber
	gaugeStationCode
	gaugeLocationCode
	gaugeInstallationDip
	gaugeVerticalOffset
	gaugeOffsetNorth
	gaugeOffsetEast
	gaugeScaleFactor
	gaugeScaleBias
	gaugeCableLength
	gaugeInstallationDate
	gaugeRemovalDate
	gaugeLast
)

type InstalledGauge struct {
	Install
	Offset
	Orientation
	Scale

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
		"Scale Factor",
		"Scale Bias",
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
			strconv.FormatFloat(v.Factor, 'g', -1, 64),
			strconv.FormatFloat(v.Bias, 'g', -1, 64),
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
			if len(d) != gaugeLast {
				return fmt.Errorf("incorrect number of installed gauge fields")
			}
			var err error

			var dip float64
			if dip, err = strconv.ParseFloat(d[gaugeInstallationDip], 64); err != nil {
				return err
			}

			var height, north, east float64
			if height, err = strconv.ParseFloat(d[gaugeVerticalOffset], 64); err != nil {
				return err
			}
			if north, err = strconv.ParseFloat(d[gaugeOffsetNorth], 64); err != nil {
				return err
			}
			if east, err = strconv.ParseFloat(d[gaugeOffsetEast], 64); err != nil {
				return err
			}

			var factor, bias float64
			if factor, err = strconv.ParseFloat(d[gaugeScaleFactor], 64); err != nil {
				return err
			}
			if bias, err = strconv.ParseFloat(d[gaugeScaleBias], 64); err != nil {
				return err
			}

			var length float64
			if length, err = strconv.ParseFloat(d[gaugeCableLength], 64); err != nil {
				return err
			}

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[gaugeInstallationDate]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[gaugeRemovalDate]); err != nil {
				return err
			}

			gauges = append(gauges, InstalledGauge{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(d[gaugeGaugeMake]),
						Model:  strings.TrimSpace(d[gaugeGaugeModel]),
						Serial: strings.TrimSpace(d[gaugeSerialNumber]),
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
				Scale: Scale{
					Factor: factor,
					Bias:   bias,
				},
				StationCode:  strings.TrimSpace(d[gaugeStationCode]),
				LocationCode: strings.TrimSpace(d[gaugeLocationCode]),
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
