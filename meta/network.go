package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	networkCode int = iota
	networkExternal
	networkDescription
	networkRestricted
	networkLast
)

type Network struct {
	NetworkCode  string
	ExternalCode string
	Description  string
	Restricted   bool
}

type NetworkList []Network

func (n NetworkList) Len() int           { return len(n) }
func (n NetworkList) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n NetworkList) Less(i, j int) bool { return n[i].NetworkCode < n[j].NetworkCode }

func (n NetworkList) encode() [][]string {
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

func (n *NetworkList) decode(data [][]string) error {
	var networks []Network
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != networkLast {
				return fmt.Errorf("incorrect number of installed network fields")
			}
			var err error

			var restricted bool
			if restricted, err = strconv.ParseBool(d[networkRestricted]); err != nil {
				return err
			}

			networks = append(networks, Network{
				NetworkCode:  strings.TrimSpace(d[networkCode]),
				ExternalCode: strings.TrimSpace(d[networkExternal]),
				Description:  strings.TrimSpace(d[networkDescription]),
				Restricted:   restricted,
			})
		}

		*n = NetworkList(networks)
	}
	return nil
}

func LoadNetworks(path string) ([]Network, error) {
	var n []Network

	if err := LoadList(path, (*NetworkList)(&n)); err != nil {
		return nil, err
	}

	sort.Sort(NetworkList(n))

	return n, nil
}
