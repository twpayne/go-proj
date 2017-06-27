package proj

import (
	"testing"
)

func TestSphericalMercator(t *testing.T) {
	for _, c := range []struct {
		sm       *SphericalMercator
		lat, lon float64
		e, n     float64
	}{
		{
			sm:  EPSG3857,
			lat: rad(dms(52, 39, 27.2531)),
			lon: rad(dms(1, 43, 4.5177)),
			e:   191238.15587944098,
			n:   6919907.180756924,
		},
		{
			sm:  EPSG3857,
			lat: rad(10),
			lon: rad(50),
			e:   5565974.539664,
			n:   1118889.974858,
		},
		{
			sm:  EPSG3857,
			lat: rad(45.677),
			lon: rad(-111.0429),
			e:   -12361239.084208,
			n:   5728738.469095,
		},
	} {
		if e, n := c.sm.Forward(c.lat, c.lon); !near(e, c.e, 1e-4) || !near(n, c.n, 1e-3) {
			t.Errorf("%v.Forward(%v, %v) == %v, %v, want %v, %v", c.sm, c.lat, c.lon, e, n, c.e, c.n)
		}
		if lat, lon := c.sm.Reverse(c.e, c.n); !near(lat, c.lat, 1e-10) || !near(lon, c.lon, 1e-10) {
			t.Errorf("%v.Reverse(%v, %v) == %v, %v, want %v, %v", c.sm, c.e, c.n, lat, lon, c.lat, c.lon)
		}
	}
}

func TestEPSG3857(t *testing.T) {
	t.Skip()
	for _, c := range epsg3857TestData {
		if e, n := EPSG3857.Forward(c.lat, c.lon); !near(e, c.e, 1e-4) || !near(n, c.n, 1e-3) {
			t.Errorf("EPSG3857.Forward(%v, %v) == %v, %v, want %v, %v", c.lat, c.lon, e, n, c.e, c.n)
		}
		if lat, lon := EPSG3857.Reverse(c.e, c.n); !near(lat, c.lat, 1e-10) || !near(lon, c.lon, 1e-10) {
			t.Errorf("EPSG3857.Reverse(%v, %v) == %v, %v, want %v, %v", c.e, c.n, lat, lon, c.lat, c.lon)
		}
	}
}
