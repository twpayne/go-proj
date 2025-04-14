package proj

// #include <stdlib.h>
// #include "go-proj.h"
// #cgo nocallback proj_context_create
// #cgo nocallback proj_context_destroy
// #cgo nocallback proj_context_errno
// #cgo nocallback proj_context_errno_string
// #cgo nocallback proj_context_set_search_paths
// #cgo nocallback proj_create
// #cgo nocallback proj_create_argv
// #cgo nocallback proj_create_crs_to_crs
// #cgo nocallback proj_create_crs_to_crs_from_pj
// #cgo nocallback proj_destroy
// #cgo nocallback proj_log_level
// #cgo noescape proj_context_create
// #cgo noescape proj_context_destroy
// #cgo noescape proj_context_errno
// #cgo noescape proj_context_errno_string
// #cgo noescape proj_context_set_search_paths
// #cgo noescape proj_create
// #cgo noescape proj_create_argv
// #cgo noescape proj_create_crs_to_crs
// #cgo noescape proj_create_crs_to_crs_from_pj
// #cgo noescape proj_destroy
// #cgo noescape proj_log_level
import "C"

import (
	"runtime"
	"sync"
	"unsafe"
)

type LogLevel C.PJ_LOG_LEVEL

const (
	LogLevelNone  LogLevel = C.PJ_LOG_NONE
	LogLevelError LogLevel = C.PJ_LOG_ERROR
	LogLevelDebug LogLevel = C.PJ_LOG_DEBUG
	LogLevelTrace LogLevel = C.PJ_LOG_TRACE
	LogLevelTell  LogLevel = C.PJ_LOG_TELL
)

var defaultContext = &Context{}

func init() {
	C.proj_log_level(nil, C.PJ_LOG_NONE)
}

// A Context is a context.
type Context struct {
	mutex      sync.Mutex
	cPJContext *C.PJ_CONTEXT
}

// NewContext returns a new Context.
func NewContext() *Context {
	pjContext := C.proj_context_create()
	C.proj_log_level(pjContext, C.PJ_LOG_NONE)
	c := &Context{
		cPJContext: pjContext,
	}
	runtime.AddCleanup(c, func(pjContext *C.PJ_CONTEXT) {
		C.proj_context_destroy(pjContext)
	}, pjContext)
	return c
}

// SetLogLevel sets the log level.
func (c *Context) SetLogLevel(logLevel LogLevel) {
	c.Lock()
	defer c.Unlock()
	C.proj_log_level(c.cPJContext, C.PJ_LOG_LEVEL(logLevel))
}

// SetSearchPaths sets the paths PROJ should be exploring to find the PROJ Data files.
func (c *Context) SetSearchPaths(paths []string) {
	c.Lock()
	defer c.Unlock()
	cPaths := make([]*C.char, len(paths))
	var pathPtr unsafe.Pointer
	for i, path := range paths {
		cPaths[i] = C.CString(path)
		defer C.free(unsafe.Pointer(cPaths[i]))
	}
	if len(paths) > 0 {
		pathPtr = unsafe.Pointer(&cPaths[0])
	}
	C.proj_context_set_search_paths(c.cPJContext, C.int(len(cPaths)), (**C.char)(pathPtr))
}

func (c *Context) Lock() {
	c.mutex.Lock()
}

// NewCRSToCRS returns a new PJ from sourceCRS to targetCRS and optional area.
func (c *Context) NewCRSToCRS(sourceCRS, targetCRS string, area *Area) (*PJ, error) {
	c.Lock()
	defer c.Unlock()

	cSourceCRS := C.CString(sourceCRS)
	defer C.free(unsafe.Pointer(cSourceCRS))

	cTargetCRS := C.CString(targetCRS)
	defer C.free(unsafe.Pointer(cTargetCRS))

	var cArea *C.PJ_AREA
	if area != nil {
		cArea = area.cPJArea
	}

	return c.newPJ(C.proj_create_crs_to_crs(c.cPJContext, cSourceCRS, cTargetCRS, cArea))
}

// NewCRSToCRSFromPJ returns a new PJ from two CRSs.
func (c *Context) NewCRSToCRSFromPJ(sourcePJ, targetPJ *PJ, area *Area, options string) (*PJ, error) {
	c.Lock()
	defer c.Unlock()

	if sourcePJ.context != c {
		sourcePJ.context.Lock()
		defer sourcePJ.context.Unlock()
	}

	if targetPJ.context != c && targetPJ.context != sourcePJ.context {
		targetPJ.context.Lock()
		defer targetPJ.context.Unlock()
	}

	var cOptionsPtr **C.char
	if options != "" {
		cOptions := C.CString(options)
		defer C.free(unsafe.Pointer(cOptions))
		cOptionsPtr = &cOptions
	}

	var cArea *C.PJ_AREA
	if area != nil {
		cArea = area.cPJArea
	}

	return c.newPJ(C.proj_create_crs_to_crs_from_pj(c.cPJContext, sourcePJ.cPJ, targetPJ.cPJ, cArea, cOptionsPtr))
}

// New returns a new PJ with the given definition.
func (c *Context) New(definition string) (*PJ, error) {
	c.Lock()
	defer c.Unlock()

	cDefinition := C.CString(definition)
	defer C.free(unsafe.Pointer(cDefinition))

	return c.newPJ(C.proj_create(c.cPJContext, cDefinition))
}

// NewFromArgs returns a new PJ from args.
func (c *Context) NewFromArgs(args ...string) (*PJ, error) {
	c.Lock()
	defer c.Unlock()

	cArgs := make([]*C.char, len(args))
	for i := range cArgs {
		cArg := C.CString(args[i])
		defer C.free(unsafe.Pointer(cArg))
		cArgs[i] = cArg
	}

	return c.newPJ(C.proj_create_argv(c.cPJContext, (C.int)(len(cArgs)), (**C.char)(unsafe.Pointer(&cArgs[0]))))
}

func (c *Context) Unlock() {
	c.mutex.Unlock()
}

// errnoString returns the text representation of errno.
func (c *Context) errnoString(errno int) string {
	c.Lock()
	defer c.Unlock()
	return C.GoString(C.proj_context_errno_string(c.cPJContext, (C.int)(errno)))
}

// newError returns a new error with number errno.
func (c *Context) newError(errno int) *Error {
	return &Error{
		context: c,
		errno:   errno,
	}
}

// newPJ returns a new PJ or an error.
func (c *Context) newPJ(cPJ *C.PJ) (*PJ, error) {
	if cPJ == nil {
		return nil, c.newError(int(C.proj_context_errno(c.cPJContext)))
	}

	pj := &PJ{
		context: c,
		cPJ:     cPJ,
	}
	runtime.AddCleanup(pj, func(cPJ *C.PJ) {
		C.proj_destroy(cPJ)
	}, cPJ)
	return pj, nil
}

// SetLogLevel sets the log level for the default context.
func SetLogLevel(logLevel LogLevel) {
	defaultContext.SetLogLevel(logLevel)
}

// New returns a PJ with the given definition.
func New(definition string) (*PJ, error) {
	return defaultContext.New(definition)
}

// NewFromArgs returns a PJ with the given args.
func NewFromArgs(args ...string) (*PJ, error) {
	return defaultContext.NewFromArgs(args...)
}

// NewCRSToCRS returns a new PJ from sourceCRS to targetCRS and optional area.
func NewCRSToCRS(sourceCRS, targetCRS string, area *Area) (*PJ, error) {
	return defaultContext.NewCRSToCRS(sourceCRS, targetCRS, area)
}

// NewCRSToCRSFromPJ returns a new PJ from two CRSs.
func NewCRSToCRSFromPJ(sourcePJ, targetPJ *PJ, area *Area, options string) (*PJ, error) {
	return defaultContext.NewCRSToCRSFromPJ(sourcePJ, targetPJ, area, options)
}
