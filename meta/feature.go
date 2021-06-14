package meta

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	featureStation = iota
	featureLocation
	featureDescription
	featureStart
	featureEnd
	featureLast
)

type Feature struct {
	Span

	Station     string
	Location    string
	Description string
}

func (f Feature) Less(feature Feature) bool {
	switch {
	case f.Station < feature.Station:
		return true
	case f.Station > feature.Station:
		return false
	case f.Location < feature.Location:
		return true
	case f.Location > feature.Location:
		return false
	case f.Start.Before(feature.Start):
		return true
	default:
		return false
	}
}

type FeatureList []Feature

func (f FeatureList) Len() int           { return len(f) }
func (f FeatureList) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FeatureList) Less(i, j int) bool { return f[i].Less(f[j]) }

func (f FeatureList) encode() [][]string {
	data := [][]string{{
		"Station",
		"Location",
		"Description",
		"Start Date",
		"End Date",
	}}

	for _, v := range f {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.Description),
			v.Start.Format(DateTimeFormat),
			v.End.Format(DateTimeFormat),
		})
	}

	return data
}

func (f *FeatureList) decode(data [][]string) error {
	if !(len(data) > 1) {
		return nil
	}

	var features []Feature

	for _, d := range data[1:] {
		if len(d) != featureLast {
			return fmt.Errorf("incorrect number of feature fields")
		}

		start, err := time.Parse(DateTimeFormat, d[featureStart])
		if err != nil {
			return err
		}
		end, err := time.Parse(DateTimeFormat, d[featureEnd])
		if err != nil {
			return err
		}

		features = append(features, Feature{
			Span: Span{
				Start: start,
				End:   end,
			},
			Station:     strings.TrimSpace(d[featureStation]),
			Location:    strings.TrimSpace(d[featureLocation]),
			Description: strings.TrimSpace(d[featureDescription]),
		})
	}

	*f = FeatureList(features)

	return nil
}

func LoadFeatures(path string) ([]Feature, error) {
	var f []Feature

	if err := LoadList(path, (*FeatureList)(&f)); err != nil {
		return nil, err
	}

	sort.Sort(FeatureList(f))

	return f, nil
}
