package proj_test

import (
	"runtime"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-proj/v11"
)

func TestContext_NewCRSToCRS(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

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
			pj, err := context.NewCRSToCRS(tc.sourceCRS, tc.targetCRS, nil)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr[proj.VersionMajor])
				assert.Zero(t, pj)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, pj)
			}
		})
	}
}

func TestContext_NewCRSToCRSFromPJ(t *testing.T) {
	defer runtime.GC()

	sourceCRS, err := proj.New("epsg:4326")
	assert.NoError(t, err)
	assert.True(t, sourceCRS.IsCRS())

	targetCRS, err := proj.New("epsg:3857")
	assert.NoError(t, err)
	assert.True(t, targetCRS.IsCRS())

	pj, err := proj.NewCRSToCRSFromPJ(sourceCRS, targetCRS, nil, "")
	assert.NoError(t, err)
	assert.NotZero(t, pj)
}

func TestContext_New(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	for _, tc := range []struct {
		definition  string
		expectedErr map[int]string
	}{
		{
			definition: "epsg:4326",
		},
		{
			definition: "+proj=etmerc +lat_0=38 +lon_0=125 +ellps=bessel",
		},
		{
			definition: "invalid",
			expectedErr: map[int]string{
				6: "generic error of unknown origin",
				8: "Unknown error (code 4096)",
				9: "Invalid PROJ string syntax",
			},
		},
	} {
		t.Run(tc.definition, func(t *testing.T) {
			pj, err := context.New(tc.definition)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr[proj.VersionMajor])
				assert.Zero(t, pj)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, pj)
			}
		})
	}
}

func TestContext_NewFromArgs(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	for i, tc := range []struct {
		args        []string
		expectedErr map[int]string
	}{
		{
			args: []string{"proj=utm", "zone=32", "ellps=GRS80"},
		},
		{
			args: []string{"proj=utm", "zone=0", "ellps=GRS80"},
			expectedErr: map[int]string{
				6: "invalid UTM zone number",
				8: "Invalid value for an argument",
				9: "Invalid value for an argument",
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			pj, err := context.NewFromArgs(tc.args...)
			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr[proj.VersionMajor])
				assert.Zero(t, pj)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, pj)
			}
		})
	}
}

func TestContext_SetSearchPaths(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	assert.NotZero(t, context)

	// The C function does not return any error so we only validate
	// that executing the SetSearchPaths function call
	// does not panic considering various boundary conditions
	context.SetSearchPaths(nil)
	context.SetSearchPaths([]string{})
	context.SetSearchPaths([]string{"/tmp/data"})
	context.SetSearchPaths([]string{"/tmp/data", "/tmp/data2"})
}
