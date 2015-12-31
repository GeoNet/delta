package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Network struct {
	NetworkCode  string
	ExternalCode string
	Description  string
	Restricted   bool
}

type Networks []Network

func (n Networks) Len() int           { return len(n) }
func (n Networks) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Networks) Less(i, j int) bool { return n[i].NetworkCode < n[j].NetworkCode }

func (n Networks) encode() [][]string {
	data := [][]string{{
		"Network Code",
		"External Code",
		"Description",
		"Restricted",
	}}
	for _, v := range n {
		data = append(data, []string{
			strings.TrimSpace(v.NetworkCode),
			strings.TrimSpace(v.ExternalCode),
			strings.TrimSpace(v.Description),
			strconv.FormatBool(v.Restricted),
		})
	}
	return data
}

func (n *Networks) decode(data [][]string) error {
	var networks []Network
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 4 {
				return fmt.Errorf("incorrect number of installed network fields")
			}
			var err error

			var restricted bool
			if restricted, err = strconv.ParseBool(d[3]); err != nil {
				return err
			}

			networks = append(networks, Network{
				NetworkCode:  strings.TrimSpace(d[0]),
				ExternalCode: strings.TrimSpace(d[1]),
				Description:  strings.TrimSpace(d[2]),
				Restricted:   restricted,
			})
		}

		*n = Networks(networks)
	}
	return nil
}

func LoadNetworks(path string) ([]Network, error) {
	var n []Network

	if err := LoadList(path, (*Networks)(&n)); err != nil {
		return nil, err
	}

	sort.Sort(Networks(n))

	return n, nil
}
