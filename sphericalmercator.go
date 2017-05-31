package proj

import (
	"math"
)

// A SphericalMercator represents a spherical Mercator projection.
type SphericalMercator struct {
	code int
	r    float64
}

// Code returns sm's EPSG code.
func (sm *SphericalMercator) Code() int {
	return sm.code
}

// Forward converts latitude φ and longitude λ to easting E and northing N.
func (sm *SphericalMercator) Forward(φ, λ float64) (E, N float64) {
	E = sm.r * λ
	N = sm.r * math.Log(math.Tan((φ+math.Pi/2)/2))
	return
}

// Reverse converts easting E and northing N to latitude φ and longitude λ.
func (sm *SphericalMercator) Reverse(E, N float64) (φ, λ float64) {
	φ = 2*math.Atan(math.Exp(N/sm.r)) - math.Pi/2
	λ = E / sm.r
	return
}
