package proj

// #include <stdlib.h>
// #include "go-proj.h"
import "C"

import (
	"runtime"
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

// A Transformation is a transformation.
type Transformation struct {
	context *Context
	pj      *C.PJ
}

// NewCRSToCRSTransformation returns a new Transformation from sourceCRS to
// targetCRS and optional area.
func (c *Context) NewCRSToCRSTransformation(sourceCRS, targetCRS string, area *Area) (*Transformation, error) {
	c.Lock()
	defer c.Unlock()

	cSourceCRS := C.CString(sourceCRS)
	defer C.free(unsafe.Pointer(cSourceCRS))

	cTargetCRS := C.CString(targetCRS)
	defer C.free(unsafe.Pointer(cTargetCRS))

	var cArea *C.PJ_AREA
	if area != nil {
		cArea = area.pjArea
	}

	pj := C.proj_create_crs_to_crs(c.pjContext, cSourceCRS, cTargetCRS, cArea)
	if pj == nil {
		return nil, c.newError(int(C.proj_context_errno(c.pjContext)))
	}

	transformation := &Transformation{
		context: c,
		pj:      pj,
	}
	runtime.SetFinalizer(transformation, (*Transformation).Destroy)
	return transformation, nil
}

// Destroy releases all resources associated with t.
func (t *Transformation) Destroy() {
	t.context.Lock()
	defer t.context.Unlock()
	if t.pj != nil {
		C.proj_destroy(t.pj)
		t.pj = nil
	}
}

// Forward transforms coord in the forward direction.
func (t *Transformation) Forward(coord Coord) (Coord, error) {
	return t.Trans(DirectionFwd, coord)
}

// ForwardArray transforms coorsd in the forward direction.
func (t *Transformation) ForwardArray(coords []Coord) error {
	return t.TransArray(DirectionFwd, coords)
}

// ForwardFlatCoords transforms flatCoords in the forward direction.
func (t *Transformation) ForwardFlatCoords(flatCoords []float64, stride, zIndex, mIndex int) error {
	return t.TransFlatCoords(DirectionFwd, flatCoords, stride, zIndex, mIndex)
}

// Inverse transforms coord in the inverse direction.
func (t *Transformation) Inverse(coord Coord) (Coord, error) {
	return t.Trans(DirectionInv, coord)
}

// InverseArray transforms coorsd in the inverse direction.
func (t *Transformation) InverseArray(coords []Coord) error {
	return t.TransArray(DirectionInv, coords)
}

// InverseFlatCoords transforms flatCoords in the inverse direction.
func (t *Transformation) InverseFlatCoords(flatCoords []float64, stride, zIndex, mIndex int) error {
	return t.TransFlatCoords(DirectionInv, flatCoords, stride, zIndex, mIndex)
}

// Trans transforms a single Coord.
func (t *Transformation) Trans(direction Direction, coord Coord) (Coord, error) {
	t.context.Lock()
	defer t.context.Unlock()

	lastErrno := C.proj_errno_reset(t.pj)
	defer C.proj_errno_restore(t.pj, lastErrno)

	pjCoord := C.proj_trans(t.pj, (C.PJ_DIRECTION)(direction), *(*C.PJ_COORD)(unsafe.Pointer(&coord)))
	if errno := int(C.proj_errno(t.pj)); errno != 0 {
		return Coord{}, t.context.newError(errno)
	}
	return *(*Coord)(unsafe.Pointer(&pjCoord)), nil
}

// TransArray transforms an array of Coords.
func (t *Transformation) TransArray(direction Direction, coords []Coord) error {
	if len(coords) == 0 {
		return nil
	}

	t.context.Lock()
	defer t.context.Unlock()

	lastErrno := C.proj_errno_reset(t.pj)
	defer C.proj_errno_restore(t.pj, lastErrno)

	if errno := int(C.proj_trans_array(t.pj, (C.PJ_DIRECTION)(direction), (C.ulong)(len(coords)), (*C.PJ_COORD)(unsafe.Pointer(&coords[0])))); errno != 0 {
		return t.context.newError(errno)
	}
	return nil
}

// TransFlatCoords transforms an array of flat coordinates.
func (t *Transformation) TransFlatCoords(direction Direction, flatCoords []float64, stride, zIndex, mIndex int) error {
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

	return t.TransGeneric(direction, x, sx, nx, y, sy, ny, z, sz, nz, m, sm, nm)
}

// TransGeneric transforms a series of coordinates.
func (t *Transformation) TransGeneric(direction Direction, x *float64, sx, nx int, y *float64, sy, ny int, z *float64, sz, nz int, m *float64, sm, nm int) error {
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

	t.context.Lock()
	defer t.context.Unlock()

	lastErrno := C.proj_errno_reset(t.pj)
	defer C.proj_errno_restore(t.pj, lastErrno)

	if int(C.proj_trans_generic(t.pj, (C.PJ_DIRECTION)(direction),
		(*C.double)(x), C.size_t(sx), C.size_t(nx),
		(*C.double)(y), C.size_t(sy), C.size_t(ny),
		(*C.double)(z), C.size_t(sz), C.size_t(nz),
		(*C.double)(m), C.size_t(sm), C.size_t(nm),
	)) != maxN {
		return t.context.newError(int(C.proj_errno(t.pj)))
	}

	return nil
}
