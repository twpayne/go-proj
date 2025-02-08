package proj_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-proj/v11"
)

func TestCoordsToFloat64Slices(t *testing.T) {
	for i, tc := range []struct {
		coords        []proj.Coord
		float64Slices [][]float64
	}{
		{
			coords:        nil,
			float64Slices: nil,
		},
		{
			coords:        []proj.Coord{},
			float64Slices: [][]float64{},
		},
		{
			coords:        []proj.Coord{{1, 2}},
			float64Slices: [][]float64{{1, 2, 0, 0}},
		},
		{
			coords:        []proj.Coord{{1, 2}, {3, 4}},
			float64Slices: [][]float64{{1, 2, 0, 0}, {3, 4, 0, 0}},
		},
		{
			coords:        []proj.Coord{{1, 2, 3, 4}, {5, 6, 7, 8}},
			float64Slices: [][]float64{{1, 2, 3, 4}, {5, 6, 7, 8}},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.float64Slices, proj.CoordsToFloat64Slices(slices.Clone(tc.coords)))
			assert.Equal(t, tc.coords, proj.Float64SlicesToCoords(slices.Clone(tc.float64Slices)))
		})
	}
}
