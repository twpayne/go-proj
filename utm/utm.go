package utm

import (
	"errors"
	"fmt"
	"sync"

	"github.com/twpayne/go-proj/v11"
)

// A Coord is a UTM coordinate.
type Coord struct {
	Zone int
	E    float64
	N    float64
}

var (
	errUndefined = errors.New("undefined")

	transformationCache sync.Map
)

func (c Coord) String() string {
	if c.N >= 0 {
		return fmt.Sprintf("%dN %d %d", c.Zone, int(c.E+0.5), int(c.N+0.5))
	}
	return fmt.Sprintf("%dS %d %d", c.Zone, int(c.E+0.5), int(c.N-0.5)+10_000_000)
}

// Forward returns the forward transformation of (lon, lat) to UTM.
func Forward(lon, lat float64) (Coord, error) {
	zone := Zone(lon, lat)
	if zone < 0 {
		return Coord{}, errUndefined
	}
	pj, err := ZoneTransformation(zone)
	if err != nil {
		return Coord{}, err
	}
	utmCoord, err := pj.Forward(proj.NewCoord(lat, lon, 0, 0))
	if err != nil {
		return Coord{}, err
	}
	return Coord{
		Zone: zone,
		E:    utmCoord.X(),
		N:    utmCoord.Y(),
	}, nil
}

// Zone returns the UTM zone at lon and lat, or -1 if there is no zone at that
// coordinate.
func Zone(lon, lat float64) int {
	// See https://en.wikipedia.org/wiki/Universal_Transverse_Mercator_coordinate_system#Exceptions.
	switch {
	case lat < -80 || 84 <= lat:
		return -1
	case 56 <= lat && lat < 64 && 3 <= lon && lon < 9:
		return 32
	case 72 <= lat && lat < 84:
		switch {
		case 0 <= lon && lon < 9:
			return 31
		case 9 <= lon && lon < 21:
			return 33
		case 21 <= lon && lon < 30:
			return 35
		}
	}
	return int((180+lon)/6) + 1
}

// ZoneTransformation returns the transformation from EPSG:4326 to the given UTM
// zone.
func ZoneTransformation(zone int) (*proj.PJ, error) {
	if pj, ok := transformationCache.Load(zone); ok {
		return pj.(*proj.PJ), nil //nolint:forcetypeassert
	}
	pj, err := proj.NewCRSToCRS("epsg:4326", fmt.Sprintf("+proj=utm +zone=%d", zone), nil)
	if err != nil {
		return nil, err
	}
	actual, _ := transformationCache.LoadOrStore(zone, pj)
	return actual.(*proj.PJ), nil //nolint:forcetypeassert
}
