package meta

import (
	"fmt"
	"sort"
	"strings"
)

type Asset struct {
	Equipment
	Manufacturer string
	AssetNumber  string
	Notes        string
}

type Assets []Asset

func (a Assets) Len() int           { return len(a) }
func (a Assets) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Assets) Less(i, j int) bool { return a[i].Equipment.Less(a[j].Equipment) }

func (a Assets) encode() [][]string {
	data := [][]string{{"Manufacturer", "Make", "Model", "Serial Number", "Asset Number", "Notes"}}
	for _, v := range a {
		data = append(data, []string{
			v.Manufacturer,
			v.Make,
			v.Model,
			v.Serial,
			v.AssetNumber,
			v.Notes,
		})
	}
	return data
}

func (a *Assets) decode(data [][]string) error {
	var assets []Asset
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 6 {
				return fmt.Errorf("incorrect number of asset fields")
			}
			assets = append(assets, Asset{
				Equipment: Equipment{
					Make:   strings.TrimSpace(d[1]),
					Model:  strings.TrimSpace(d[2]),
					Serial: strings.TrimSpace(d[3]),
				},
				Manufacturer: strings.TrimSpace(d[0]),
				AssetNumber:  strings.TrimSpace(d[4]),
				Notes:        strings.TrimSpace(d[5]),
			})
		}

		*a = Assets(assets)
	}
	return nil
}

func LoadAssets(path string) ([]Asset, error) {
	var a []Asset

	if err := LoadList(path, (*Assets)(&a)); err != nil {
		return nil, err
	}

	sort.Sort(Assets(a))

	return a, nil
}
