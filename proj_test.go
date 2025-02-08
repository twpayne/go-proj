package proj_test

import (
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-proj/v11"
)

func TestRadToDeg(t *testing.T) {
	coord := proj.NewCoord(math.Pi, -math.Pi/4, 1, 2)
	actual := coord.RadToDeg()
	assert.Equal(t, 180., actual.X())
	assert.Equal(t, -45., actual.Y())
	assert.Equal(t, 1., actual.Z())
	assert.Equal(t, 2., actual.M())
}
