package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	deployedDataloggerMake int = iota
	deployedDataloggerModel
	deployedDataloggerSerial
	deployedDataloggerPlace
	deployedDataloggerRole
	deployedDataloggerStart
	deployedDataloggerEnd
	deployedDataloggerLast
)

var deployedDataloggerHeaders Header = map[string]int{
	"Make":       deployedDataloggerMake,
	"Model":      deployedDataloggerModel,
	"Serial":     deployedDataloggerSerial,
	"Place":      deployedDataloggerPlace,
	"Role":       deployedDataloggerRole,
	"Start Date": deployedDataloggerStart,
	"End Date":   deployedDataloggerEnd,
}

var DeployedDataloggerTable Table = Table{
	name:    "Datalogger",
	headers: deployedDataloggerHeaders,
	primary: []string{"Make", "Model", "Serial", "Place", "Role", "Start Date"},
	native:  []string{},
	foreign: map[string][]string{
		"Asset":      {"Make", "Model", "Serial"},
		"Connection": {"Place", "Role"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type DeployedDatalogger struct {
	Install

	Place string `json:"place"`
	Role  string `json:"role,omitempty"`
}

type DeployedDataloggerList []DeployedDatalogger

func (dd DeployedDataloggerList) Len() int           { return len(dd) }
func (dd DeployedDataloggerList) Swap(i, j int)      { dd[i], dd[j] = dd[j], dd[i] }
func (dd DeployedDataloggerList) Less(i, j int) bool { return dd[i].Install.Less(dd[j].Install) }

func (dd DeployedDataloggerList) encode() [][]string {

	var data [][]string

	data = append(data, deployedDataloggerHeaders.Columns())

	for _, row := range dd {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Place),
			strings.TrimSpace(row.Role),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (dd *DeployedDataloggerList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var deployedDataloggers []DeployedDatalogger

	fields := deployedDataloggerHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[deployedDataloggerStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[deployedDataloggerEnd])
		if err != nil {
			return err
		}

		deployedDataloggers = append(deployedDataloggers, DeployedDatalogger{
			Install: Install{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[deployedDataloggerMake]),
					Model:  strings.TrimSpace(d[deployedDataloggerModel]),
					Serial: strings.TrimSpace(d[deployedDataloggerSerial]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
			},
			Place: strings.TrimSpace(d[deployedDataloggerPlace]),
			Role:  strings.TrimSpace(d[deployedDataloggerRole]),
		})
	}

	*dd = DeployedDataloggerList(deployedDataloggers)

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
