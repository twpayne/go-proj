// Package proj implements the algorithms described in "A Guide to Coordinate Systems in Great Britain".
package proj

import (
	"math"
)

// dms converts d degrees, m minutes, and s seconds to degrees.
func dms(d, m, s float64) float64 {
	return d + m/60 + s/3600
}

// deg converts x from radians to degrees.
func deg(x float64) float64 {
	return 180 * x / math.Pi
}

// rad converts x from degrees to radians.
func rad(x float64) float64 {
	return math.Pi * x / 180
}

// near returns true iff the absolute difference between x and y is less than e.
func near(x, y, e float64) bool {
	return math.Abs(x-y) < e
}
