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
	featureSublocation
	featureProperty
	featureDescription
	featureAspect
	featureStart
	featureEnd
)

var featureHeaders Header = map[string]int{
	"Station":     featureStation,
	"Location":    featureLocation,
	"Sublocation": featureSublocation,
	"Property":    featureProperty,
	"Description": featureDescription,
	"Aspect":      featureAspect,
	"Start Date":  featureStart,
	"End Date":    featureEnd,
}

type Feature struct {
	Span

	Station     string
	Location    string
	Sublocation string
	Property    string
	Description string
	Aspect      string
}

// Id is a shorthand reference to the sample for debugging or testing.
func (f Feature) Id() string {
	return fmt.Sprintf("%s_%s_%s:%s", f.Station, f.Location, f.Sublocation, f.Start.Format(DateTimeFormat))
}

// Less allows samples to be sorted.
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
	case f.Sublocation < feature.Sublocation:
		return true
	case f.Sublocation > feature.Sublocation:
		return false
	case f.Property < feature.Property:
		return true
	case f.Property > feature.Property:
		return false
	case f.Start.Before(feature.Start):
		return true
	default:
		return false
	}
}

// Overlaps allows samples to be tested.
func (f Feature) Overlaps(feature Feature) bool {
	switch {
	case f.Station != feature.Station:
		return false
	case f.Location != feature.Location:
		return false
	case f.Sublocation != feature.Sublocation:
		return false
	case f.Property != feature.Property:
		return false
	case !f.Span.Overlaps(feature.Span):
		return false
	default:
		return true
	}
}

type FeatureList []Feature

func (f FeatureList) Len() int           { return len(f) }
func (f FeatureList) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f FeatureList) Less(i, j int) bool { return f[i].Less(f[j]) }

func (f FeatureList) encode() [][]string {
	var data [][]string

	data = append(data, featureHeaders.Columns())
	for _, v := range f {
		data = append(data, []string{
			strings.TrimSpace(v.Station),
			strings.TrimSpace(v.Location),
			strings.TrimSpace(v.Sublocation),
			strings.TrimSpace(v.Property),
			strings.TrimSpace(v.Description),
			strings.TrimSpace(v.Aspect),
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

	fields := featureHeaders.Fields(data[0])
	for _, v := range data[1:] {
		d := fields.Remap(v)

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
			Sublocation: strings.TrimSpace(d[featureSublocation]),
			Property:    strings.TrimSpace(d[featureProperty]),
			Description: strings.TrimSpace(d[featureDescription]),
			Aspect:      strings.TrimSpace(d[featureAspect]),
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
