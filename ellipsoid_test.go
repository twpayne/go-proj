package proj

import (
	"testing"
)

func TestEllipsoid(t *testing.T) {
	for _, c := range []struct {
		e              *Ellipsoid
		phi, lambda, H float64
		x, y, z        float64
	}{
		{
			e:      Airy1830,
			phi:    rad(dms(52, 39, 27.2531)),
			lambda: rad(dms(1, 43, 4.5177)),
			H:      24.700,
			x:      3874938.849,
			y:      116218.624,
			z:      5047168.208,
		},
	} {
		if x, y, z := c.e.Cartesian(c.phi, c.lambda, c.H); !near(x, c.x, 1e-3) || !near(y, c.y, 1e-3) || !near(z, c.z, 1e-3) {
			t.Errorf("%s.Cartesian(%v, %v, %v) == %v, %v, %v, want %v, %v, %v", c.e, c.phi, c.lambda, c.H, x, y, z, c.x, c.y, c.z)
		}
		if phi, lambda, H := c.e.Polar(c.x, c.y, c.z, 1e-10); !near(phi, c.phi, 1e-9) || !near(lambda, c.lambda, 1e-10) || !near(H, c.H, 1e-3) {
			t.Errorf("%s.Polar(%v, %v, %v) == %v, %v, %v, want %v, %v, %v", c.e, c.x, c.y, c.z, phi, lambda, H, c.phi, c.lambda, c.H)
		}
	}
}
