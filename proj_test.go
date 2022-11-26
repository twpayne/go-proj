package proj_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"github.com/twpayne/go-proj/v9"
)

var (
	newYorkEPSG3857 = proj.Coord{-8238322.592110482, 4970068.348185822, 10, 0}
	newYorkEPSG4326 = proj.Coord{40.712778, -74.006111, 10, 0}
	parisEPSG3857   = proj.Coord{261848.15527273554, 6250566.54904563, 78, 0}
	parisEPSG4326   = proj.Coord{48.856613, 2.352222, 78, 0}
)

func TestContextNewCRSToCRSTransformation(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	for _, tc := range []struct {
		name        string
		sourceCRS   string
		targetCRS   string
		expectedErr map[int]string
	}{
		{
			name:      "EPSG:4326_to_EPSG;3857",
			sourceCRS: "EPSG:4326",
			targetCRS: "EPSG:3857",
		},
		{
			name:      "EPSG:4326_to_invalid",
			sourceCRS: "EPSG:4326",
			targetCRS: "invalid",
			expectedErr: map[int]string{
				6: "generic error of unknown origin",
				8: "Unknown error (code 4096)",
				9: "Invalid PROJ string syntax",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			transformation, err := context.NewCRSToCRSTransformation(tc.sourceCRS, tc.targetCRS, nil)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr[proj.VersionMajor])
				assert.Nil(t, transformation)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, transformation)
			}
		})
	}
}

func TestTransformationTrans(t *testing.T) {
	for _, tc := range []struct {
		name        string
		sourceCRS   string
		targetCRS   string
		area        *proj.Area
		sourceCoord proj.Coord
		targetCoord proj.Coord
	}{
		{
			name:        "EPSG:4326_to_EPSG:3857_origin",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:3857",
			sourceCoord: proj.Coord{},
			targetCoord: proj.Coord{},
		},
		{
			name:        "EPSG:4326_to_EPSG:3857_origin_with_area",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:3857",
			area:        proj.NewArea(-180, -85, 180, 85),
			sourceCoord: proj.Coord{},
			targetCoord: proj.Coord{},
		},
		{
			name:        "EPSG:4326_to_EPSG:3857_new_york",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:3857",
			sourceCoord: newYorkEPSG4326,
			targetCoord: newYorkEPSG3857,
		},
		{
			name:        "EPSG:4326_to_EPSG:3857_paris",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:3857",
			sourceCoord: parisEPSG4326,
			targetCoord: parisEPSG3857,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC()

			context := proj.NewContext()
			require.NotNil(t, context)

			transformation, err := context.NewCRSToCRSTransformation(tc.sourceCRS, tc.targetCRS, tc.area)
			require.NoError(t, err)
			require.NotNil(t, transformation)

			actualTargetCoord, err := transformation.Forward(tc.sourceCoord)
			require.NoError(t, err)
			assert.InDeltaSlice(t, tc.targetCoord[:], actualTargetCoord[:], 1e-13)

			actualSourceCoord, err := transformation.Trans(proj.DirectionInv, tc.targetCoord)
			require.NoError(t, err)
			assert.InDeltaSlice(t, tc.sourceCoord[:], actualSourceCoord[:], 1e-13)
		})
	}
}

func TestTransformationTransArray(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	transformation, err := context.NewCRSToCRSTransformation("EPSG:4326", "EPSG:3857", nil)
	require.NoError(t, err)
	require.NotNil(t, transformation)

	for _, tc := range []struct {
		name         string
		sourceCoords []proj.Coord
		targetCoords []proj.Coord
	}{
		{
			name: "empty",
		},
		{
			name: "one_element",
			sourceCoords: []proj.Coord{
				newYorkEPSG4326,
			},
			targetCoords: []proj.Coord{
				newYorkEPSG3857,
			},
		},
		{
			name: "two_elements",
			sourceCoords: []proj.Coord{
				newYorkEPSG4326,
				parisEPSG4326,
			},
			targetCoords: []proj.Coord{
				newYorkEPSG3857,
				parisEPSG3857,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, len(tc.targetCoords), len(tc.sourceCoords))

			actualTargetCoords := slices.Clone(tc.sourceCoords)
			require.NoError(t, transformation.ForwardArray(actualTargetCoords))
			for i, actualTargetCoord := range actualTargetCoords {
				assert.InDeltaSlice(t, tc.targetCoords[i][:], actualTargetCoord[:], 1e-13)
			}

			actualSourceCoords := slices.Clone(tc.targetCoords)
			require.NoError(t, transformation.TransArray(proj.DirectionInv, actualSourceCoords))
			for i, actualSourceCoord := range actualSourceCoords {
				assert.InDeltaSlice(t, tc.sourceCoords[i][:], actualSourceCoord[:], 1e-13)
			}
		})
	}
}

func TestTransformationTransError(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	transformation, err := context.NewCRSToCRSTransformation("EPSG:4326", "EPSG:3857", nil)
	require.NoError(t, err)
	require.NotNil(t, transformation)

	for _, tc := range []struct {
		name        string
		direction   proj.Direction
		coord       proj.Coord
		expectedErr map[int]string
	}{
		{
			name:      "invalid_coordinate",
			direction: proj.DirectionFwd,
			coord:     proj.Coord{91, 0, 0, 0},
			expectedErr: map[int]string{
				6: "latitude or longitude exceeded limits",
				8: "Invalid coordinate",
				9: "Invalid coordinate",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actualCoord, err := transformation.Trans(tc.direction, tc.coord)
			assert.EqualError(t, err, tc.expectedErr[proj.VersionMajor])
			assert.Equal(t, proj.Coord{}, actualCoord)

			_, err = transformation.Trans(tc.direction, proj.Coord{})
			require.NoError(t, err)
		})
	}
}

func TestTransformationTransFlatCoords(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	transformation, err := context.NewCRSToCRSTransformation("EPSG:4326", "EPSG:3857", nil)
	require.NoError(t, err)
	require.NotNil(t, transformation)

	for _, tc := range []struct {
		name             string
		sourceFlatCoords []float64
		targetFlatCoords []float64
		stride           int
		zIndex           int
		mIndex           int
	}{
		{
			name: "empty",
		},
		{
			name: "xy",
			sourceFlatCoords: []float64{
				40.712778, -74.006111,
				48.856613, 2.352222,
			},
			targetFlatCoords: []float64{
				-8238322.592110482, 4970068.348185822,
				261848.15527273554, 6250566.54904563,
			},
			stride: 2,
			zIndex: -1,
			mIndex: -1,
		},
		{
			name: "xyz",
			sourceFlatCoords: []float64{
				40.712778, -74.006111, 10,
				48.856613, 2.352222, 78,
			},
			targetFlatCoords: []float64{
				-8238322.592110482, 4970068.348185822, 10,
				261848.15527273554, 6250566.54904563, 78,
			},
			stride: 3,
			zIndex: 2,
			mIndex: -1,
		},
		{
			name: "xym",
			sourceFlatCoords: []float64{
				40.712778, -74.006111, 0,
				48.856613, 2.352222, 1,
			},
			targetFlatCoords: []float64{
				-8238322.592110482, 4970068.348185822, 0,
				261848.15527273554, 6250566.54904563, 1,
			},
			stride: 3,
			zIndex: -1,
			mIndex: 2,
		},
		{
			name: "xyzm",
			sourceFlatCoords: []float64{
				40.712778, -74.006111, 10, 0,
				48.856613, 2.352222, 78, 1,
			},
			targetFlatCoords: []float64{
				-8238322.592110482, 4970068.348185822, 10, 0,
				261848.15527273554, 6250566.54904563, 78, 1,
			},
			stride: 4,
			zIndex: 2,
			mIndex: 3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			actualTargetFlatCoords := slices.Clone(tc.sourceFlatCoords)
			require.NoError(t, transformation.ForwardFlatCoords(actualTargetFlatCoords, tc.stride, tc.zIndex, tc.mIndex))
			assert.InDeltaSlice(t, tc.targetFlatCoords, actualTargetFlatCoords, 1e-9)
		})
	}
}
