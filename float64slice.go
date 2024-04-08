package proj

// ForwardFloat64Slice transforms float64 in place in the forward direction.
func (pj *PJ) ForwardFloat64Slice(float64Slice []float64) ([]float64, error) {
	return pj.TransFloat64Slice(DirectionFwd, float64Slice)
}

// InverseFloat64Slice transforms float64 in place in the forward direction.
func (pj *PJ) InverseFloat64Slice(float64Slice []float64) ([]float64, error) {
	return pj.TransFloat64Slice(DirectionInv, float64Slice)
}

// TransFloat64Slice transforms a []float64 in place.
func (pj *PJ) TransFloat64Slice(direction Direction, float64Slice []float64) ([]float64, error) {
	var coord Coord
	copy(coord[:], float64Slice)
	transCoord, err := pj.Trans(direction, coord)
	if err != nil {
		return nil, err
	}
	copy(float64Slice, transCoord[:])
	return float64Slice, nil
}
