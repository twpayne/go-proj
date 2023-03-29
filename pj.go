package proj

// #include "go-proj.h"
import "C"

import (
	"unsafe"
)

// A Direction is a direction.
type Direction C.PJ_DIRECTION

// Directions.
const (
	DirectionFwd   Direction = C.PJ_FWD
	DirectionIdent Direction = C.PJ_IDENT
	DirectionInv   Direction = C.PJ_INV
)

// A PJ is a projection or a transformation.
type PJ struct {
	context *Context
	pj      *C.PJ
}

// A PJInfo contains information about a Proj.
type PJInfo struct {
	ID          string
	Description string
	Definition  string
	HasInverse  bool
	Accuracy    float64
}

// Destroy releases all resources associated with p.
func (p *PJ) Destroy() {
	p.context.Lock()
	defer p.context.Unlock()
	if p.pj != nil {
		C.proj_destroy(p.pj)
		p.pj = nil
	}
}

// Forward transforms coord in the forward direction.
func (p *PJ) Forward(coord Coord) (Coord, error) {
	return p.Trans(DirectionFwd, coord)
}

// ForwardBounds transforms bounds in the forward direction.
func (p *PJ) ForwardBounds(bounds Bounds, densifyPoints int) (Bounds, error) {
	return p.TransBounds(DirectionFwd, bounds, densifyPoints)
}

// ForwardArray transforms coorsd in the forward direction.
func (p *PJ) ForwardArray(coords []Coord) error {
	return p.TransArray(DirectionFwd, coords)
}

// ForwardFlatCoords transforms flatCoords in the forward direction.
func (p *PJ) ForwardFlatCoords(flatCoords []float64, stride, zIndex, mIndex int) error {
	return p.TransFlatCoords(DirectionFwd, flatCoords, stride, zIndex, mIndex)
}

// Geod returns the distance, forward azimuth, and reverse azimuth between a and b.
func (p *PJ) Geod(a, b Coord) (float64, float64, float64) {
	p.context.Lock()
	defer p.context.Unlock()
	cCoord := C.proj_geod(p.pj, *(*C.PJ_COORD)(unsafe.Pointer(&a)), *(*C.PJ_COORD)(unsafe.Pointer(&b)))
	cGeod := *(*C.PJ_GEOD)(unsafe.Pointer(&cCoord))
	return (float64)(cGeod.s), (float64)(cGeod.a1), (float64)(cGeod.a2)
}

// GetLastUsedOperation returns the operation used in the last call to Trans.
func (p *PJ) GetLastUsedOperation() (*PJ, error) {
	p.context.Lock()
	defer p.context.Unlock()
	return p.context.newPJ(C.proj_trans_get_last_used_operation(p.pj))
}

// Info returns information about p.
func (p *PJ) Info() PJInfo {
	p.context.Lock()
	defer p.context.Unlock()

	cProjInfo := C.proj_pj_info(p.pj)
	return PJInfo{
		ID:          C.GoString(cProjInfo.id),
		Description: C.GoString(cProjInfo.description),
		Definition:  C.GoString(cProjInfo.definition),
		HasInverse:  cProjInfo.has_inverse != 0,
		Accuracy:    (float64)(cProjInfo.accuracy),
	}
}

// IsCRS returns whether p is a CRS.
func (p *PJ) IsCRS() bool {
	return C.proj_is_crs(p.pj) != 0
}

// Inverse transforms coord in the inverse direction.
func (p *PJ) Inverse(coord Coord) (Coord, error) {
	return p.Trans(DirectionInv, coord)
}

// InverseArray transforms coorsd in the inverse direction.
func (p *PJ) InverseArray(coords []Coord) error {
	return p.TransArray(DirectionInv, coords)
}

// InverseBounds transforms bounds in the forward direction.
func (p *PJ) InverseBounds(bounds Bounds, densifyPoints int) (Bounds, error) {
	return p.TransBounds(DirectionInv, bounds, densifyPoints)
}

// InverseFlatCoords transforms flatCoords in the inverse direction.
func (p *PJ) InverseFlatCoords(flatCoords []float64, stride, zIndex, mIndex int) error {
	return p.TransFlatCoords(DirectionInv, flatCoords, stride, zIndex, mIndex)
}

// LPDist returns the geodesic distance between a and b in geodetic coordinates.
func (p *PJ) LPDist(a, b Coord) float64 {
	p.context.Lock()
	defer p.context.Unlock()
	return (float64)(C.proj_lp_dist(p.pj, *(*C.PJ_COORD)(unsafe.Pointer(&a)), *(*C.PJ_COORD)(unsafe.Pointer(&b))))
}

// LPZDist returns the geodesic distance between a and b in geodetic
// coordinates, taking height above the ellipsoid into account.q
func (p *PJ) LPZDist(a, b Coord) float64 {
	p.context.Lock()
	defer p.context.Unlock()
	return (float64)(C.proj_lpz_dist(p.pj, *(*C.PJ_COORD)(unsafe.Pointer(&a)), *(*C.PJ_COORD)(unsafe.Pointer(&b))))
}

// Trans transforms a single Coord.
func (p *PJ) Trans(direction Direction, coord Coord) (Coord, error) {
	p.context.Lock()
	defer p.context.Unlock()

	lastErrno := C.proj_errno_reset(p.pj)
	defer C.proj_errno_restore(p.pj, lastErrno)

	pjCoord := C.proj_trans(p.pj, (C.PJ_DIRECTION)(direction), *(*C.PJ_COORD)(unsafe.Pointer(&coord)))
	if errno := int(C.proj_errno(p.pj)); errno != 0 {
		return Coord{}, p.context.newError(errno)
	}
	return *(*Coord)(unsafe.Pointer(&pjCoord)), nil
}

// TransArray transforms an array of Coords.
func (p *PJ) TransArray(direction Direction, coords []Coord) error {
	if len(coords) == 0 {
		return nil
	}

	p.context.Lock()
	defer p.context.Unlock()

	lastErrno := C.proj_errno_reset(p.pj)
	defer C.proj_errno_restore(p.pj, lastErrno)

	if errno := int(C.proj_trans_array(p.pj, (C.PJ_DIRECTION)(direction), (C.ulong)(len(coords)), (*C.PJ_COORD)(unsafe.Pointer(&coords[0])))); errno != 0 {
		return p.context.newError(errno)
	}
	return nil
}

func (p *PJ) TransBounds(direction Direction, bounds Bounds, densifyPoints int) (Bounds, error) {
	p.context.Lock()
	defer p.context.Unlock()

	var transBounds Bounds
	if C.proj_trans_bounds(p.context.pjContext, p.pj, (C.PJ_DIRECTION)(direction),
		(C.double)(bounds.XMin), (C.double)(bounds.YMin), (C.double)(bounds.XMax), (C.double)(bounds.YMax),
		(*C.double)(&transBounds.XMin), (*C.double)(&transBounds.YMin), (*C.double)(&transBounds.XMax), (*C.double)(&transBounds.YMax),
		C.int(densifyPoints)) == 0 {
		return Bounds{}, p.context.newError(int(C.proj_errno(p.pj)))
	}
	return transBounds, nil
}

// TransFlatCoords transforms an array of flat coordinates.
func (p *PJ) TransFlatCoords(direction Direction, flatCoords []float64, stride, zIndex, mIndex int) error {
	if len(flatCoords) == 0 {
		return nil
	}
	n := len(flatCoords) / stride

	var x, y, z, m *float64
	var sx, sy, sz, sm int
	var nx, ny, nz, nm int

	x = &flatCoords[0]
	y = &flatCoords[1]
	sx = 8 * stride
	sy = 8 * stride
	nx = n
	ny = n

	if zIndex != -1 {
		z = &flatCoords[zIndex]
		sz = 8 * stride
		nz = n
	}

	if mIndex != -1 {
		m = &flatCoords[mIndex]
		sm = 8 * stride
		nm = n
	}

	return p.TransGeneric(direction, x, sx, nx, y, sy, ny, z, sz, nz, m, sm, nm)
}

// TransGeneric transforms a series of coordinates.
func (p *PJ) TransGeneric(direction Direction, x *float64, sx, nx int, y *float64, sy, ny int, z *float64, sz, nz int, m *float64, sm, nm int) error {
	maxN := nx
	if ny > maxN {
		maxN = ny
	}
	if nz > maxN {
		maxN = nz
	}
	if nm > maxN {
		maxN = nm
	}

	p.context.Lock()
	defer p.context.Unlock()

	lastErrno := C.proj_errno_reset(p.pj)
	defer C.proj_errno_restore(p.pj, lastErrno)

	if int(C.proj_trans_generic(p.pj, (C.PJ_DIRECTION)(direction),
		(*C.double)(x), C.size_t(sx), C.size_t(nx),
		(*C.double)(y), C.size_t(sy), C.size_t(ny),
		(*C.double)(z), C.size_t(sz), C.size_t(nz),
		(*C.double)(m), C.size_t(sm), C.size_t(nm),
	)) != maxN {
		return p.context.newError(int(C.proj_errno(p.pj)))
	}

	return nil
}
