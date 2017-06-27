// Package proj implements common geospatial projections.
package proj

import (
	"math"
)

// A T is a generic projection.
type T interface {
	Code() int
	Forward(lat, lon float64) (E float64, N float64)
	Reverse(E, N float64) (lat float64, lon float64)
}

// dms converts d degrees, m minutes, and s seconds to degrees.
func dms(d, m, s float64) float64 {
	return d + m/60 + s/3600
}

// rad converts x from degrees to radians.
func rad(x float64) float64 {
	return math.Pi * x / 180
}

func init() {
	initUTM()
	initEPSG()
}
