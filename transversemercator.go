package proj

import (
	"math"
)

// A TransverseMercator represents a transverse Mercator projection.
type TransverseMercator struct {
	code   int
	f0     float64
	φ0, λ0 float64
	e0, n0 float64
	e      *Ellipsoid
}

// Code returns tm's EPSG code.
func (tm *TransverseMercator) Code() int {
	return tm.code
}

// Forward converts latitude φ and longitude λ to easting E and northing N.
func (tm *TransverseMercator) Forward(φ, λ float64) (E, N float64) {
	sinφ, cosφ := math.Sincos(φ)
	ν := tm.e.a * tm.f0 / math.Sqrt(1-tm.e.e2*sinφ*sinφ)
	ρ := tm.e.a * tm.f0 * (1 - tm.e.e2) * math.Pow(1-tm.e.e2*sinφ*sinφ, -1.5)
	η2 := ν/ρ - 1
	n := tm.e.n
	n2 := n * n
	M := tm.e.b * tm.f0 * ((1+n+5*n2/4+5*n*n2/4)*(φ-tm.φ0) - (3*n+3*n2+21*n*n2/8)*math.Sin((φ-tm.φ0))*math.Cos(φ+tm.φ0) + (15*n2/8+15*n*n2/8)*math.Sin(2*(φ-tm.φ0))*math.Cos(2*(φ+tm.φ0)) - (35*n*n2/24)*math.Sin(3*(φ-tm.φ0))*math.Cos(3*(φ+tm.φ0)))
	I := M + tm.n0
	II := ν * sinφ * cosφ / 2
	cosφ2 := cosφ * cosφ
	cosφ4 := cosφ2 * cosφ2
	tanφ := math.Tan(φ)
	tanφ2 := tanφ * tanφ
	tanφ4 := tanφ2 * tanφ2
	III := ν * sinφ * cosφ * cosφ2 * (5 - tanφ2 + 9*η2) / 24
	IIIA := ν * sinφ * cosφ * cosφ4 * (61 - 58*tanφ2 + tanφ4) / 720
	IV := ν * cosφ
	V := ν * cosφ * cosφ2 * (ν/ρ - tanφ2) / 6
	VI := ν * cosφ * cosφ4 * (5 - 18*tanφ2 + tanφ4 + 14*η2 - 58*tanφ2*η2) / 120
	δλ := λ - tm.λ0
	δλ2 := δλ * δλ
	δλ4 := δλ2 * δλ2
	N = I + II*δλ2 + III*δλ4 + IIIA*δλ2*δλ4
	E = tm.e0 + IV*δλ + V*δλ*δλ2 + VI*δλ*δλ4
	return
}

// Reverse converts easting E and northing N to latitude φ and longitude λ.
func (tm *TransverseMercator) Reverse(E, N float64) (φ, λ float64) {
	φ1 := (N-tm.e.n)/(tm.e.a*tm.f0) + tm.φ0
	n := tm.e.n
	n2 := n * n
	var M float64
	for {
		M = tm.e.b * tm.f0 * ((1+n+5*n2/4+5*n*n2/4)*(φ1-tm.φ0) - (3*n+3*n2+21*n*n2/8)*math.Sin((φ1-tm.φ0))*math.Cos(φ1+tm.φ0) + (15*n2/8+15*n*n2/8)*math.Sin(2*(φ1-tm.φ0))*math.Cos(2*(φ1+tm.φ0)) - (35*n*n2/24)*math.Sin(3*(φ1-tm.φ0))*math.Cos(3*(φ1+tm.φ0)))
		if math.Abs(N-tm.n0-M) < 1e-8 {
			break
		} else {
			φ1 += (N - tm.n0 - M) / (tm.e.a * tm.f0)
		}
	}
	sinφ1, cosφ1 := math.Sincos(φ1)
	ν := tm.e.a * tm.f0 / math.Sqrt(1-tm.e.e2*sinφ1*sinφ1)
	ρ := tm.e.a * tm.f0 * (1 - tm.e.e2) * math.Pow(1-tm.e.e2*sinφ1*sinφ1, -1.5)
	η2 := ν/ρ - 1
	tanφ1 := math.Tan(φ1)
	tanφ12 := tanφ1 * tanφ1
	tanφ14 := tanφ12 * tanφ12
	VII := tanφ1 / (2 * ρ * ν)
	ν2 := ν * ν
	ν4 := ν2 * ν2
	VIII := tanφ1 * (5 + 3*tanφ12 + η2 - 9*tanφ12*η2) / (24 * ρ * ν * ν2)
	IX := tanφ1 * (61 + 90*tanφ12 + 45*tanφ12*tanφ12) / (720 * ρ * ν * ν4)
	secφ1 := 1 / cosφ1
	X := secφ1 / ν
	XI := secφ1 * (ν/ρ + 2*tanφ12) / (6 * ν * ν2)
	XII := secφ1 * (5 + 28*tanφ12 + 24*tanφ14) / (120 * ν * ν4)
	XIIA := secφ1 * (61 + 662*tanφ12 + 1320*tanφ14 + 720*tanφ12*tanφ14) / (5040 * ν * ν2 * ν4)
	δE := E - tm.e0
	δE2 := δE * δE
	δE4 := δE2 * δE2
	φ = φ1 - VII*δE2 + VIII*δE4 - IX*δE2*δE4
	λ = tm.λ0 + X*δE - XI*δE*δE2 + XII*δE*δE4 - XIIA*δE*δE2*δE4
	return
}
