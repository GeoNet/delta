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
func (d DeployedDataloggers) Less(i, j int) bool { return d[i].Install.less(d[j].Install) }

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

func (a *DeployedDataloggers) decode(data [][]string) error {
	var dataloggers []DeployedDatalogger
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 7 {
				return fmt.Errorf("incorrect number of installed datalogger fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[5]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[6]); err != nil {
				return err
			}

			dataloggers = append(dataloggers, DeployedDatalogger{
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
				Place: strings.TrimSpace(d[3]),
				Role:  strings.TrimSpace(d[4]),
			})
		}

		*a = DeployedDataloggers(dataloggers)
	}
	return nil
}

func LoadDeployedDataloggers(path string) ([]DeployedDatalogger, error) {
	var a []DeployedDatalogger

	if err := LoadList(path, (*DeployedDataloggers)(&a)); err != nil {
		return nil, err
	}

	sort.Sort(DeployedDataloggers(a))

	return a, nil
}
