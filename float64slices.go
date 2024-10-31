package proj

// CoordsToFloat64Slices is a convenience function that converts a slice of
// Coords to a slice of []float64s. For performance, the returned []float64s
// alias coords.
func CoordsToFloat64Slices(coords []Coord) [][]float64 {
	if coords == nil {
		return nil
	}
	float64Slices := make([][]float64, len(coords))
	for i := range float64Slices {
		float64Slices[i] = coords[i][:]
	}
	return float64Slices
}

// Float64Slices is a convenience function that converts a slice of []float64s
// to a slice of Coords.
func Float64SlicesToCoords(float64Slices [][]float64) []Coord {
	if float64Slices == nil {
		return nil
	}
	coords := make([]Coord, len(float64Slices))
	for i := range coords {
		var coord Coord
		copy(coord[:], float64Slices[i])
		coords[i] = coord
	}
	return coords
}
