package meta

import (
	"testing"
)

func TestDatasetList(t *testing.T) {
	t.Run("check datasets", testListFunc("testdata/datasets.csv", &DatasetList{
		Dataset{
			Domain:  "acoustic",
			Network: "HA",
			Key:     "Gns2022a",
			Tilde:   false,

			tilde: "false",
		},
		Dataset{
			Domain:  "camera",
			Network: "VC",
			Key:     "",
			Tilde:   false,

			tilde: "false",
		},
		Dataset{
			Domain:  "coastal",
			Network: "LG",
			Key:     "Gns2007a",
			Tilde:   true,

			tilde: "true",
		},
	}))
}
