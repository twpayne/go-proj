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
	cPJ     *C.PJ
}

// A PJInfo contains information about a PJ.
type PJInfo struct {
	ID          string
	Description string
	Definition  string
	HasInverse  bool
	Accuracy    float64
}

// Returns a new PJ instance whose axis order is the one expected for
// visualization purposes. If the axis order of its source or target CRS is
// northing, easting, then an axis swap operation will be inserted.
//
// The axis order of geographic CRS will be longitude, latitude[, height], and
// the one of projected CRS will be easting, northing [, height].
func (pj *PJ) NormalizeForVisualization() (*PJ, error) {
	pj.context.Lock()
	defer pj.context.Unlock()
	return pj.context.newPJ(C.proj_normalize_for_visualization(pj.context.cPJContext, pj.cPJ))
}

// Forward transforms coord in the forward direction.
func (pj *PJ) Forward(coord Coord) (Coord, error) {
	return pj.Trans(DirectionFwd, coord)
}

// ForwardBounds transforms bounds in the forward direction.
func (pj *PJ) ForwardBounds(bounds Bounds, densifyPoints int) (Bounds, error) {
	return pj.TransBounds(DirectionFwd, bounds, densifyPoints)
}

// ForwardArray transforms coords in the forward direction.
func (pj *PJ) ForwardArray(coords []Coord) error {
	return pj.TransArray(DirectionFwd, coords)
}

// ForwardFlatCoords transforms flatCoords in the forward direction.
func (pj *PJ) ForwardFlatCoords(flatCoords []float64, stride, zIndex, mIndex int) error {
	return pj.TransFlatCoords(DirectionFwd, flatCoords, stride, zIndex, mIndex)
}

// ForwardFloat64Slice transforms float64 in place in the forward direction.
func (pj *PJ) ForwardFloat64Slice(float64Slice []float64) ([]float64, error) {
	return pj.TransFloat64Slice(DirectionFwd, float64Slice)
}

// ForwardFloat64Slices transforms float64Slices in the forward direction.
func (pj *PJ) ForwardFloat64Slices(float64Slices [][]float64) error {
	return pj.TransFloat64Slices(DirectionFwd, float64Slices)
}

// Geod returns the distance, forward azimuth, and reverse azimuth between a and b.
func (pj *PJ) Geod(a, b Coord) (float64, float64, float64) {
	pj.context.Lock()
	defer pj.context.Unlock()
	cCoord := C.proj_geod(pj.cPJ, *(*C.PJ_COORD)(unsafe.Pointer(&a)), *(*C.PJ_COORD)(unsafe.Pointer(&b)))
	cGeod := *(*C.PJ_GEOD)(unsafe.Pointer(&cCoord))
	return (float64)(cGeod.s), (float64)(cGeod.a1), (float64)(cGeod.a2)
}

// GetLastUsedOperation returns the operation used in the last call to Trans.
func (pj *PJ) GetLastUsedOperation() (*PJ, error) {
	pj.context.Lock()
	defer pj.context.Unlock()
	return pj.context.newPJ(C.proj_trans_get_last_used_operation(pj.cPJ))
}

// Info returns information about pj.
func (pj *PJ) Info() PJInfo {
	pj.context.Lock()
	defer pj.context.Unlock()

	cProjInfo := C.proj_pj_info(pj.cPJ)
	return PJInfo{
		ID:          C.GoString(cProjInfo.id),
		Description: C.GoString(cProjInfo.description),
		Definition:  C.GoString(cProjInfo.definition),
		HasInverse:  cProjInfo.has_inverse != 0,
		Accuracy:    (float64)(cProjInfo.accuracy),
	}
}

// IsCRS returns whether pj is a CRS.
func (pj *PJ) IsCRS() bool {
	return C.proj_is_crs(pj.cPJ) != 0
}

// Inverse transforms coord in the inverse direction.
func (pj *PJ) Inverse(coord Coord) (Coord, error) {
	return pj.Trans(DirectionInv, coord)
}

// InverseArray transforms coords in the inverse direction.
func (pj *PJ) InverseArray(coords []Coord) error {
	return pj.TransArray(DirectionInv, coords)
}

// InverseBounds transforms bounds in the forward direction.
func (pj *PJ) InverseBounds(bounds Bounds, densifyPoints int) (Bounds, error) {
	return pj.TransBounds(DirectionInv, bounds, densifyPoints)
}

// InverseFlatCoords transforms flatCoords in the inverse direction.
func (pj *PJ) InverseFlatCoords(flatCoords []float64, stride, zIndex, mIndex int) error {
	return pj.TransFlatCoords(DirectionInv, flatCoords, stride, zIndex, mIndex)
}

// InverseFloat64Slice transforms float64 in place in the forward direction.
func (pj *PJ) InverseFloat64Slice(float64Slice []float64) ([]float64, error) {
	return pj.TransFloat64Slice(DirectionInv, float64Slice)
}

// InverseFloat64Slices transforms float64Slices in the inverse direction.
func (pj *PJ) InverseFloat64Slices(float64Slices [][]float64) error {
	return pj.TransFloat64Slices(DirectionInv, float64Slices)
}

// LPDist returns the geodesic distance between a and b in geodetic coordinates.
func (pj *PJ) LPDist(a, b Coord) float64 {
	pj.context.Lock()
	defer pj.context.Unlock()
	return (float64)(C.proj_lp_dist(pj.cPJ, *(*C.PJ_COORD)(unsafe.Pointer(&a)), *(*C.PJ_COORD)(unsafe.Pointer(&b))))
}

// LPZDist returns the geodesic distance between a and b in geodetic
// coordinates, taking height above the ellipsoid into account.
func (pj *PJ) LPZDist(a, b Coord) float64 {
	pj.context.Lock()
	defer pj.context.Unlock()
	return (float64)(C.proj_lpz_dist(pj.cPJ, *(*C.PJ_COORD)(unsafe.Pointer(&a)), *(*C.PJ_COORD)(unsafe.Pointer(&b))))
}

// Trans transforms a single Coord in place.
func (pj *PJ) Trans(direction Direction, coord Coord) (Coord, error) {
	pj.context.Lock()
	defer pj.context.Unlock()

	lastErrno := C.proj_errno_reset(pj.cPJ)
	defer C.proj_errno_restore(pj.cPJ, lastErrno)

	pjCoord := C.proj_trans(pj.cPJ, (C.PJ_DIRECTION)(direction), *(*C.PJ_COORD)(unsafe.Pointer(&coord)))
	if errno := int(C.proj_errno(pj.cPJ)); errno != 0 {
		return Coord{}, pj.context.newError(errno)
	}
	return *(*Coord)(unsafe.Pointer(&pjCoord)), nil
}

// TransArray transforms an array of Coords.
func (pj *PJ) TransArray(direction Direction, coords []Coord) error {
	if len(coords) == 0 {
		return nil
	}

	pj.context.Lock()
	defer pj.context.Unlock()

	lastErrno := C.proj_errno_reset(pj.cPJ)
	defer C.proj_errno_restore(pj.cPJ, lastErrno)

	if errno := int(C.proj_trans_array(pj.cPJ, (C.PJ_DIRECTION)(direction), (C.size_t)(len(coords)), (*C.PJ_COORD)(unsafe.Pointer(&coords[0])))); errno != 0 {
		return pj.context.newError(errno)
	}
	return nil
}

// TransBounds transforms bounds.
func (pj *PJ) TransBounds(direction Direction, bounds Bounds, densifyPoints int) (Bounds, error) {
	pj.context.Lock()
	defer pj.context.Unlock()

	var transBounds Bounds
	if C.proj_trans_bounds(pj.context.cPJContext, pj.cPJ, (C.PJ_DIRECTION)(direction),
		(C.double)(bounds.XMin), (C.double)(bounds.YMin), (C.double)(bounds.XMax), (C.double)(bounds.YMax),
		(*C.double)(&transBounds.XMin), (*C.double)(&transBounds.YMin), (*C.double)(&transBounds.XMax), (*C.double)(&transBounds.YMax),
		C.int(densifyPoints)) == 0 {
		return Bounds{}, pj.context.newError(int(C.proj_errno(pj.cPJ)))
	}
	return transBounds, nil
}

// TransFlatCoords transforms an array of flat coordinates.
func (pj *PJ) TransFlatCoords(direction Direction, flatCoords []float64, stride, zIndex, mIndex int) error {
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

	return pj.TransGeneric(direction, x, sx, nx, y, sy, ny, z, sz, nz, m, sm, nm)
}

// TransFloat64Slice transforms a []float64 in place.
func (pj *PJ) TransFloat64Slice(direction Direction, float64Slice []float64) ([]float64, error) {
	var coord Coord
	copy(coord[:], float64Slice)
	transCoord, err := pj.Trans(direction, coord)
	if err != nil {
		return nil, err
	}
	copy(float64Slice, transCoord[:])
	return float64Slice, nil
}

// TransFloat64Slices transforms float64Slices.
func (pj *PJ) TransFloat64Slices(direction Direction, float64Slices [][]float64) error {
	coords := Float64SlicesToCoords(float64Slices)
	if err := pj.TransArray(direction, coords); err != nil {
		return err
	}
	for i, coord := range coords {
		copy(float64Slices[i], coord[:])
	}
	return nil
}

// TransGeneric transforms a series of coordinates.
func (pj *PJ) TransGeneric(direction Direction, x *float64, sx, nx int, y *float64, sy, ny int, z *float64, sz, nz int, m *float64, sm, nm int) error {
	pj.context.Lock()
	defer pj.context.Unlock()

	lastErrno := C.proj_errno_reset(pj.cPJ)
	defer C.proj_errno_restore(pj.cPJ, lastErrno)

	if int(C.proj_trans_generic(pj.cPJ, (C.PJ_DIRECTION)(direction),
		(*C.double)(x), C.size_t(sx), C.size_t(nx),
		(*C.double)(y), C.size_t(sy), C.size_t(ny),
		(*C.double)(z), C.size_t(sz), C.size_t(nz),
		(*C.double)(m), C.size_t(sm), C.size_t(nm),
	)) != max(nx, ny, nz, nm) {
		return pj.context.newError(int(C.proj_errno(pj.cPJ)))
	}

	return nil
}
