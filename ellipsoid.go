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

// Cartesian converts polar coordinates φ, λ, and H to Cartesian coordinates x, y, and z.
func (e *Ellipsoid) Cartesian(φ, λ, H float64) (x, y, z float64) {
	sinφ, cosφ := math.Sincos(φ)
	sinλ, cosλ := math.Sincos(λ)
	v := e.a / math.Sqrt(1-e.e2*sinφ*sinφ)
	x = (v + H) * cosφ * cosλ
	y = (v + H) * cosφ * sinλ
	z = ((1-e.e2)*v + H) * sinφ
	return
}

// Polar converts Cartesian coordinates x, y, z to polar coordinates φ, λ, and H with precision prec.
func (e *Ellipsoid) Polar(x, y, z, prec float64) (φ, λ, H float64) {
	λ = math.Atan2(y, x)
	p := math.Hypot(x, y)
	φ0 := math.Atan2(z, p*(1-e.e2))
	for {
		sinφ0 := math.Sin(φ0)
		v := e.a / math.Sqrt(1-e.e2*sinφ0*sinφ0)
		φ = math.Atan2(z+e.e2*v*sinφ0, p)
		if math.Abs(φ-φ0) < prec {
			H = p/math.Cos(φ) - v
			return
		}
		φ0 = φ
	}
}

var (
	Airy1830          = NewEllipsoid("Airy1830", 6377563.396, 6356256.909)
	Airy1830Modified  = NewEllipsoid("Airy1830Modified", 6377340.189, 6356034.447)
	International1924 = NewEllipsoid("International1924", 6378388.000, 6356911.946)
	Hayford1909       = NewEllipsoid("Hayford1909", 6378388.000, 6356911.946)
	GRS80             = NewEllipsoid("GRS80", 6378137.000, 6356752.3141)
	WGS84             = NewEllipsoid("WGS84", 6378137.000, 6356752.3141)
)
