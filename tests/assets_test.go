package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var assetChecks = map[string]func(*meta.Set) func(t *testing.T){

	"check for duplicate asset numbers": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			reference := make(map[string]string)
			for _, a := range set.Assets() {
				if a.Number == "" {
					continue
				}
				if x, ok := reference[a.Number]; ok {
					t.Errorf("duplicate asset number %s: %s [%s]", a.String(), a.Number, x)
				}
				reference[a.Number] = a.String()
			}
		}
	},

	"check for duplicate equipment": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			assets := set.Assets()
			for i := 0; i < len(assets); i++ {
				for j := i + 1; j < len(assets); j++ {
					switch {
					case assets[i].Model != assets[j].Model:
					case assets[i].Serial != assets[j].Serial:
					default:
						t.Errorf("duplicate equipment %s: %s", assets[i].String(), assets[j].String())
					}
				}
			}
		}
	},
}

func TestAssets(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range assetChecks {
		t.Run(k, v(set))
	}
}
