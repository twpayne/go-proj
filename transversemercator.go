package proj

import (
	"math"
)

// A TransverseMercator represents a transverse Mercator projection.
type TransverseMercator struct {
	code       int
	f0         float64
	lat0, lon0 float64
	e0, n0     float64
	e          *Ellipsoid
}

// Code returns tm's EPSG code.
func (tm *TransverseMercator) Code() int {
	return tm.code
}

// Forward converts latitude lat and longitude lon to easting E and northing N.
func (tm *TransverseMercator) Forward(lat, lon float64) (E, N float64) {
	sinLat, cosLat := math.Sincos(lat)
	nu := tm.e.a * tm.f0 / math.Sqrt(1-tm.e.e2*sinLat*sinLat)
	rho := tm.e.a * tm.f0 * (1 - tm.e.e2) * math.Pow(1-tm.e.e2*sinLat*sinLat, -1.5)
	eta2 := nu/rho - 1
	n := tm.e.n
	n2 := n * n
	M := tm.e.b * tm.f0 * ((1+n+5*n2/4+5*n*n2/4)*(lat-tm.lat0) - (3*n+3*n2+21*n*n2/8)*math.Sin((lat-tm.lat0))*math.Cos(lat+tm.lat0) + (15*n2/8+15*n*n2/8)*math.Sin(2*(lat-tm.lat0))*math.Cos(2*(lat+tm.lat0)) - (35*n*n2/24)*math.Sin(3*(lat-tm.lat0))*math.Cos(3*(lat+tm.lat0)))
	I := M + tm.n0
	II := nu * sinLat * cosLat / 2
	cosLat2 := cosLat * cosLat
	cosLat4 := cosLat2 * cosLat2
	tanLat := math.Tan(lat)
	tanLat2 := tanLat * tanLat
	tanLat4 := tanLat2 * tanLat2
	III := nu * sinLat * cosLat * cosLat2 * (5 - tanLat2 + 9*eta2) / 24
	IIIA := nu * sinLat * cosLat * cosLat4 * (61 - 58*tanLat2 + tanLat4) / 720
	IV := nu * cosLat
	V := nu * cosLat * cosLat2 * (nu/rho - tanLat2) / 6
	VI := nu * cosLat * cosLat4 * (5 - 18*tanLat2 + tanLat4 + 14*eta2 - 58*tanLat2*eta2) / 120
	deltaLon := lon - tm.lon0
	deltaLon2 := deltaLon * deltaLon
	deltaLon4 := deltaLon2 * deltaLon2
	N = I + II*deltaLon2 + III*deltaLon4 + IIIA*deltaLon2*deltaLon4
	E = tm.e0 + IV*deltaLon + V*deltaLon*deltaLon2 + VI*deltaLon*deltaLon4
	return
}

// Reverse converts easting E and northing N to latitude lat and longitude lon.
func (tm *TransverseMercator) Reverse(E, N float64) (lat, lon float64) {
	lat1 := (N-tm.e.n)/(tm.e.a*tm.f0) + tm.lat0
	n := tm.e.n
	n2 := n * n
	var M float64
	for {
		M = tm.e.b * tm.f0 * ((1+n+5*n2/4+5*n*n2/4)*(lat1-tm.lat0) - (3*n+3*n2+21*n*n2/8)*math.Sin((lat1-tm.lat0))*math.Cos(lat1+tm.lat0) + (15*n2/8+15*n*n2/8)*math.Sin(2*(lat1-tm.lat0))*math.Cos(2*(lat1+tm.lat0)) - (35*n*n2/24)*math.Sin(3*(lat1-tm.lat0))*math.Cos(3*(lat1+tm.lat0)))
		if math.Abs(N-tm.n0-M) < 1e-8 {
			break
		} else {
			lat1 += (N - tm.n0 - M) / (tm.e.a * tm.f0)
		}
	}
	sinLat1, cosLat1 := math.Sincos(lat1)
	nu := tm.e.a * tm.f0 / math.Sqrt(1-tm.e.e2*sinLat1*sinLat1)
	rho := tm.e.a * tm.f0 * (1 - tm.e.e2) * math.Pow(1-tm.e.e2*sinLat1*sinLat1, -1.5)
	eta2 := nu/rho - 1
	tanLat1 := math.Tan(lat1)
	tanLat12 := tanLat1 * tanLat1
	tanLat14 := tanLat12 * tanLat12
	VII := tanLat1 / (2 * rho * nu)
	nu2 := nu * nu
	nu4 := nu2 * nu2
	VIII := tanLat1 * (5 + 3*tanLat12 + eta2 - 9*tanLat12*eta2) / (24 * rho * nu * nu2)
	IX := tanLat1 * (61 + 90*tanLat12 + 45*tanLat12*tanLat12) / (720 * rho * nu * nu4)
	secLat1 := 1 / cosLat1
	X := secLat1 / nu
	XI := secLat1 * (nu/rho + 2*tanLat12) / (6 * nu * nu2)
	XII := secLat1 * (5 + 28*tanLat12 + 24*tanLat14) / (120 * nu * nu4)
	XIIA := secLat1 * (61 + 662*tanLat12 + 1320*tanLat14 + 720*tanLat12*tanLat14) / (5040 * nu * nu2 * nu4)
	deltaE := E - tm.e0
	deltaE2 := deltaE * deltaE
	deltaE4 := deltaE2 * deltaE2
	lat = lat1 - VII*deltaE2 + VIII*deltaE4 - IX*deltaE2*deltaE4
	lon = tm.lon0 + X*deltaE - XI*deltaE*deltaE2 + XII*deltaE*deltaE4 - XIIA*deltaE*deltaE2*deltaE4
	return
}
