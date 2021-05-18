package main

import (
	"math"
)

func WGS842ITRF(lat, lon, elev float64) (float64, float64, float64) {

	a, b := 6378137.0, 6356752.31424518
	e := math.Sqrt((a*a - b*b) / (a * a))

	phi, lambda := math.Pi*lat/180.0, math.Pi*lon/180.0

	N := a / math.Sqrt(float64(1.0)-e*e*math.Sin(phi)*math.Sin(phi))

	x := (N + elev) * math.Cos(phi) * math.Cos(lambda)
	y := (N + elev) * math.Cos(phi) * math.Sin(lambda)
	z := (((b * b * N) / (a * a)) + elev) * math.Sin(phi)

	return x, y, z
}
