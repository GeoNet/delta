package meta

import (
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

var networkHeaders Header = map[string]int{
	"Network":     networkCode,
	"External":    networkExternal,
	"Description": networkDescription,
	"Restricted":  networkRestricted,
}

type Network struct {
	Code        string
	External    string
	Description string
	Restricted  bool
}

type NetworkList []Network

func (n NetworkList) Len() int           { return len(n) }
func (n NetworkList) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n NetworkList) Less(i, j int) bool { return n[i].Code < n[j].Code }

func (n NetworkList) encode() [][]string {
	var data [][]string

	data = append(data, networkHeaders.Columns())
	for _, row := range n {
		data = append(data, []string{
			strings.TrimSpace(row.Code),
			strings.TrimSpace(row.External),
			strings.TrimSpace(row.Description),
			strconv.FormatBool(row.Restricted),
		})
	}

	return data
}

func (n *NetworkList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var networks []Network

	fields := networkHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		restricted, err := strconv.ParseBool(d[networkRestricted])
		if err != nil {
			return err
		}

		networks = append(networks, Network{
			Code:        strings.TrimSpace(d[networkCode]),
			External:    strings.TrimSpace(d[networkExternal]),
			Description: strings.TrimSpace(d[networkDescription]),
			Restricted:  restricted,
		})
	}

	*n = NetworkList(networks)

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
