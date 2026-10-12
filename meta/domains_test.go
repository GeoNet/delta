package meta

import (
	"testing"
)

func TestDomain(t *testing.T) {
	t.Run("check domains", testListFunc("testdata/domains.csv", &DomainList{
		Domain{
			Name:        "coastal",
			Description: "Coastal Tsunami Gauge Network",
		},
		Domain{
			Name: "gnss",
		},
	}))
}
