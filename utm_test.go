package proj

import (
	"testing"
)

func TestUTM(t *testing.T) {
	for _, c := range []struct {
		φ, λ   float64
		zone   int
		letter string
	}{
		{φ: -80, λ: -180, zone: 1, letter: "C"},
		{φ: -73, λ: -175, zone: 1, letter: "C"},
		{φ: -66, λ: -170, zone: 2, letter: "D"},
		{φ: -59, λ: -165, zone: 3, letter: "E"},
		{φ: -52, λ: -160, zone: 4, letter: "F"},
		{φ: -45, λ: -155, zone: 5, letter: "G"},
		{φ: -38, λ: -150, zone: 6, letter: "H"},
		{φ: -31, λ: -145, zone: 6, letter: "J"},
		{φ: -24, λ: -140, zone: 7, letter: "K"},
		{φ: -17, λ: -135, zone: 8, letter: "K"},
		{φ: -10, λ: -130, zone: 9, letter: "L"},
		{φ: -3, λ: -125, zone: 10, letter: "M"},
		{φ: 4, λ: -120, zone: 11, letter: "N"},
		{φ: 11, λ: -115, zone: 11, letter: "P"},
		{φ: 18, λ: -110, zone: 12, letter: "Q"},
		{φ: 25, λ: -105, zone: 13, letter: "R"},
		{φ: 32, λ: -100, zone: 14, letter: "S"},
		{φ: 39, λ: -95, zone: 15, letter: "S"},
		{φ: 46, λ: -90, zone: 16, letter: "T"},
		{φ: 53, λ: -85, zone: 16, letter: "U"},
		{φ: 60, λ: -80, zone: 17, letter: "V"},
		{φ: 67, λ: -75, zone: 18, letter: "W"},
		{φ: 74, λ: -70, zone: 19, letter: "X"},
		{φ: 81, λ: -65, zone: 20, letter: "X"},

		{φ: 57, λ: 1, zone: 31, letter: "V"},
		{φ: 57, λ: 2, zone: 31, letter: "V"},
		{φ: 63, λ: 1, zone: 31, letter: "V"},
		{φ: 63, λ: 2, zone: 31, letter: "V"},

		{φ: 57, λ: 4, zone: 32, letter: "V"},
		{φ: 57, λ: 11, zone: 32, letter: "V"},
		{φ: 63, λ: 4, zone: 32, letter: "V"},
		{φ: 63, λ: 11, zone: 32, letter: "V"},

		{φ: 73, λ: 1, zone: 31, letter: "X"},
		{φ: 73, λ: 8, zone: 31, letter: "X"},
		{φ: 83, λ: 1, zone: 31, letter: "X"},
		{φ: 83, λ: 8, zone: 31, letter: "X"},

		{φ: 73, λ: 10, zone: 33, letter: "X"},
		{φ: 73, λ: 20, zone: 33, letter: "X"},
		{φ: 83, λ: 10, zone: 33, letter: "X"},
		{φ: 83, λ: 20, zone: 33, letter: "X"},

		{φ: 73, λ: 22, zone: 35, letter: "X"},
		{φ: 73, λ: 32, zone: 35, letter: "X"},
		{φ: 83, λ: 22, zone: 35, letter: "X"},
		{φ: 83, λ: 32, zone: 35, letter: "X"},

		{φ: 73, λ: 34, zone: 37, letter: "X"},
		{φ: 73, λ: 41, zone: 37, letter: "X"},
		{φ: 83, λ: 34, zone: 37, letter: "X"},
		{φ: 83, λ: 41, zone: 37, letter: "X"},
	} {
		if zone := LatLongUTMZone(c.φ, c.λ); zone != c.zone {
			t.Errorf("LatLongUTMZone(%v, %v) == %v, want %v", c.φ, c.λ, zone, c.zone)
		}
		if letter := LatLongUTMLetter(c.φ, c.λ); letter != c.letter {
			t.Errorf("LatLongUTMLetter(%v, %v) == %v, want %v", c.φ, c.λ, letter, c.letter)
		}
	}
}
