package proj_test

import (
	"fmt"

	"github.com/twpayne/go-proj/v9"
)

func ExampleTransformation_Forward() {
	transformation, err := proj.NewContext().NewCRSToCRSTransformation("EPSG:4326", "EPSG:3857", nil)
	if err != nil {
		panic(err)
	}

	// Convert Zürich's WGS84 latitude/longitude to Web Mercator.
	zurichEPSG4326 := proj.Coord{47.374444, 8.541111, 408, 0}
	zurichEPSG3857, err := transformation.Forward(zurichEPSG4326)
	if err != nil {
		panic(err)
	}
	fmt.Printf("x=%.6f y=%.6f z=%.6f", zurichEPSG3857[0], zurichEPSG4326[1], zurichEPSG3857[2])
	// Output: x=950792.127329 y=8.541111 z=408.000000
}
