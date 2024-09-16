package meta

import (
	"testing"
)

func TestMethodList(t *testing.T) {
	t.Run("check methods", testListFunc("testdata/methods.csv", &MethodList{
		Method{
			Domain:      "envirosensor",
			Name:        "max",
			Description: "The maximum value of the observation over the output interval",
			Reference:   "https://help.campbellsci.com/crbasic/cr1000x/Content/Instructions/maximum.htm",
		},
		Method{
			Domain:      "manualcollect",
			Name:        "accumulation-chamber",
			Description: "An upside down bucket-shaped collection device is placed on the ground surface. A sensor in the chamber measures the concentration of gas coming from the ground beneath the chamber",
			Reference:   "https://www.geonet.org.nz/volcano/how",
		},
	}))
}
