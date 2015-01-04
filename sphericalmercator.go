package proj

import (
	"math"
)

var (
	EPSG3857 = NewSphericalMercator("EPSG3857", 6378137)
)

// A SphericalMercator represents a spherical Mercator projection.
type SphericalMercator struct {
	name string
	r    float64
}

// NewSphericalMercator returns a new SphericalMercator with the given name and radius r.
func NewSphericalMercator(name string, r float64) *SphericalMercator {
	return &SphericalMercator{
		name: name,
		r:    r,
	}
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

// String returns the name of the projection.
func (sm *SphericalMercator) String() string {
	return sm.name
}
