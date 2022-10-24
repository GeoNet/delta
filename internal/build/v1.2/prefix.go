package build

import (
	"fmt"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

func Stages(prefix string, freq float64, stages ...stationxml.ResponseStageType) ([]stationxml.ResponseStageType, error) {
	var res []stationxml.ResponseStageType
	for n, s := range stages {
		var stage stationxml.ResponseStageType
		if err := clone(&s, &stage); err != nil {
			return nil, err
		}

		if stage.PolesZeros != nil {
			stage.PolesZeros.Name = fmt.Sprintf("%sstage_%d", prefix, n+1)
		}
		if stage.Coefficients != nil {
			stage.Coefficients.Name = fmt.Sprintf("%sstage_%d", prefix, n+1)
		}
		if stage.Polynomial != nil {
			stage.Polynomial.Name = fmt.Sprintf("%sstage_%d", prefix, n+1)
		}

		stage.Number = stationxml.CounterType(n + 1)

		if stage.StageGain != nil {
			scale := stage.StageGain.Value
			if stage.PolesZeros != nil {
				g, z := Gain(stage.PolesZeros, freq), Gain(stage.PolesZeros, stage.PolesZeros.NormalizationFrequency.Value)
				stage.PolesZeros.NormalizationFactor = 1.0 / g
				stage.PolesZeros.NormalizationFrequency = stationxml.FrequencyType{stationxml.FloatType{Value: freq}}
				scale /= (z / g)
			}
			stage.StageGain = &stationxml.GainType{
				Value:     scale,
				Frequency: freq,
			}
		}
		res = append(res, stage)

	}
	return res, nil
}

func (r *Response) Normalise() error {
	stages, err := Stages(r.prefix, r.freq, r.Stages()...)
	if err != nil {
		return err
	}
	r.stages = stages
	return nil
}
