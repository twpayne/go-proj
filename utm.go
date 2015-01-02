package proj

import (
	"fmt"
	"sync"
)

var (
	utmZoneCacheMutex sync.RWMutex
	utmZoneCache      = make(map[int]*TransverseMercator)
)

// UTMZone returns the projection for zone z, or nil if no such zone exists.
func UTMZone(z int) *TransverseMercator {
	if z < 1 || 60 < z {
		return nil
	}
	utmZoneCacheMutex.RLock()
	tm, ok := utmZoneCache[z]
	utmZoneCacheMutex.RUnlock()
	if ok {
		return tm
	}
	tm = &TransverseMercator{
		name: fmt.Sprintf("UTMZone(%d)", z),
		f0:   0.9996,
		φ0:   rad(0),
		λ0:   rad(6*(30-float64(z)) + 3),
		e0:   500000,
		n0:   0,
		e:    International1924,
	}
	utmZoneCacheMutex.Lock()
	utmZoneCache[z] = tm
	utmZoneCacheMutex.Unlock()
	return tm
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
func LatLongUTMLetter(φ, λ float64) string {
	if φ < -80 || 84 <= φ {
		return ""
	}
	const letters = "CDEFGHJKLMNPQRSTUVWXX"
	i := int((φ + 80) / 8)
	return letters[i : i+1]
}
