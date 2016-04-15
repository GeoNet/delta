package meta

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	//	"time"
)

type Monument struct {
	MarkCode           string
	DomesNumber        string
	MarkType           string
	MonumentType       string
	GroundRelationship float64
	Bedrock            string
	Geology            string
}

const (
	monumentMarkCode int = iota
	monumentDomesNumber
	monumentMarkType
	monumentMonumentType
	monumentGroundRelationship
	monumentBedrock
	monumentGeology
)

type MonumentList []Monument

func (m MonumentList) Len() int           { return len(m) }
func (m MonumentList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MonumentList) Less(i, j int) bool { return m[i].MarkCode < m[j].MarkCode }

func (m MonumentList) encode() [][]string {
	data := [][]string{{
		"Mark Code",
		"Domes Number",
		"Mark Type",
		"Monument Type",
		"Ground Relationship",
		"Bedrock",
		"Geology",
	}}
	for _, v := range m {
		data = append(data, []string{
			strings.TrimSpace(v.MarkCode),
			strings.TrimSpace(v.DomesNumber),
			strings.TrimSpace(v.MarkType),
			strings.TrimSpace(v.MonumentType),
			strconv.FormatFloat(v.GroundRelationship, 'g', -1, 64),
			strings.TrimSpace(v.Bedrock),
			strings.TrimSpace(v.Geology),
		})
	}
	return data
}

func (m *MonumentList) decode(data [][]string) error {
	var monuments []Monument
	if len(data) > 1 {
		for _, d := range data[1:] {
			var err error

			if len(d) != 7 {
				return fmt.Errorf("incorrect number of monument fields")
			}
			var ground float64
			if ground, err = strconv.ParseFloat(d[monumentGroundRelationship], 64); err != nil {
				return err
			}
			monuments = append(monuments, Monument{
				MarkCode:           strings.TrimSpace(d[monumentMarkCode]),
				DomesNumber:        strings.TrimSpace(d[monumentDomesNumber]),
				MarkType:           strings.TrimSpace(d[monumentMarkType]),
				MonumentType:       strings.TrimSpace(d[monumentMonumentType]),
				GroundRelationship: ground,
				Bedrock:            strings.TrimSpace(d[monumentBedrock]),
				Geology:            strings.TrimSpace(d[monumentGeology]),
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
