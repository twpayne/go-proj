package proj_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-proj/v10"
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

func TestTransFloat64Slices(t *testing.T) {
	for i, tc := range []struct {
		float64Slices [][]float64
		expected      [][]float64
		delta         float64
	}{
		{
			float64Slices: nil,
			expected:      nil,
		},
		{
			float64Slices: [][]float64{},
			expected:      [][]float64{},
		},
		{
			float64Slices: [][]float64{
				{48.856613, 2.352222, 78},
				{40.712778, -74.006111, 10},
			},
			expected: [][]float64{
				{261848.15527273554, 6250566.54904563, 78},
				{-8238322.592110482, 4970068.348185822, 10},
			},
		},
		{
			float64Slices: [][]float64{
				{48.856613, 2.352222, 78, 1},
				{40.712778, -74.006111, 10, 2},
			},
			expected: [][]float64{
				{261848.15527273554, 6250566.54904563, 78, 1},
				{-8238322.592110482, 4970068.348185822, 10, 2},
			},
		},
		{
			float64Slices: [][]float64{
				{48.856613, 2.352222},
				{40.712778, -74.006111, 10, 2},
			},
			expected: [][]float64{
				{261848.15527273554, 6250566.54904563},
				{-8238322.592110482, 4970068.348185822, 10, 2},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			pj, err := proj.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
			assert.NoError(t, err)
			float64Slices := slices.Clone(tc.float64Slices)
			assert.NoError(t, pj.ForwardFloat64Slices(float64Slices))
			assertInDeltaFloat64Slices(t, tc.expected, float64Slices, tc.delta)
		})
	}
}

func assertInDeltaFloat64Slices(tb testing.TB, expected, actual [][]float64, delta float64) {
	tb.Helper()
	assert.Equal(tb, len(expected), len(actual))
	for i := range expected {
		assertInDeltaFloat64Slice(tb, expected[i], actual[i], delta)
	}
}
