package delta_test

import (
	"testing"

	"github.com/GeoNet/delta"
	"github.com/GeoNet/delta/meta"
)

var antennaChecks = map[string]func(*meta.Set) func(t *testing.T){

	// check for session overlaps, there can't be two sessions running at the same mark for the same sampling interval.
	"check antenna installation overlap": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {

			installs := make(map[string]meta.InstalledAntennaList)
			for _, s := range set.InstalledAntennas() {
				installs[s.Model] = append(installs[s.Model], s)
			}

			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
						if v[i].Serial != v[j].Serial {
							continue
						}
						if v[i].End.Before(v[j].Start) {
							continue
						}
						if v[i].Start.After(v[j].End) {
							continue
						}
						if v[i].End.Equal(v[j].Start) {
							continue
						}
						if v[i].Start.Equal(v[j].End) {
							continue
						}
						t.Errorf("installed antennas %s [%s] at %s has overlap at %s between %s and %s",
							v[i].Model, v[i].Serial, v[i].Mark, v[j].Mark,
							v[i].Start.Format(meta.DateTimeFormat),
							v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for invalid installation dates": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, i := range set.InstalledAntennas() {
				if i.End.After(i.Start) {
					continue
				}
				t.Errorf("installed antenna is removed before it has been installed: %s", i.String())
			}
		}
	},

	"check for antenna installation mark overlaps": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			installs := make(map[string]meta.InstalledAntennaList)
			for _, s := range set.InstalledAntennas() {
				installs[s.Mark] = append(installs[s.Mark], s)
			}

			for _, v := range installs {
				for i := 0; i < len(v); i++ {
					for j := i + 1; j < len(v); j++ {
						if v[i].End.Before(v[j].Start) {
							continue
						}
						if v[i].Start.After(v[j].End) {
							continue
						}
						if v[i].End.Equal(v[j].Start) {
							continue
						}
						if v[i].Start.Equal(v[j].End) {
							continue
						}
						t.Errorf("installed antennas %s [%s] and %s [%s] at %s has an overlap between %s and %s",
							v[i].Model, v[i].Serial, v[j].Model, v[j].Serial, v[i].Mark,
							v[i].Start.Format(meta.DateTimeFormat), v[i].End.Format(meta.DateTimeFormat))
					}
				}
			}
		}
	},

	"check for missing antenna marks": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			keys := make(map[string]interface{})
			for _, m := range set.Marks() {
				keys[m.Code] = true
			}

			for _, c := range set.InstalledAntennas() {
				if _, ok := keys[c.Mark]; !ok {
					t.Errorf("unable to find antenna mark %-5s", c.Mark)
				}
			}
		}
	},

	"check for missing antenna assets": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, r := range set.InstalledAntennas() {
				var found bool
				for _, a := range set.Assets() {
					if a.Model != r.Model {
						continue
					}
					if a.Serial != r.Serial {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("unable to find antenna asset: %s [%s]", r.Model, r.Serial)
			}
		}
	},

	"check for missing antenna sessions": func(set *meta.Set) func(t *testing.T) {
		return func(t *testing.T) {
			for _, r := range set.InstalledAntennas() {
				var found bool
				for _, s := range set.Sessions() {
					if s.End.Before(r.Start) {
						continue
					}
					if s.Start.After(r.End) {
						continue
					}
					found = true
				}
				if found {
					continue
				}
				t.Errorf("unable to find antenna session: %s [%s]", r.Model, r.Serial)
			}
		}
	},
}

func TestAntennas(t *testing.T) {

	set, err := delta.New()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range antennaChecks {
		t.Run(k, v(set))
	}
}
