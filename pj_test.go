package proj_test

import (
	"math"
	"runtime"
	"slices"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

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
	gdanskEPSG4326  = proj.Coord{54.371652, 18.612462, 11.1, 0}
	gdanskEPSG2180  = proj.Coord{723134.1266446244, 474831.4869142064, 11.1, 0}
)

func TestPJ_Info(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	pj, err := context.New("epsg:2056")
	assert.NoError(t, err)
	assert.NotZero(t, pj)

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
	assert.NotZero(t, context)

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
			distDelta:                  1e-9,
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
			assert.NoError(t, err)
			assert.NotZero(t, pj)

			assertInDelta(t, tc.expectedLPDist, pj.LPDist(tc.a, tc.b), tc.distDelta)
			assertInDelta(t, tc.expectedLPDist, pj.LPDist(tc.b, tc.a), tc.distDelta)
			assertInDelta(t, tc.expectedLPZDist, pj.LPZDist(tc.a, tc.b), tc.distDelta)
			assertInDelta(t, tc.expectedLPZDist, pj.LPZDist(tc.b, tc.a), tc.distDelta)

			actualGeodDist, actualGeodForwardAzimuth, actualGeodReverseAzimuth := pj.Geod(tc.a, tc.b)
			assertInDelta(t, tc.expectedGeodDist, actualGeodDist, tc.distDelta)
			assertInDelta(t, tc.expectedGeodForwardAzimuth, actualGeodForwardAzimuth, tc.azimuthDelta)
			assertInDelta(t, tc.expectedGeodReverseAzimuth, actualGeodReverseAzimuth, tc.azimuthDelta)

			actualReverseGeodDist, actualReverseGeodForwardAzimuth, actualReverseGeodReverseAzimuth := pj.Geod(tc.b, tc.a)
			assertInDelta(t, tc.expectedGeodDist, actualReverseGeodDist, tc.distDelta)
			assertInDelta(t, tc.expectedGeodForwardAzimuth, 180+actualReverseGeodReverseAzimuth, tc.azimuthDelta)
			assertInDelta(t, tc.expectedGeodReverseAzimuth, 180+actualReverseGeodForwardAzimuth, tc.azimuthDelta)
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
			assert.NotZero(t, context)

			pj, err := context.NewCRSToCRS(tc.sourceCRS, tc.targetCRS, tc.area)
			assert.NoError(t, err)
			assert.NotZero(t, pj)

			actualTargetCoord, err := pj.Forward(tc.sourceCoord)
			assert.NoError(t, err)
			assertInDeltaFloat64Slice(t, tc.targetCoord[:], actualTargetCoord[:], tc.targetDelta)

			actualSourceCoord, err := pj.Inverse(tc.targetCoord)
			assert.NoError(t, err)
			assertInDeltaFloat64Slice(t, tc.sourceCoord[:], actualSourceCoord[:], tc.sourceDelta)
		})
	}
}

func TestPJ_Trans_error(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	assert.NoError(t, err)
	assert.NotZero(t, pj)

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
			assert.NoError(t, err)
		})
	}
}

func TestPJ_TransArray(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	assert.NoError(t, err)
	assert.NotZero(t, pj)

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
			assert.Equal(t, len(tc.targetCoords), len(tc.sourceCoords))

			actualTargetCoords := slices.Clone(tc.sourceCoords)
			assert.NoError(t, pj.ForwardArray(actualTargetCoords))
			for i, actualTargetCoord := range actualTargetCoords {
				assertInDeltaFloat64Slice(t, tc.targetCoords[i][:], actualTargetCoord[:], 1e1)
			}

			actualSourceCoords := slices.Clone(tc.targetCoords)
			assert.NoError(t, pj.InverseArray(actualSourceCoords))
			for i, actualSourceCoord := range actualSourceCoords {
				assertInDeltaFloat64Slice(t, tc.sourceCoords[i][:], actualSourceCoord[:], 1e-13)
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
	assert.NotZero(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:2056", nil)
	assert.NoError(t, err)
	assert.NotZero(t, pj)

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
			assertInDelta(t, tc.targetBounds.XMin, targetBounds.XMin, tc.targetDelta)
			assertInDelta(t, tc.targetBounds.YMin, targetBounds.YMin, tc.targetDelta)
			assertInDelta(t, tc.targetBounds.XMax, targetBounds.XMax, tc.targetDelta)
			assertInDelta(t, tc.targetBounds.YMax, targetBounds.YMax, tc.targetDelta)

			sourceBounds, err := pj.InverseBounds(tc.targetBounds, 21)
			assert.NoError(t, err)
			assertInDelta(t, tc.sourceBounds.XMin, sourceBounds.XMin, tc.sourceDelta)
			assertInDelta(t, tc.sourceBounds.YMin, sourceBounds.YMin, tc.sourceDelta)
			assertInDelta(t, tc.sourceBounds.XMax, sourceBounds.XMax, tc.sourceDelta)
			assertInDelta(t, tc.sourceBounds.YMax, sourceBounds.YMax, tc.sourceDelta)
		})
	}
}

func TestPJ_TransFlatCoords(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	assert.NoError(t, err)
	assert.NotZero(t, pj)

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
			assert.NoError(t, pj.ForwardFlatCoords(actualTargetFlatCoords, tc.stride, tc.zIndex, tc.mIndex))
			assertInDeltaFloat64Slice(t, tc.targetFlatCoords, actualTargetFlatCoords, 1e1)

			actualSourceFlatCoords := slices.Clone(tc.targetFlatCoords)
			assert.NoError(t, pj.InverseFlatCoords(actualSourceFlatCoords, tc.stride, tc.zIndex, tc.mIndex))
			assertInDeltaFloat64Slice(t, tc.sourceFlatCoords, actualSourceFlatCoords, 1e-9)
		})
	}
}

func TestPJ_TransFloat64Slice(t *testing.T) {
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

func TestPJ_NormalizeForVisualizationForNorthingEastingCRS(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:2180", nil)
	assert.NoError(t, err)
	assert.NotZero(t, pj)

	t.Run("original axis order", func(t *testing.T) {
		actualCoord, err := pj.Forward(gdanskEPSG4326)
		// Original axis order. X is northing, Y is easting.
		assert.NoError(t, err)
		assertInDeltaFloat64Slice(t, gdanskEPSG2180[:], actualCoord[:], 1e-7)
	})

	t.Run("normalized axis order", func(t *testing.T) {
		// Create a new PJ with the axis swap operation.
		pj, err = pj.NormalizeForVisualization()
		assert.NoError(t, err)

		// The axis order of geographic CRS is now longitude, latitude.
		swappedGdanskEPSG4326 := proj.Coord{gdanskEPSG4326[1], gdanskEPSG4326[0], gdanskEPSG4326[2], gdanskEPSG4326[3]}

		actualCoord, err := pj.Forward(swappedGdanskEPSG4326)

		// Normalized axis order. X is easting, Y is northing.
		swappedGdanskEPSG2180 := proj.Coord{gdanskEPSG2180[1], gdanskEPSG2180[0], gdanskEPSG2180[2], gdanskEPSG2180[3]}
		assert.NoError(t, err)
		assertInDeltaFloat64Slice(t, swappedGdanskEPSG2180[:], actualCoord[:], 1e-7)
	})
}

func TestPJ_NormalizeForVisualizationForEastingNorthingCRS(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	pj, err := context.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	assert.NoError(t, err)
	assert.NotZero(t, pj)

	// Create a new PJ with the axis swap operation.
	pj, err = pj.NormalizeForVisualization()
	assert.NoError(t, err)

	// The axis order of geographic CRS is now longitude, latitude.
	swappedNewYorkEPSG4326 := proj.Coord{newYorkEPSG4326[1], newYorkEPSG4326[0], newYorkEPSG4326[2], newYorkEPSG4326[3]}

	actualCoord, err := pj.Forward(swappedNewYorkEPSG4326)

	// The output axis order is not changed.
	assert.NoError(t, err)
	assertInDeltaFloat64Slice(t, newYorkEPSG3857[:], actualCoord[:], 1e-7)
}

func TestPJ_TransFloat64Slices(t *testing.T) {
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
			delta: 1e-9,
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
			delta: 1e-9,
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
			delta: 1e-9,
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

func assertInDelta(tb testing.TB, expected, actual, delta float64) {
	tb.Helper()
	if actualDelta := math.Abs(expected - actual); actualDelta > delta {
		tb.Fatalf("Expected %e to be within %e of %e, but delta is %e", actual, delta, expected, actualDelta)
	}
}

func assertInDeltaFloat64Slice(tb testing.TB, expected, actual []float64, delta float64) {
	tb.Helper()
	assert.Equal(tb, len(expected), len(actual))
	for i := range expected {
		if actualDelta := math.Abs(expected[i] - actual[i]); actualDelta > delta {
			tb.Fatalf("Expected %e to be within %e of %e at index %d, but delta is %e", actual[i], delta, expected[i], i, actualDelta)
		}
	}
}

func assertInDeltaFloat64Slices(tb testing.TB, expected, actual [][]float64, delta float64) {
	tb.Helper()
	assert.Equal(tb, len(expected), len(actual))
	for i := range expected {
		assertInDeltaFloat64Slice(tb, expected[i], actual[i], delta)
	}
}
