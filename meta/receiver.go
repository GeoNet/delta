package meta

import (
	"fmt"
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

type DeployedReceiver struct {
	Install

	Mark string
}

type DeployedReceiverList []DeployedReceiver

func (r DeployedReceiverList) Len() int           { return len(r) }
func (r DeployedReceiverList) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r DeployedReceiverList) Less(i, j int) bool { return r[i].Install.less(r[j].Install) }

func (r DeployedReceiverList) encode() [][]string {
	data := [][]string{{
		"Make",
		"Model",
		"Serial",
		"Mark",
		"Start Date",
		"End Date",
	}}
	for _, v := range r {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Mark),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}
	return data
}

func (r *DeployedReceiverList) decode(data [][]string) error {
	var receivers []DeployedReceiver
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != deployedReceiverLast {
				return fmt.Errorf("incorrect number of installed receiver fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[deployedReceiverStart]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[deployedReceiverEnd]); err != nil {
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

		*r = DeployedReceiverList(receivers)
	}
	return nil
}

func LoadDeployedReceivers(path string) ([]DeployedReceiver, error) {
	var r []DeployedReceiver

	if err := LoadList(path, (*DeployedReceiverList)(&r)); err != nil {
		return nil, err
	}

	sort.Sort(DeployedReceiverList(r))

	return r, nil
}
