package utm_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-proj/v11/utm"
)

func TestCoord_String(t *testing.T) {
	for _, tc := range []struct {
		name     string
		coord    utm.Coord
		expected string
	}{
		{
			name: "CN Tower",
			coord: utm.Coord{
				Zone: 17,
				E:    630_084,
				N:    4_833_438,
			},
			expected: "17N 630084 4833438",
		},
		{
			name: "Sydney",
			coord: utm.Coord{
				Zone: 56,
				E:    333_504,
				N:    6_251_170 - 10_000_000,
			},
			expected: "56S 333504 6251170",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.coord.String())
		})
	}
}

func TestForward(t *testing.T) {
	for _, tc := range []struct {
		name     string
		lon, lat float64
		expected utm.Coord
	}{
		{
			name: "CN Tower",
			lon:  -degrees(79, 23, 13.7),
			lat:  degrees(43, 38, 33.24),
			expected: utm.Coord{
				Zone: 17,
				E:    630_084,
				N:    4_833_438,
			},
		},
		{
			name: "Toronto",
			lon:  -degrees(79, 22, 48),
			lat:  degrees(43, 39, 0),
			expected: utm.Coord{
				Zone: 17,
				E:    630_644,
				N:    4_834_275,
			},
		},
		{
			name: "Sydney",
			lon:  degrees(151, 12, 0),
			lat:  -degrees(33, 52, 0),
			expected: utm.Coord{
				Zone: 56,
				E:    333_504,
				N:    6_251_170 - 10_000_000,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := utm.Forward(tc.lon, tc.lat)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected.Zone, actual.Zone)
			assertWithin(t, 1, tc.expected.E, actual.E)
			assertWithin(t, 1, tc.expected.N, actual.N)
		})
	}
}

func TestZone(t *testing.T) {
	for _, tc := range []struct {
		lon, lat float64
		expected int
	}{
		{
			lon:      0,
			lat:      0,
			expected: 31,
		},
		{
			lon:      -118.246765,
			lat:      34.030811,
			expected: 11,
		},
		{
			lon:      -43.199158,
			lat:      -22.911851,
			expected: 23,
		},
		{
			lon:      -0.111694,
			lat:      51.489593,
			expected: 30,
		},
		{
			lon:      77.208252,
			lat:      28.598027,
			expected: 43,
		},
		{
			lon:      5.727882,
			lat:      58.961764,
			expected: 32,
		},
		{
			lon:      151.192017,
			lat:      -33.895953,
			expected: 56,
		},
	} {
		t.Run(fmt.Sprintf("%f_%f", tc.lon, tc.lat), func(t *testing.T) {
			assert.Equal(t, tc.expected, utm.Zone(tc.lon, tc.lat))
		})
	}
}

func assertWithin(tb testing.TB, maxDelta, actual, expected float64) {
	tb.Helper()
	delta := math.Abs(expected - actual)
	if delta <= maxDelta {
		return
	}
	tb.Fatalf("Expected %v to be within %v of %v, but delta is %v", actual, maxDelta, expected, delta)
}

func degrees(degrees, minutes int, seconds float64) float64 {
	return float64(degrees) + float64(minutes)/60 + seconds/3600
}
