package proj

import (
	"testing"
)

func TestTransverseMercator(t *testing.T) {
	for _, c := range []struct {
		tm   *TransverseMercator
		φ, λ float64
		e, n float64
	}{
		{
			tm: NationalGrid,
			φ:  rad(dms(52, 39, 27.2531)),
			λ:  rad(dms(1, 43, 4.5177)),
			e:  651409.903,
			n:  313177.270,
		},
	} {
		if e, n := c.tm.Forward(c.φ, c.λ); !near(e, c.e, 1e-4) || !near(n, c.n, 1e-3) {
			t.Errorf("%v.Forward(%v, %v) == %v, %v, want %v, %v", c.tm, c.φ, c.λ, e, n, c.e, c.n)
		}
		if φ, λ := c.tm.Reverse(c.e, c.n); !near(φ, c.φ, 1e-10) || !near(λ, c.λ, 1e-10) {
			t.Errorf("%v.Reverse(%v, %v) == %v, %v, want %v, %v", c.tm, c.e, c.n, φ, λ, c.φ, c.λ)
		}
	}
}
