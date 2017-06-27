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

// Forward converts latitude lat and longitude lon to easting E and northing N.
func (sm *SphericalMercator) Forward(lat, lon float64) (E, N float64) {
	E = sm.r * lon
	N = sm.r * math.Log(math.Tan((lat+math.Pi/2)/2))
	return
}

// Reverse converts easting E and northing N to latitude lat and longitude lon.
func (sm *SphericalMercator) Reverse(E, N float64) (lat, lon float64) {
	lat = 2*math.Atan(math.Exp(N/sm.r)) - math.Pi/2
	lon = E / sm.r
	return
}
