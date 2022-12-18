package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
	"github.com/GeoNet/delta/resp"
)

var respsChecks = map[string]func(*meta.Set) func(t *testing.T){
	"check for component response files": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, c := range set.Components() {
				if c.Response == "" {
					continue
				}
				switch data, err := resp.Lookup(c.Response); {
				case err != nil:
					t.Errorf("unable to lookup component response %s: %v", c.Response, err)
				case data == nil:
					t.Errorf("unable to find component response %s", c.Response)
				}
			}
		}
	},
	"check for channel response files": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, c := range set.Channels() {
				if c.Response == "" {
					continue
				}
				switch data, err := resp.Lookup(c.Response); {
				case err != nil:
					t.Errorf("unable to lookup channel response %s: %v", c.Response, err)
				case data == nil:
					t.Errorf("unable to find channel response %s", c.Response)
				}
			}
		}
	},
}

func TestResps(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range respsChecks {
		t.Run(k, v(set))
	}
}
