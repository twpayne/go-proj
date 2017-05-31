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
		φ0:   rad(0),
		λ0:   rad(6*(30-float64(zone)) + 3),
		e0:   500000,
		n0:   0,
		e:    International1924,
	}
}

// LatLongUTMZone returns the UTM zone for latitude φ and longitude λ.
func LatLongUTMZone(φ, λ float64) int {
	switch {
	case 56 <= φ && φ < 64 && 3 <= λ && λ < 12:
		return 32
	case 72 <= φ && φ < 84 && 0 <= λ && λ < 42:
		switch {
		case λ < 9:
			return 31
		case λ < 21:
			return 33
		case λ < 33:
			return 35
		default:
			return 37
		}
	default:
		return int(λ+180)/6 + 1
	}
}

// LatLongUTMLetter returns the UTM letter for latitude φ and longitude λ.
func LatLongUTMLetter(φ, λ float64) byte {
	if φ < -80 || 84 <= φ {
		return 0
	}
	i := int((φ + 80) / 8)
	return utmLetters[i]
}

func initUTM() {
	for z := 1; z <= 60; z++ {
		UTMZones[z] = utmZone(z)
	}
}
