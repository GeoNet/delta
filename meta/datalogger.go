package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type DeployedDatalogger struct {
	Install

	Place string
	Role  string
}

type DeployedDataloggers []DeployedDatalogger

func (d DeployedDataloggers) Len() int           { return len(d) }
func (d DeployedDataloggers) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DeployedDataloggers) Less(i, j int) bool { return d[i].Install.Less(d[j].Install) }

func (d DeployedDataloggers) encode() [][]string {
	data := [][]string{{
		"Datalogger Make",
		"Datalogger Model",
		"Serial Number",
		"Deployment Place",
		"Deployment Role",
		"Installation Date",
		"Removal Date",
	}}
	for _, v := range d {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Place),
			strings.TrimSpace(v.Role),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (d *DeployedDataloggers) decode(data [][]string) error {
	var dataloggers []DeployedDatalogger
	if len(data) > 1 {
		for _, v := range data[1:] {
			if len(v) != 7 {
				return fmt.Errorf("incorrect number of installed datalogger fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, v[5]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, v[6]); err != nil {
				return err
			}

			dataloggers = append(dataloggers, DeployedDatalogger{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(v[0]),
						Model:  strings.TrimSpace(v[1]),
						Serial: strings.TrimSpace(v[2]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Place: strings.TrimSpace(v[3]),
				Role:  strings.TrimSpace(v[4]),
			})
		}

		*d = DeployedDataloggers(dataloggers)
	}
	return nil
}

func LoadDeployedDataloggers(path string) ([]DeployedDatalogger, error) {
	var d []DeployedDatalogger

	if err := LoadList(path, (*DeployedDataloggers)(&d)); err != nil {
		return nil, err
	}

	sort.Sort(DeployedDataloggers(d))

	return d, nil
}
