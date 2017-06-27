package proj

import (
	"testing"
)

func TestSphericalMercator(t *testing.T) {
	for _, c := range []struct {
		sm   *SphericalMercator
		φ, λ float64
		e, n float64
	}{
		{
			sm: EPSG3857,
			φ:  rad(dms(52, 39, 27.2531)),
			λ:  rad(dms(1, 43, 4.5177)),
			e:  191238.15587944098,
			n:  6919907.180756924,
		},
		{
			sm: EPSG3857,
			φ:  rad(10),
			λ:  rad(50),
			e:  5565974.539664,
			n:  1118889.974858,
		},
		{
			sm: EPSG3857,
			φ:  rad(45.677),
			λ:  rad(-111.0429),
			e:  -12361239.084208,
			n:  5728738.469095,
		},
	} {
		if e, n := c.sm.Forward(c.φ, c.λ); !near(e, c.e, 1e-4) || !near(n, c.n, 1e-3) {
			t.Errorf("%v.Forward(%v, %v) == %v, %v, want %v, %v", c.sm, c.φ, c.λ, e, n, c.e, c.n)
		}
		if φ, λ := c.sm.Reverse(c.e, c.n); !near(φ, c.φ, 1e-10) || !near(λ, c.λ, 1e-10) {
			t.Errorf("%v.Reverse(%v, %v) == %v, %v, want %v, %v", c.sm, c.e, c.n, φ, λ, c.φ, c.λ)
		}
	}
}
