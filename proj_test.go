package proj

import (
	"math"
)

// near returns true iff the absolute difference between x and y is less than e.
func near(x, y, e float64) bool {
	return math.Abs(x-y) < e
}
