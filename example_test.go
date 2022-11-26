package proj_test

import (
	"fmt"

	"github.com/twpayne/go-proj/v9"
)

func ExampleTransformation_Trans() {
	transformation, err := proj.NewContext().NewCRSToCRSTransformation("EPSG:4326", "EPSG:3857", nil)
	if err != nil {
		panic(err)
	}

	// Convert Zürich's WGS84 latitude/longitude to Web Mercator.
	zurichEPSG4326 := proj.Coord{47.374444, 8.541111, 408, 0}
	zuriceEPSG3857, err := transformation.Trans(proj.DirectionFwd, zurichEPSG4326)
	if err != nil {
		panic(err)
	}
	fmt.Println(zuriceEPSG3857)
	// Output: [950792.1273288276 6.003408475803397e+06 408 0]
}
