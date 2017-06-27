package proj

const (
	utmLetters = "CDEFGHJKLMNPQRSTUVWXX"
)

var (
	// UTMZones is a map from UTM zones to projections.
	UTMZones = map[int]*TransverseMercator{}
)

// utmZone returns the UTM projection for zone.
func utmZone(zone int) *TransverseMercator {
	return &TransverseMercator{
		code: 32600 + zone,
		f0:   0.9996,
		lat0: rad(0),
		lon0: rad(6*(30-float64(zone)) + 3),
		e0:   500000,
		n0:   0,
		e:    International1924,
	}
}

// LatLongUTMZone returns the UTM zone for latitude lat and longitude lon.
func LatLongUTMZone(lat, lon float64) int {
	switch {
	case 56 <= lat && lat < 64 && 3 <= lon && lon < 12:
		return 32
	case 72 <= lat && lat < 84 && 0 <= lon && lon < 42:
		switch {
		case lon < 9:
			return 31
		case lon < 21:
			return 33
		case lon < 33:
			return 35
		default:
			return 37
		}
	default:
		return int(lon+180)/6 + 1
	}
}

// LatLongUTMLetter returns the UTM letter for latitude lat and longitude lon.
func LatLongUTMLetter(lat, lon float64) byte {
	if lat < -80 || 84 <= lat {
		return 0
	}
	i := int((lat + 80) / 8)
	return utmLetters[i]
}

func initUTM() {
	for z := 1; z <= 60; z++ {
		UTMZones[z] = utmZone(z)
	}
}
