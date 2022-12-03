package proj_test

import (
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twpayne/go-proj/v9"
)

func TestContext_NewCRSToCRSTransformation(t *testing.T) {
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

func TestContext_NewTransformation(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

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
			transformation, err := context.NewTransformation(tc.definition)
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

func TestContext_NewTransformationFromArgs(t *testing.T) {
	defer runtime.GC()

	context := proj.NewContext()
	require.NotNil(t, context)

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
			transformation, err := context.NewTransformationFromArgs(tc.args...)
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
