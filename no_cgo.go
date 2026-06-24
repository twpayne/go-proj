//go:build !cgo

package proj

const noCGoMessage = "package proj requires CGo and the PROJ C library to be installed.\n" +
	"Install libproj-dev (Debian/Ubuntu), proj-devel (RHEL/Fedora), proj (Alpine), or proj (Homebrew/macOS).\n" +
	"See https://proj.org for more information."

func init() {
	panic(noCGoMessage)
}

// proj.go stubs

const (
	VersionMajor = 0
	VersionMinor = 0
	VersionPatch = 0
)

type Area struct{}

func NewArea(westLonDegree, southLatDegree, eastLonDegree, northLatDegree float64) *Area {
	panic(noCGoMessage)
}

type Error struct{}

func (e *Error) Error() string {
	panic(noCGoMessage)
}

// context.go stubs

type LogLevel int

const (
	LogLevelNone  LogLevel = 0
	LogLevelError LogLevel = 1
	LogLevelDebug LogLevel = 2
	LogLevelTrace LogLevel = 3
	LogLevelTell  LogLevel = 4
)

type Context struct{}

func NewContext() *Context {
	panic(noCGoMessage)
}

func (c *Context) SetLogLevel(logLevel LogLevel) {
	panic(noCGoMessage)
}

func (c *Context) SetSearchPaths(paths []string) {
	panic(noCGoMessage)
}

func (c *Context) Lock() {
	panic(noCGoMessage)
}

func (c *Context) Unlock() {
	panic(noCGoMessage)
}

func (c *Context) NewCRSToCRS(sourceCRS, targetCRS string, area *Area) (*PJ, error) {
	panic(noCGoMessage)
}

func (c *Context) NewCRSToCRSFromPJ(sourcePJ, targetPJ *PJ, area *Area, options string) (*PJ, error) {
	panic(noCGoMessage)
}

func (c *Context) New(definition string) (*PJ, error) {
	panic(noCGoMessage)
}

func (c *Context) NewFromArgs(args ...string) (*PJ, error) {
	panic(noCGoMessage)
}

func SetLogLevel(logLevel LogLevel) {
	panic(noCGoMessage)
}

func New(definition string) (*PJ, error) {
	panic(noCGoMessage)
}

func NewFromArgs(args ...string) (*PJ, error) {
	panic(noCGoMessage)
}

func NewCRSToCRS(sourceCRS, targetCRS string, area *Area) (*PJ, error) {
	panic(noCGoMessage)
}

func NewCRSToCRSFromPJ(sourcePJ, targetPJ *PJ, area *Area, options string) (*PJ, error) {
	panic(noCGoMessage)
}

// pj.go stubs

type Direction int

const (
	DirectionFwd   Direction = 1
	DirectionIdent Direction = 0
	DirectionInv   Direction = -1
)

type PJ struct{}

type PJInfo struct {
	ID          string
	Description string
	Definition  string
	HasInverse  bool
	Accuracy    float64
}

func (pj *PJ) NormalizeForVisualization() (*PJ, error) {
	panic(noCGoMessage)
}

func (pj *PJ) Forward(coord Coord) (Coord, error) {
	panic(noCGoMessage)
}

func (pj *PJ) ForwardBounds(bounds Bounds, densifyPoints int) (Bounds, error) {
	panic(noCGoMessage)
}

func (pj *PJ) ForwardArray(coords []Coord) error {
	panic(noCGoMessage)
}

func (pj *PJ) ForwardFlatCoords(flatCoords []float64, stride, zIndex, mIndex int) error {
	panic(noCGoMessage)
}

func (pj *PJ) ForwardFloat64Slice(float64Slice []float64) ([]float64, error) {
	panic(noCGoMessage)
}

func (pj *PJ) ForwardFloat64Slices(float64Slices [][]float64) error {
	panic(noCGoMessage)
}

func (pj *PJ) Geod(a, b Coord) (float64, float64, float64) {
	panic(noCGoMessage)
}

func (pj *PJ) GetLastUsedOperation() (*PJ, error) {
	panic(noCGoMessage)
}

func (pj *PJ) Info() PJInfo {
	panic(noCGoMessage)
}

func (pj *PJ) IsCRS() bool {
	panic(noCGoMessage)
}

func (pj *PJ) Inverse(coord Coord) (Coord, error) {
	panic(noCGoMessage)
}

func (pj *PJ) InverseArray(coords []Coord) error {
	panic(noCGoMessage)
}

func (pj *PJ) InverseBounds(bounds Bounds, densifyPoints int) (Bounds, error) {
	panic(noCGoMessage)
}

func (pj *PJ) InverseFlatCoords(flatCoords []float64, stride, zIndex, mIndex int) error {
	panic(noCGoMessage)
}

func (pj *PJ) InverseFloat64Slice(float64Slice []float64) ([]float64, error) {
	panic(noCGoMessage)
}

func (pj *PJ) InverseFloat64Slices(float64Slices [][]float64) error {
	panic(noCGoMessage)
}

func (pj *PJ) LPDist(a, b Coord) float64 {
	panic(noCGoMessage)
}

func (pj *PJ) LPZDist(a, b Coord) float64 {
	panic(noCGoMessage)
}

func (pj *PJ) Trans(direction Direction, coord Coord) (Coord, error) {
	panic(noCGoMessage)
}

func (pj *PJ) TransArray(direction Direction, coords []Coord) error {
	panic(noCGoMessage)
}

func (pj *PJ) TransBounds(direction Direction, bounds Bounds, densifyPoints int) (Bounds, error) {
	panic(noCGoMessage)
}

func (pj *PJ) TransFlatCoords(direction Direction, flatCoords []float64, stride, zIndex, mIndex int) error {
	panic(noCGoMessage)
}

func (pj *PJ) TransFloat64Slice(direction Direction, float64Slice []float64) ([]float64, error) {
	panic(noCGoMessage)
}

func (pj *PJ) TransFloat64Slices(direction Direction, float64Slices [][]float64) error {
	panic(noCGoMessage)
}

func (pj *PJ) TransGeneric(direction Direction, x *float64, sx, nx int, y *float64, sy, ny int, z *float64, sz, nz int, m *float64, sm, nm int) error {
	panic(noCGoMessage)
}
