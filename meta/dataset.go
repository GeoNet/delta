package meta

import (
	"sort"
	"strconv"
	"strings"
)

const (
	datasetDomain int = iota
	datasetNetwork
	datasetKey
	datasetTilde
	datasetLast
)

var datasetHeaders Header = map[string]int{
	"Domain":  datasetDomain,
	"Network": datasetNetwork,
	"Key":     datasetKey,
	"Tilde":   datasetTilde,
}

type Dataset struct {
	Domain  string
	Network string
	Key     string
	Tilde   bool

	tilde string
}

type DatasetList []Dataset

func (d DatasetList) Len() int      { return len(d) }
func (d DatasetList) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d DatasetList) Less(i, j int) bool {
	switch {
	case d[i].Domain < d[j].Domain:
		return true
	case d[i].Domain > d[j].Domain:
		return false
	case d[i].Network < d[j].Network:
		return true
	case d[i].Network > d[j].Network:
		return false
	case d[i].Key < d[j].Key:
		return true
	default:
		return false
	}
}

func (d DatasetList) encode() [][]string {
	var data [][]string

	data = append(data, datasetHeaders.Columns())

	for _, row := range d {
		data = append(data, []string{
			row.Domain,
			row.Network,
			row.Key,
			row.tilde,
		})
	}

	return data
}

func (d *DatasetList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var datasets []Dataset

	fields := datasetHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		tilde, err := strconv.ParseBool(d[datasetTilde])
		if err != nil {
			return err
		}

		datasets = append(datasets, Dataset{
			Domain:  strings.TrimSpace(d[datasetDomain]),
			Network: strings.TrimSpace(d[datasetNetwork]),
			Key:     strings.TrimSpace(d[datasetKey]),
			Tilde:   tilde,

			tilde: strings.TrimSpace(d[datasetTilde]),
		})
	}

	*d = DatasetList(datasets)

	return nil
}

func LoadDatasets(path string) ([]Dataset, error) {
	var c []Dataset

	if err := LoadList(path, (*DatasetList)(&c)); err != nil {
		return nil, err
	}

	sort.Sort(DatasetList(c))

	return c, nil
}
