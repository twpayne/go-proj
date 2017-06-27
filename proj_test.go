package proj

import (
	"math"
)

// deg converts x from radians to degrees.
func deg(x float64) float64 {
	return 180 * x / math.Pi
}

// near returns true iff the absolute difference between x and y is less than e.
func near(x, y, e float64) bool {
	return math.Abs(x-y) < e
}
