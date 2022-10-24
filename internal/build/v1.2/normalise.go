package build

import (
	"fmt"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

func (r *Response) Normalise() error {

	var stages []stationxml.ResponseStageType
	for n, s := range append(r.sensor.Stage, r.datalogger.Stage...) {
		var stage stationxml.ResponseStageType
		if err := clone(&s, &stage); err != nil {
			return err
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

		if stage.StageGain != nil {
			scale := stage.StageGain.Value
			if stage.PolesZeros != nil {
				g, z := Gain(stage.PolesZeros, r.freq), Gain(stage.PolesZeros, stage.PolesZeros.NormalizationFrequency.Value)
				stage.PolesZeros.NormalizationFactor = 1.0 / g
				stage.PolesZeros.NormalizationFrequency = stationxml.FrequencyType{stationxml.FloatType{Value: r.freq}}
				scale /= (z / g)
			}
			stage.StageGain = &stationxml.GainType{
				Value:     scale,
				Frequency: r.freq,
			}
		}
		stages = append(stages, stage)
	}

	r.stages = stages

	return nil
}
