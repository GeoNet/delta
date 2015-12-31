package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type DeployedReceiver struct {
	Install

	Place string
}

type DeployedReceiverList []DeployedReceiver

func (r DeployedReceiverList) Len() int           { return len(r) }
func (r DeployedReceiverList) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r DeployedReceiverList) Less(i, j int) bool { return r[i].Install.less(r[j].Install) }

func (r DeployedReceiverList) encode() [][]string {
	data := [][]string{{
		"Receiver Make",
		"Receiver Model",
		"Serial Number",
		"Deployment Place",
		"Installation Date",
		"Removal Date",
	}}
	for _, v := range r {
		data = append(data, []string{
			strings.TrimSpace(v.Make),
			strings.TrimSpace(v.Model),
			strings.TrimSpace(v.Serial),
			strings.TrimSpace(v.Place),
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
			if len(d) != 6 {
				return fmt.Errorf("incorrect number of installed receiver fields")
			}
			var err error

			var start, end time.Time
			if start, err = time.Parse(DateTimeFormat, d[4]); err != nil {
				return err
			}
			if end, err = time.Parse(DateTimeFormat, d[5]); err != nil {
				return err
			}

			receivers = append(receivers, DeployedReceiver{
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
