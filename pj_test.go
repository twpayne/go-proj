package proj_test

import (
	"math"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"

	"github.com/twpayne/go-proj/v10"
)

var (
	bernEPSG4326    = proj.Coord{46.948056, 7.4475, 540, 0}
	bernEPSG2056    = proj.Coord{2600675.0876650945, 1199663.542715189, 540, 0}
	zurichEPSG4326  = proj.Coord{47.374444, 8.541111, 408, 0}
	zurichEPSG2056  = proj.Coord{2683263.251826082, 1247651.9664695852, 408, 0}
	newYorkEPSG3857 = proj.Coord{-8238322.592110482, 4970068.348185822, 10, 0}
	newYorkEPSG4326 = proj.Coord{40.712778, -74.006111, 10, 0}
	parisEPSG3857   = proj.Coord{261848.15527273554, 6250566.54904563, 78, 0}
	parisEPSG4326   = proj.Coord{48.856613, 2.352222, 78, 0}
)

func TestPJ_Info(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	pj, err := context.New("epsg:2056")
	require.NoError(t, err)
	require.NotNil(t, pj)

	expectedInfo := proj.PJInfo{
		Description: "CH1903+ / LV95",
		Accuracy:    -1,
	}
	actualInfo := pj.Info()
	assert.Equal(t, expectedInfo, actualInfo)
}

func TestPJ_LPDist(t *testing.T) {
	if proj.VersionMajor < 7 {
		t.Skip("distance functions not tested")
	}

	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	for i, tc := range []struct {
		definition                 string
		a                          proj.Coord
		b                          proj.Coord
		expectedLPDist             float64
		expectedLPZDist            float64
		expectedGeodDist           float64
		expectedGeodForwardAzimuth float64
		expectedGeodReverseAzimuth float64
		distDelta                  float64
		azimuthDelta               float64
	}{
		{
			definition:                 "epsg:4326",
			a:                          bernEPSG4326.DegToRad(),
			b:                          zurichEPSG4326.DegToRad(),
			expectedLPDist:             129762.08359988699,
			expectedLPZDist:            129762.15073812571,
			expectedGeodDist:           129762.08359988699,
			expectedGeodForwardAzimuth: 21.20947946541022,
			expectedGeodReverseAzimuth: 21.268782222540885,
			distDelta:                  math.SmallestNonzeroFloat64,
			azimuthDelta:               1e-13,
		},
		{
			definition:                 "epsg:4326",
			a:                          newYorkEPSG4326.DegToRad(),
			b:                          parisEPSG4326.DegToRad(),
			expectedLPDist:             8494402.471778858,
			expectedLPZDist:            8494402.472051037,
			expectedGeodDist:           8494402.471778858,
			expectedGeodForwardAzimuth: 8.381709060115105,
			expectedGeodReverseAzimuth: 2.310935629050629,
			distDelta:                  math.SmallestNonzeroFloat64,
			azimuthDelta:               1e-13,
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			pj, err := context.New(tc.definition)
			require.NoError(t, err)
			require.NotNil(t, pj)

			assert.InDelta(t, tc.expectedLPDist, pj.LPDist(tc.a, tc.b), tc.distDelta)
			assert.InDelta(t, tc.expectedLPDist, pj.LPDist(tc.b, tc.a), tc.distDelta)
			assert.InDelta(t, tc.expectedLPZDist, pj.LPZDist(tc.a, tc.b), tc.distDelta)
			assert.InDelta(t, tc.expectedLPZDist, pj.LPZDist(tc.b, tc.a), tc.distDelta)

			actualGeodDist, actualGeodForwardAzimuth, actualGeodReverseAzimuth := pj.Geod(tc.a, tc.b)
			assert.Equal(t, tc.expectedGeodDist, actualGeodDist)
			assert.InDelta(t, tc.expectedGeodForwardAzimuth, actualGeodForwardAzimuth, tc.azimuthDelta)
			assert.InDelta(t, tc.expectedGeodReverseAzimuth, actualGeodReverseAzimuth, tc.azimuthDelta)

			actualReverseGeodDist, actualReverseGeodForwardAzimuth, actualReverseGeodReverseAzimuth := pj.Geod(tc.b, tc.a)
			assert.Equal(t, tc.expectedGeodDist, actualReverseGeodDist)
			assert.InDelta(t, tc.expectedGeodForwardAzimuth, 180+actualReverseGeodReverseAzimuth, tc.azimuthDelta)
			assert.InDelta(t, tc.expectedGeodReverseAzimuth, 180+actualReverseGeodForwardAzimuth, tc.azimuthDelta)
		})
	}
}

func TestPJ_Trans(t *testing.T) {
	for _, tc := range []struct {
		name        string
		sourceCRS   string
		targetCRS   string
		area        *proj.Area
		sourceCoord proj.Coord
		targetCoord proj.Coord
		sourceDelta float64
		targetDelta float64
	}{
		{
			name:        "EPSG:4326_to_EPSG:3857_origin",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:3857",
			sourceCoord: proj.Coord{},
			targetCoord: proj.Coord{},
			sourceDelta: 1e-13,
			targetDelta: 1e1,
		},
		{
			name:        "EPSG:4326_to_EPSG:3857_origin_with_area",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:3857",
			area:        proj.NewArea(-180, -85, 180, 85),
			sourceCoord: proj.Coord{},
			targetCoord: proj.Coord{},
			sourceDelta: 1e-13,
			targetDelta: 1e1,
		},
		{
			name:        "EPSG:4326_to_EPSG:2056_bern",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:2056",
			sourceCoord: bernEPSG4326,
			targetCoord: bernEPSG2056,
			sourceDelta: 1e-6,
			targetDelta: 1e1,
		},
		{
			name:        "EPSG:4326_to_EPSG:2056_zurich",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:2056",
			sourceCoord: zurichEPSG4326,
			targetCoord: zurichEPSG2056,
			sourceDelta: 1e-6,
			targetDelta: 1e1,
		},
		{
			name:        "EPSG:4326_to_EPSG:3857_new_york",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:3857",
			sourceCoord: newYorkEPSG4326,
			targetCoord: newYorkEPSG3857,
			sourceDelta: 1e-13,
			targetDelta: 1e1,
		},
		{
			name:        "EPSG:4326_to_EPSG:3857_paris",
			sourceCRS:   "EPSG:4326",
			targetCRS:   "EPSG:3857",
			sourceCoord: parisEPSG4326,
			targetCoord: parisEPSG3857,
			sourceDelta: 1e-13,
			targetDelta: 1e1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer runtime.GC()

			context := proj.NewContext()
			require.NotNil(t, context)

			pj, err := context.NewCRSToCRS(tc.sourceCRS, tc.targetCRS, tc.area)
			require.NoError(t, err)
			require.NotNil(t, pj)

			actualTargetCoord, err := pj.Forward(tc.sourceCoord)
			require.NoError(t, err)
			assert.InDeltaSlice(t, tc.targetCoord[:], actualTargetCoord[:], tc.targetDelta)

			actualSourceCoord, err := pj.Inverse(tc.targetCoord)
			require.NoError(t, err)
			assert.InDeltaSlice(t, tc.sourceCoord[:], actualSourceCoord[:], tc.sourceDelta)
		})
	}
}

func TestPJ_Trans_error(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	require.NoError(t, err)
	require.NotNil(t, pj)

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
			actualCoord, err := pj.Trans(tc.direction, tc.coord)
			assert.EqualError(t, err, tc.expectedErr[proj.VersionMajor])
			assert.Equal(t, proj.Coord{}, actualCoord)

			_, err = pj.Trans(tc.direction, proj.Coord{})
			require.NoError(t, err)
		})
	}
}

func TestPJ_TransArray(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	require.NoError(t, err)
	require.NotNil(t, pj)

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
			require.NoError(t, pj.ForwardArray(actualTargetCoords))
			for i, actualTargetCoord := range actualTargetCoords {
				assert.InDeltaSlice(t, tc.targetCoords[i][:], actualTargetCoord[:], 1e1)
			}

			actualSourceCoords := slices.Clone(tc.targetCoords)
			require.NoError(t, pj.InverseArray(actualSourceCoords))
			for i, actualSourceCoord := range actualSourceCoords {
				assert.InDeltaSlice(t, tc.sourceCoords[i][:], actualSourceCoord[:], 1e-13)
			}
		})
	}
}

func TestPJ_TransBounds(t *testing.T) {
	if proj.VersionMajor < 8 || proj.VersionMajor == 8 && proj.VersionMinor < 2 {
		t.Skip()
	}

	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:2056", nil)
	require.NoError(t, err)
	require.NotNil(t, pj)

	for _, tc := range []struct {
		name         string
		sourceBounds proj.Bounds
		targetBounds proj.Bounds
		sourceDelta  float64
		targetDelta  float64
	}{
		{
			name: "EPSG:4326_to_EPSG:2056",
			sourceBounds: proj.Bounds{
				XMin: bernEPSG4326.X(),
				YMin: bernEPSG4326.Y(),
				XMax: zurichEPSG4326.X(),
				YMax: zurichEPSG4326.Y(),
			},
			targetBounds: proj.Bounds{
				XMin: bernEPSG2056.X(),
				YMin: bernEPSG2056.Y(),
				XMax: zurichEPSG2056.X(),
				YMax: zurichEPSG2056.Y(),
			},
			sourceDelta: 1e-2,
			targetDelta: 1e3,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			targetBounds, err := pj.ForwardBounds(tc.sourceBounds, 21)
			assert.NoError(t, err)
			assert.InDelta(t, tc.targetBounds.XMin, targetBounds.XMin, tc.targetDelta)
			assert.InDelta(t, tc.targetBounds.YMin, targetBounds.YMin, tc.targetDelta)
			assert.InDelta(t, tc.targetBounds.XMax, targetBounds.XMax, tc.targetDelta)
			assert.InDelta(t, tc.targetBounds.YMax, targetBounds.YMax, tc.targetDelta)

			sourceBounds, err := pj.InverseBounds(tc.targetBounds, 21)
			assert.NoError(t, err)
			assert.InDelta(t, tc.sourceBounds.XMin, sourceBounds.XMin, tc.sourceDelta)
			assert.InDelta(t, tc.sourceBounds.YMin, sourceBounds.YMin, tc.sourceDelta)
			assert.InDelta(t, tc.sourceBounds.XMax, sourceBounds.XMax, tc.sourceDelta)
			assert.InDelta(t, tc.sourceBounds.YMax, sourceBounds.YMax, tc.sourceDelta)
		})
	}
}

func TestPJ_TransFlatCoords(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	require.NoError(t, err)
	require.NotNil(t, pj)

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
			require.NoError(t, pj.ForwardFlatCoords(actualTargetFlatCoords, tc.stride, tc.zIndex, tc.mIndex))
			assert.InDeltaSlice(t, tc.targetFlatCoords, actualTargetFlatCoords, 1e1)

			actualSourceFlatCoords := slices.Clone(tc.targetFlatCoords)
			require.NoError(t, pj.InverseFlatCoords(actualSourceFlatCoords, tc.stride, tc.zIndex, tc.mIndex))
			assert.InDeltaSlice(t, tc.sourceFlatCoords, actualSourceFlatCoords, 1e-9)
		})
	}
}
