// Package proj provides an interface to PROJ. See https://proj.org.
package proj

// #cgo pkg-config: proj
// #include "go-proj.h"
// #cgo nocallback proj_area_create
// #cgo nocallback proj_area_destroy
// #cgo nocallback proj_area_set_bbox
// #cgo noescape proj_area_create
// #cgo noescape proj_area_destroy
// #cgo noescape proj_area_set_bbox
import "C"

import (
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
	cPJArea *C.PJ_AREA
}

// An Error is an error.
type Error struct {
	context *Context
	errno   int
}

// NewArea returns a new Area.
func NewArea(westLonDegree, southLatDegree, eastLonDegree, northLatDegree float64) *Area {
	cPJArea := C.proj_area_create()
	C.proj_area_set_bbox(cPJArea, C.double(westLonDegree), C.double(southLatDegree), C.double(eastLonDegree), C.double(northLatDegree))
	a := &Area{
		cPJArea: cPJArea,
	}
	runtime.AddCleanup(a, func(cPJArea *C.PJ_AREA) {
		C.proj_area_destroy(cPJArea)
	}, cPJArea)
	return a
}

func (e *Error) Error() string {
	return e.context.errnoString(e.errno)
}
