package build

import (
	"math"
	"math/cmplx"

	"github.com/GeoNet/delta/internal/stationxml/v1.2"
)

func Gain(pz *stationxml.PolesZerosType, freq float64) float64 {

	var w complex128
	switch pz.PzTransferFunctionType {
	case stationxml.LaplaceRadiansSecondPzTransferFunction:
		w = complex(0.0, 2.0*math.Pi*freq)
	default:
		w = complex(0.0, freq)
	}

	h := complex(float64(1.0), float64(0.0))

	for _, zero := range pz.Zero {
		h *= (w - complex(zero.Real.Value, zero.Imaginary.Value))
	}

	for _, pole := range pz.Pole {
		h /= (w - complex(pole.Real.Value, pole.Imaginary.Value))
	}

	return cmplx.Abs(h)
}
