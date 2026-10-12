package main

import (
	"encoding/json"
	"fmt"
	"strings"

	stationxml "github.com/GeoNet/delta/internal/stationxml/v1.2"
)

type PoleZero complex128

func (pz PoleZero) MarshalJSON() ([]byte, error) {
	if imag(pz) < 0 {
		return json.Marshal(fmt.Sprintf("(%g%gi)", real(pz), imag(pz)))
	}
	return json.Marshal(fmt.Sprintf("(%g+%gi)", real(pz), imag(pz)))
}

// PolesZeros holds the response details for an analogue sensor.
type PolesZeros struct {
	Name string `json:"name"`

	TransferFunction       string  `json:"transfer_function"`
	NormalizationFactor    float64 `json:"normalization_factor"`
	NormalizationFrequency float64 `json:"normalization_frequency"`

	Poles []PoleZero `json:"poles"`
	Zeros []PoleZero `json:"zeros"`
}

func GetPolesZeros(stage stationxml.ResponseStageType) (PolesZeros, bool) {
	if stage.PolesZeros == nil {
		return PolesZeros{}, false
	}

	var poles, zeros []PoleZero
	for _, p := range stage.PolesZeros.Pole {
		poles = append(poles, PoleZero(complex(p.Real.Value, p.Imaginary.Value)))
	}
	for _, z := range stage.PolesZeros.Zero {
		zeros = append(zeros, PoleZero(complex(z.Real.Value, z.Imaginary.Value)))
	}

	name := strings.TrimPrefix(stage.PolesZeros.ResourceId, "PolesZeros#")

	paz := PolesZeros{
		Name: name,

		TransferFunction:       stage.PolesZeros.PzTransferFunctionType.String(),
		NormalizationFactor:    stage.PolesZeros.NormalizationFactor,
		NormalizationFrequency: stage.PolesZeros.NormalizationFrequency.Value,

		Poles: poles,
		Zeros: zeros,
	}

	return paz, true
}
