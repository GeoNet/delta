package meta

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	monumentMark int = iota
	monumentDomesNumber
	monumentMarkType
	monumentType
	monumentGroundRelationship
	monumentFoundationType
	monumentFoundationDepth
	monumentStart
	monumentEnd
	monumentBedrock
	monumentGeology
	monumentLast
)

var monumentHeaders Header = map[string]int{
	"Mark":                monumentMark,
	"Domes Number":        monumentDomesNumber,
	"Mark Type":           monumentMarkType,
	"Type":                monumentType,
	"Ground Relationship": monumentGroundRelationship,
	"Foundation Type":     monumentFoundationType,
	"Foundation Depth":    monumentFoundationDepth,
	"Start Date":          monumentStart,
	"End Date":            monumentEnd,
	"Bedrock":             monumentBedrock,
	"Geology":             monumentGeology,
}

type Monument struct {
	Span

	Mark               string
	DomesNumber        string
	MarkType           string
	Type               string
	GroundRelationship float64
	FoundationType     string
	FoundationDepth    float64
	Bedrock            string
	Geology            string

	groundRelationship string // shadow value to maintain formatting
	foundationDepth    string // shadow value to maintain formatting
}

type MonumentList []Monument

func (m MonumentList) Len() int           { return len(m) }
func (m MonumentList) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MonumentList) Less(i, j int) bool { return m[i].Mark < m[j].Mark }

func (m MonumentList) encode() [][]string {
	var data [][]string

	data = append(data, monumentHeaders.Columns())
	for _, row := range m {
		data = append(data, []string{
			strings.TrimSpace(row.Mark),
			strings.TrimSpace(row.DomesNumber),
			strings.TrimSpace(row.MarkType),
			strings.TrimSpace(row.Type),
			strings.TrimSpace(row.groundRelationship),
			strings.TrimSpace(row.FoundationType),
			strings.TrimSpace(row.foundationDepth),
			row.Start.Format(DateTimeFormat),
			row.End.Format(DateTimeFormat),
			strings.TrimSpace(row.Bedrock),
			strings.TrimSpace(row.Geology),
		})
	}

	return data
}

func (m *MonumentList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var monuments []Monument

	fields := monumentHeaders.Fields(data[0])
	for _, row := range data[1:] {
		d := fields.Remap(row)

		ground, err := strconv.ParseFloat(d[monumentGroundRelationship], 64)
		if err != nil {
			return err
		}
		depth, err := strconv.ParseFloat(d[monumentFoundationDepth], 64)
		if err != nil {
			return err
		}

		start, err := time.Parse(DateTimeFormat, d[monumentStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[monumentEnd])
		if err != nil {
			return err
		}

		monuments = append(monuments, Monument{
			Span: Span{
				Start: start,
				End:   end,
			},

			Mark:               strings.TrimSpace(d[monumentMark]),
			DomesNumber:        strings.TrimSpace(d[monumentDomesNumber]),
			MarkType:           strings.TrimSpace(d[monumentMarkType]),
			Type:               strings.TrimSpace(d[monumentType]),
			GroundRelationship: ground,
			FoundationType:     strings.TrimSpace(d[monumentFoundationType]),
			FoundationDepth:    depth,
			Bedrock:            strings.TrimSpace(d[monumentBedrock]),
			Geology:            strings.TrimSpace(d[monumentGeology]),

			groundRelationship: strings.TrimSpace(d[monumentGroundRelationship]),
			foundationDepth:    strings.TrimSpace(d[monumentFoundationDepth]),
		})
	}

	*m = MonumentList(monuments)

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
