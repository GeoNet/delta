package meta

import (
	"fmt"
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

type AssetList []Asset

func (a AssetList) Len() int           { return len(a) }
func (a AssetList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AssetList) Less(i, j int) bool { return a[i].Equipment.Less(a[j].Equipment) }

func (a AssetList) encode() [][]string {
	data := [][]string{{"Make", "Model", "Serial Number", "Asset Number", "Notes"}}
	for _, v := range a {
		data = append(data, []string{
			v.Make,
			v.Model,
			v.Serial,
			v.Number,
			v.Notes,
		})
	}
	return data
}

func (a *AssetList) decode(data [][]string) error {
	var assets []Asset
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != assetLast {
				return fmt.Errorf("incorrect number of asset fields")
			}
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
	}
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
