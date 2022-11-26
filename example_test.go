package proj_test

import (
	"fmt"

	"github.com/twpayne/go-proj/v9"
)

func ExampleTransformation_Forward() {
	transformation, err := proj.NewCRSToCRSTransformation("EPSG:4326", "EPSG:3857", nil)
	if err != nil {
		panic(err)
	}

	// Convert Zürich's WGS84 latitude/longitude to Web Mercator.
	zurichEPSG4326 := proj.NewCoord(47.374444, 8.541111, 408, 0)
	zurichEPSG3857, err := transformation.Forward(zurichEPSG4326)
	if err != nil {
		panic(err)
	}
	fmt.Printf("x=%.6f y=%.6f z=%.6f", zurichEPSG3857.X(), zurichEPSG3857.Y(), zurichEPSG3857.Z())
	// Output: x=950792.127329 y=6003408.475803 z=408.000000
}
