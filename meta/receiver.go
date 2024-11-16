package meta

import (
	"sort"
	"strings"
	"time"
)

const (
	deployedReceiverMake int = iota
	deployedReceiverModel
	deployedReceiverSerial
	deployedReceiverMark
	deployedReceiverStart
	deployedReceiverEnd
	deployedReceiverLast
)

var deployedReceiverHeaders Header = map[string]int{
	"Make":       deployedReceiverMake,
	"Model":      deployedReceiverModel,
	"Serial":     deployedReceiverSerial,
	"Mark":       deployedReceiverMark,
	"Start Date": deployedReceiverStart,
	"End Date":   deployedReceiverEnd,
}

var DeployedReceiverTable Table = Table{
	name:    "Receiver",
	headers: deployedReceiverHeaders,
	primary: []string{"Make", "Model", "Serial", "Start Date"},
	native:  []string{},
	foreign: map[string]map[string]string{
		"Asset": {"Make": "Make", "Model": "Model", "Serial": "Serial"},
		"Mark":  {"Mark": "Mark"},
	},
	remap: map[string]string{
		"Start Date": "Start",
		"End Date":   "End",
	},
	start: "Start Date",
	end:   "End Date",
}

type DeployedReceiver struct {
	Install

	Mark string
}

type DeployedReceiverList []DeployedReceiver

func (dr DeployedReceiverList) Len() int           { return len(dr) }
func (dr DeployedReceiverList) Swap(i, j int)      { dr[i], dr[j] = dr[j], dr[i] }
func (dr DeployedReceiverList) Less(i, j int) bool { return dr[i].Install.Less(dr[j].Install) }

func (dr DeployedReceiverList) encode() [][]string {
	var data [][]string

	data = append(data, deployedReceiverHeaders.Columns())

	for _, row := range dr {
		data = append(data, []string{
			strings.TrimSpace(row.Make),
			strings.TrimSpace(row.Model),
			strings.TrimSpace(row.Serial),
			strings.TrimSpace(row.Mark),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (dr *DeployedReceiverList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var receivers []DeployedReceiver

	fields := deployedReceiverHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		start, err := time.Parse(DateTimeFormat, d[deployedReceiverStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[deployedReceiverEnd])
		if err != nil {
			return err
		}

		receivers = append(receivers, DeployedReceiver{
			Install: Install{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[deployedReceiverMake]),
					Model:  strings.TrimSpace(d[deployedReceiverModel]),
					Serial: strings.TrimSpace(d[deployedReceiverSerial]),
				},
				Span: Span{
					Start: start,
					End:   end,
				},
			},
			Mark: strings.TrimSpace(d[deployedReceiverMark]),
		})
	}

	*dr = DeployedReceiverList(receivers)

	return nil
}

func LoadDeployedReceivers(path string) ([]DeployedReceiver, error) {
	var dr []DeployedReceiver

	if err := LoadList(path, (*DeployedReceiverList)(&dr)); err != nil {
		return nil, err
	}

	sort.Sort(DeployedReceiverList(dr))

	return dr, nil
}
