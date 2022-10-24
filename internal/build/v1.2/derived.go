package build

import (
	"encoding/xml"
	"fmt"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

func (r *Response) Derived(data []byte) (*stationxml.ResponseType, error) {

	var derived stationxml.ResponseType
	if err := xml.Unmarshal(data, &derived); err != nil {
		return nil, err
	}

	// must have at least an instrument sensitivity or polynomial
	if derived.InstrumentSensitivity == nil && derived.InstrumentPolynomial == nil {
		return nil, nil
	}

	var stages []stationxml.ResponseStageType
	for n, s := range derived.Stage {
		var stage stationxml.ResponseStageType
		if err := clone(&s, &stage); err != nil {
			return nil, err
		}

		if stage.PolesZeros != nil {
			stage.PolesZeros.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}
		if stage.Coefficients != nil {
			stage.Coefficients.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}
		if stage.Polynomial != nil {
			stage.Polynomial.Name = fmt.Sprintf("%sstage_%d", r.prefix, n+1)
		}
		stage.Number = stationxml.CounterType(n + 1)

		stages = append(stages, stage)
	}

	derived.Stage = stages

	return &derived, nil
}
