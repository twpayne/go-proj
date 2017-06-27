package proj

var (
	// EPSG3857 is EPSG:3857 (WGS 84 / Pseudo-Mercator).
	EPSG3857 = &SphericalMercator{
		code: 3857,
		r:    6378137,
	}

	// EPSG27700 is EPSG:27700 (OSGB 1936 / British National Grid).
	EPSG27700 = &TransverseMercator{
		code: 27700,
		f0:   0.9996012717,
		lat0: rad(49),
		lon0: rad(-2),
		e0:   400000,
		n0:   -100000,
		e:    Airy1830,
	}

	// EPSG29903 is EPSG:29903 (TM75 / Irish National Grid).
	EPSG29903 = &TransverseMercator{
		code: 29903,
		f0:   1.000035,
		lat0: rad(dms(53, 30, 0)),
		lon0: rad(-8),
		e0:   200000,
		n0:   250000,
		e:    Airy1830Modified,
	}

	// EPSG is a map from EPSG codes to projections.
	EPSG = map[int]T{
		3857:  EPSG3857,
		27700: EPSG27700,
		29903: EPSG29903,
	}
)

func initEPSG() {
	for _, p := range UTMZones {
		EPSG[p.Code()] = p
	}
}
