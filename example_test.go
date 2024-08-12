package proj_test

import (
	"fmt"

	"github.com/twpayne/go-proj/v10"
)

func ExamplePJ_Forward() {
	pj, err := proj.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	if err != nil {
		panic(err)
	}
	defer pj.Destroy()

	// Start with Zürich's WGS84 latitude/longitude.
	zurich4326 := proj.NewCoord(47.374444, 8.541111, 408, 0)
	fmt.Printf("initial: x=%.6f y=%.6f z=%.6f\n", zurich4326.X(), zurich4326.Y(), zurich4326.Z())

	// Convert Zürich's WGS84 latitude/longitude to Web Mercator.
	zurich3857, err := pj.Forward(zurich4326)
	if err != nil {
		panic(err)
	}
	fmt.Printf("forward: x=%.6f y=%.6f z=%.6f\n", zurich3857.X(), zurich3857.Y(), zurich3857.Z())

	// ...and convert back.
	zurich4326After, err := pj.Inverse(zurich3857)
	if err != nil {
		panic(err)
	}
	fmt.Printf("inverse: x=%.6f y=%.6f z=%.6f", zurich4326After.X(), zurich4326After.Y(), zurich4326After.Z())

	// Output:
	// initial: x=47.374444 y=8.541111 z=408.000000
	// forward: x=950792.127329 y=6003408.475803 z=408.000000
	// inverse: x=47.374444 y=8.541111 z=408.000000
}
