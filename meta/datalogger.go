package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	dataloggerMake int = iota
	dataloggerModel
	dataloggerSerial
	dataloggerPlace
	dataloggerRole
	dataloggerStart
	dataloggerEnd
	dataloggerLast
)

type DeployedDatalogger struct {
	Install

	Place string
	Role  string
}

type DeployedDataloggerList []DeployedDatalogger

func (d DeployedDataloggerList) Len() int           { return len(d) }
func (d DeployedDataloggerList) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DeployedDataloggerList) Less(i, j int) bool { return d[i].Install.less(d[j].Install) }

func (d DeployedDataloggerList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Place",
		"Role",
		"Start Date",
		"End Date",
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

func (d *DeployedDataloggerList) decode(data [][]string) error {
	var dataloggers []DeployedDatalogger
	if len(data) > 1 {
		for _, v := range data[1:] {
			if len(v) != dataloggerLast {
				return fmt.Errorf("incorrect number of installed datalogger fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, v[dataloggerStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, v[dataloggerEnd]); err != nil {
				return err
			}

			dataloggers = append(dataloggers, DeployedDatalogger{
				Install: Install{
					Equipment: Equipment{
						Make:   strings.TrimSpace(v[dataloggerMake]),
						Model:  strings.TrimSpace(v[dataloggerModel]),
						Serial: strings.TrimSpace(v[dataloggerSerial]),
					},
					Span: Span{
						Start: start,
						End:   end,
					},
				},
				Place: strings.TrimSpace(v[dataloggerPlace]),
				Role:  strings.TrimSpace(v[dataloggerRole]),
			})
		}

		*d = DeployedDataloggerList(dataloggers)
	}
	return nil
}

func LoadDeployedDataloggers(path string) ([]DeployedDatalogger, error) {
	var d []DeployedDatalogger

	if err := LoadList(path, (*DeployedDataloggerList)(&d)); err != nil {
		return nil, err
	}

	sort.Sort(DeployedDataloggerList(d))

	return d, nil
}
