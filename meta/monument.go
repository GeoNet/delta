package meta

import (
	"fmt"
	"sort"
	//	"strconv"
	"strings"
	//	"time"
)

type Monument struct {
	MarkCode string

	MarkType     string
	MonumentType string
	DomesNumber  string
}

type MonumentList []Monument

func (m MonumentList) Len() int           { return len(m) }
func (m MonumentList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MonumentList) Less(i, j int) bool { return m[i].MarkCode < m[j].MarkCode }

func (m MonumentList) encode() [][]string {
	data := [][]string{{
		"Mark Code",
		"Mark Type",
		"Monument Type",
		"Domes Number",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.MarkCode),
			strings.TrimSpace(v.MarkType),
			strings.TrimSpace(v.MonumentType),
			strings.TrimSpace(v.DomesNumber),
		})
	}
	return data
}

func (m *MonumentList) decode(data [][]string) error {
	var monuments []Monument
	if len(data) > 1 {
		for _, d := range data[1:] {
			if len(d) != 4 {
				return fmt.Errorf("incorrect number of installed monument fields")
			}
			monuments = append(monuments, Monument{
				MarkCode:     strings.TrimSpace(d[0]),
				MarkType:     strings.TrimSpace(d[1]),
				MonumentType: strings.TrimSpace(d[2]),
				DomesNumber:  strings.TrimSpace(d[3]),
			})
		}

		*m = MonumentList(monuments)
	}
	return nil
}

func LoadMonuments(path string) ([]Monument, error) {
	var m []Monument

	if err := LoadList(path, (*MonumentList)(&m)); err != nil {
		return nil, err
	}

	sort.Sort(MonumentList(m))

	return m, nil
}
