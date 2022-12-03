// Package proj provides an interface to PROJ. See https://proj.org.
package proj

// #cgo pkg-config: proj
// #include "go-proj.h"
import "C"

import (
	"runtime"
	"sync"
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

// A Context is a context.
type Context struct {
	sync.Mutex
	pjContext *C.PJ_CONTEXT
}

// A coord is a coordinate.
type Coord [4]float64

// An Error is an error.
type Error struct {
	context *Context
	errno   int
}

var defaultContext = &Context{}

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

// NewContext returns a new Context.
func NewContext() *Context {
	pjContext := C.proj_context_create()
	C.proj_log_level(pjContext, C.PJ_LOG_NONE)
	c := &Context{
		pjContext: pjContext,
	}
	runtime.SetFinalizer(c, (*Context).Destroy)
	return c
}

// Destroy frees all resources associated with c.
func (c *Context) Destroy() {
	c.Lock()
	defer c.Unlock()
	if c.pjContext != nil {
		C.proj_context_destroy(c.pjContext)
		c.pjContext = nil
	}
}

// errnoString returns the text representation of errno.
func (c *Context) errnoString(errno int) string {
	c.Lock()
	defer c.Unlock()
	return C.GoString(C.proj_context_errno_string(c.pjContext, (C.int)(errno)))
}

// newError returns a new error with number errno.
func (c *Context) newError(errno int) *Error {
	return &Error{
		context: c,
		errno:   errno,
	}
}

// newTransformation returns a new *Transformation or an error.
func (c *Context) newTransformation(pj *C.PJ) (*Transformation, error) {
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

// NewCoord returns a new Coord.
func NewCoord(x, y, z, m float64) Coord {
	return Coord{x, y, z, m}
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
