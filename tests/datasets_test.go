package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var datasetChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for duplicate datasets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			datasets := make(map[meta.Dataset]interface{})
			for _, d := range set.Datasets() {
				key := meta.Dataset{
					Domain:  d.Domain,
					Network: d.Network,
					Key:     d.Key,
				}
				if _, ok := datasets[key]; ok {
					t.Errorf("citation %s/%s/%s is duplicated", d.Domain, d.Network, d.Key)
				}
				datasets[key] = true
			}
		}
	},
	"check for missing network": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, d := range set.Datasets() {
				if _, ok := set.Network(d.Network); !ok {
					t.Errorf("dataset network %s is unknown for %s", d.Network, d.Domain)
				}
			}
		}
	},
	"check for missing citations": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, d := range set.Datasets() {
				if d.Key == "" {
					continue
				}
				if _, ok := set.Citation(d.Key); !ok {
					t.Errorf("dataset citation %s is unknown for %s/%s", d.Key, d.Domain, d.Network)
				}
			}
		}
	},
}

func TestDatasets(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range datasetChecks {
		t.Run(k, v(set))
	}
}
