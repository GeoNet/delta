package build

import (
	"encoding/xml"
	"fmt"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

/*
func Prefix(prefix string, count int, stage *stationxml.ResponseStageType) {
	if stage.PolesZeros != nil {
		stage.PolesZeros.Name = fmt.Sprintf("%sstage_%d", prefix, count)
	}
	if stage.Coefficients != nil {
		stage.Coefficients.Name = fmt.Sprintf("%sstage_%d", prefix, count)
	}
	if stage.Polynomial != nil {
		stage.Polynomial.Name = fmt.Sprintf("%sstage_%d", prefix, count)
	}

	stage.Number = stationxml.CounterType(count)
}
*/

/*
func (r *Response) Restage(stage *stationxml.ResponseStageType, count int) *stationxml.ResponseStageType {
	if stage.PolesZeros != nil {
		stage.PolesZeros.Name = fmt.Sprintf("%sstage_%d", r.prefix, count)
	}
	if stage.Coefficients != nil {
		stage.Coefficients.Name = fmt.Sprintf("%sstage_%d", r.prefix, count)
	}
	if stage.Polynomial != nil {
		stage.Polynomial.Name = fmt.Sprintf("%sstage_%d", r.prefix, count)
	}
	stage.Number = stationxml.CounterType(count)
	return stage
}
*/

/*
func Restage(prefix string, stages ...stationxml.ResponseStageType) ([]stationxml.ResponseStageType, error) {
	var res []stationxml.ResponseStageType
	for n, s := range stages {
		var stage stationxml.ResponseStageType
		if err := clone(&s, &stage); err != nil {
			return nil, err
		}
		Prefix(prefix, n+1, &stage)
		res = append(res, stage)
	}
	return res, nil
}

func Derived(prefix string, data []byte) (*stationxml.ResponseType, error) {

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
		Prefix(prefix, n+1, &stage)
		stages = append(stages, stage)
	}

	/*
		stages, err := Restage(prefix, derived.Stage...)
		if err != nil {
			return nil, err
		}
*/
/*

	derived.Stage = stages

	return &derived, nil
}
*/

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
