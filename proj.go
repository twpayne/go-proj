// Package proj provides an interface to PROJ. See https://proj.org.
package proj

// #cgo pkg-config: proj
// #include "go-proj.h"
import "C"

import (
	"math"
	"runtime"
)

// Version.
const (
	VersionMajor = C.PROJ_VERSION_MAJOR
	VersionMinor = C.PROJ_VERSION_MINOR
	VersionPatch = C.PROJ_VERSION_PATCH
)

// An Area is an area.
type Area struct {
	pjArea *C.PJ_AREA
}

type Bounds struct {
	XMin float64
	YMin float64
	XMax float64
	YMax float64
}

// A coord is a coordinate.
type Coord [4]float64

// An Error is an error.
type Error struct {
	context *Context
	errno   int
}

// NewArea returns a new Area.
func NewArea(westLonDegree, southLatDegree, eastLonDegree, northLatDegree float64) *Area {
	pjArea := C.proj_area_create()
	C.proj_area_set_bbox(pjArea, (C.double)(westLonDegree), (C.double)(southLatDegree), (C.double)(eastLonDegree), (C.double)(northLatDegree))
	a := &Area{
		pjArea: pjArea,
	}
	runtime.SetFinalizer(a, (*Area).Destroy)
	return a
}

// Destroy frees all resources associated with a.
func (a *Area) Destroy() {
	if a.pjArea != nil {
		C.proj_area_destroy(a.pjArea)
		a.pjArea = nil
	}
}

// NewCoord returns a new Coord.
func NewCoord(x, y, z, m float64) Coord {
	return Coord{x, y, z, m}
}

// DegToRad returns a new Coord with the first two elements transformed from
// degrees to radians.
func (c Coord) DegToRad() Coord {
	return Coord{math.Pi * c[0] / 180, math.Pi * c[1] / 180, c[2], c[3]}
}

// RadToDeg returns a new Coord with the first two elements transformed from
// radians to degrees.
func (c Coord) RadToDeg() Coord {
	return Coord{180 * c[0] / math.Pi, 180 * c[1] / math.Pi, c[2], c[3]}
}

// X returns c's X coordinate.
func (c *Coord) X() float64 { return c[0] }

// Y returns c's Y coordinate.
func (c *Coord) Y() float64 { return c[1] }

// Z returns c's Z coordinate.
func (c *Coord) Z() float64 { return c[2] }

// M returns c's M coordinate.
func (c *Coord) M() float64 { return c[3] }

func (e *Error) Error() string {
	return e.context.errnoString(e.errno)
}
