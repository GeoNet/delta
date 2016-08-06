package main

import (
	"testing"
)

func TestResponses_Response(t *testing.T) {

	for _, v := range Responses {
		for n, r := range v.Sensors {
			for _, s := range r.Sensors {
				if _, ok := SensorModels[s]; !ok {
					t.Errorf("%d: unknown sensor model lookup: %s", n, s)
				}
			}
			if r.Channels == "" {
				t.Errorf("%d: empty sensor channels list: %d", n)
			}
		}
		for n, r := range v.Dataloggers {
			for _, d := range r.Dataloggers {
				if _, ok := DataloggerModels[d]; !ok {
					t.Errorf("%d: unknown datalogger model lookup: %s", n, d)
				}
			}
			if len(r.Type) != 2 {
				t.Errorf("%d: invalid datalogger type string: %s", n, r.Type)
			}
			if len(r.Label) != 2 {
				t.Errorf("%d: invalid datalogger label string: %s", n, r.Label)
			}
			if r.StorageFormat == "" {
				t.Errorf("%d: empty datalogger storage format", n)
			}
		}
	}
}
