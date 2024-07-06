package proj

// ForwardFloat64Slices transforms float64Slices in the forward direction.
func (pj *PJ) ForwardFloat64Slices(float64Slices [][]float64) error {
	return pj.TransFloat64Slices(DirectionFwd, float64Slices)
}

// InverseFloat64Slices transforms float64Slices in the inverse direction.
func (pj *PJ) InverseFloat64Slices(float64Slices [][]float64) error {
	return pj.TransFloat64Slices(DirectionInv, float64Slices)
}

// TransFloat64Slices transforms float64Slices.
func (pj *PJ) TransFloat64Slices(direction Direction, float64Slices [][]float64) error {
	coords := Float64SlicesToCoords(float64Slices)
	if err := pj.TransArray(direction, coords); err != nil {
		return err
	}
	for i, coord := range coords {
		copy(float64Slices[i], coord[:])
	}
	return nil
}

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
