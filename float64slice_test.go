package proj_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-proj/v10"
)

func TestTransFloat64Slice(t *testing.T) {
	for i, tc := range []struct {
		float64Slice []float64
		expected     []float64
		delta        float64
	}{
		{
			float64Slice: nil,
			expected:     nil,
		},
		{
			float64Slice: []float64{},
			expected:     []float64{},
		},
		{
			float64Slice: []float64{723134.1266446244, 474831.4869142064},
			expected:     []float64{54.371652, 18.612462},
			delta:        1e-14,
		},
		{
			float64Slice: []float64{723134.1266446244, 474831.4869142064, 11.1},
			expected:     []float64{54.371652, 18.612462, 11.1},
			delta:        1e-14,
		},
		{
			float64Slice: []float64{723134.1266446244, 474831.4869142064, 11.1, 1},
			expected:     []float64{54.371652, 18.612462, 11.1, 1},
			delta:        1e-14,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			pj, err := proj.NewCRSToCRS("EPSG:2180", "EPSG:4326", nil)
			assert.NoError(t, err)
			float64Slice := slices.Clone(tc.float64Slice)
			actual, err := pj.ForwardFloat64Slice(float64Slice)
			assert.NoError(t, err)
			assertInDeltaFloat64Slice(t, tc.expected, actual, tc.delta)
		})
	}
}
