package meta

import (
	"sort"
	"strings"
)

type Asset struct {
	Equipment
	Number string
	Notes  string
}

const (
	assetMake int = iota
	assetModel
	assetSerial
	assetNumber
	assetNotes
	assetLast
)

var assetHeaders Header = map[string]int{
	"Make":   assetMake,
	"Model":  assetModel,
	"Serial": assetSerial,
	"Number": assetNumber,
	"Notes":  assetNotes,
}

type AssetList []Asset

func (a AssetList) Len() int           { return len(a) }
func (a AssetList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AssetList) Less(i, j int) bool { return a[i].Equipment.Less(a[j].Equipment) }

func (a AssetList) encode() [][]string {
	var data [][]string

	data = append(data, assetHeaders.Columns())

	for _, row := range a {
		data = append(data, []string{
			row.Make,
			row.Model,
			row.Serial,
			row.Number,
			row.Notes,
		})
	}
	return data
}

func (a *AssetList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var assets []Asset

	fields := assetHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		assets = append(assets, Asset{
			Equipment: Equipment{
				Make:   strings.TrimSpace(d[assetMake]),
				Model:  strings.TrimSpace(d[assetModel]),
				Serial: strings.TrimSpace(d[assetSerial]),
			},
			Number: strings.TrimSpace(d[assetNumber]),
			Notes:  strings.TrimSpace(d[assetNotes]),
		})
	}

	*a = AssetList(assets)

	return nil
}

func LoadAssets(path string) ([]Asset, error) {
	var a []Asset

	if err := LoadList(path, (*AssetList)(&a)); err != nil {
		return nil, err
	}

	sort.Sort(AssetList(a))

	return a, nil
}
