package proj

import (
	"math"
)

// An Ellipsoid represents an ellipsoid.
type Ellipsoid struct {
	name string
	a    float64
	b    float64
	e2   float64
	n    float64
}

// NewEllipsoid creates a new Ellipsoid with name, major axis a, and minor axis b.
func NewEllipsoid(name string, a, b float64) *Ellipsoid {
	return &Ellipsoid{
		name: name,
		a:    a,
		b:    b,
		e2:   (a*a - b*b) / (a * a),
		n:    (a - b) / (a + b),
	}
}

// Cartesian converts polar coordinates lat, lon, and H to Cartesian coordinates x, y, and z.
func (e *Ellipsoid) Cartesian(lat, lon, H float64) (x, y, z float64) {
	sinLat, cosLat := math.Sincos(lat)
	sinLon, cosLon := math.Sincos(lon)
	v := e.a / math.Sqrt(1-e.e2*sinLat*sinLat)
	x = (v + H) * cosLat * cosLon
	y = (v + H) * cosLat * sinLon
	z = ((1-e.e2)*v + H) * sinLat
	return
}

// Polar converts Cartesian coordinates x, y, z to polar coordinates lat, lon, and H with precision prec.
func (e *Ellipsoid) Polar(x, y, z, prec float64) (lat, lon, H float64) {
	lon = math.Atan2(y, x)
	p := math.Hypot(x, y)
	lat0 := math.Atan2(z, p*(1-e.e2))
	for {
		sinLat0 := math.Sin(lat0)
		v := e.a / math.Sqrt(1-e.e2*sinLat0*sinLat0)
		lat = math.Atan2(z+e.e2*v*sinLat0, p)
		if math.Abs(lat-lat0) < prec {
			H = p/math.Cos(lat) - v
			return
		}
		lat0 = lat
	}
}

var (
	// Airy1830 is the Airy 1830 ellipsoid.
	Airy1830 = NewEllipsoid("Airy1830", 6377563.396, 6356256.909)

	// Airy1830Modified is the Airy 1830 Modified ellipsoid.
	Airy1830Modified = NewEllipsoid("Airy1830Modified", 6377340.189, 6356034.447)

	// International1924 is the International 1924 ellipsoid.
	International1924 = NewEllipsoid("International1924", 6378388.000, 6356911.946)

	// Hayford1909 is the Hayford 1909 ellipsoid.
	Hayford1909 = NewEllipsoid("Hayford1909", 6378388.000, 6356911.946)

	// GRS80 is the GRS80 ellipsoid.
	GRS80 = NewEllipsoid("GRS80", 6378137.000, 6356752.3141)

	// WGS84 is the WGS84 ellipsoid.
	WGS84 = NewEllipsoid("WGS84", 6378137.000, 6356752.3141)
)
