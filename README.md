# go-proj

[![GoDoc](https://godoc.org/github.com/twpayne/go-proj/v9?status.svg)](https://godoc.org/github.com/twpayne/go-proj/v9)

Package go-proj provides an interface to [PROJ](https://proj.org).

## Features

* High performance bulk transformation of coordinates.
* Idiomatic Go API, including complete error handling.
* Supports PROJ versions 6 and upwards.
* Compatible with all geometry libraries.
* Automatically handles C memory management.
* Well tested.

## Example

```go
func ExamplePJ_Forward() {
	pj, err := proj.NewCRSToCRS("EPSG:4326", "EPSG:3857", nil)
	if err != nil {
		panic(err)
	}

	// Convert Zürich's WGS84 latitude/longitude to Web Mercator.
	zurichEPSG4326 := proj.NewCoord(47.374444, 8.541111, 408, 0)
	zurichEPSG3857, err := pj.Forward(zurichEPSG4326)
	if err != nil {
		panic(err)
	}
	fmt.Printf("x=%.6f y=%.6f z=%.6f", zurichEPSG3857.X(), zurichEPSG3857.Y(), zurichEPSG3857.Z())
	// Output: x=950792.127329 y=6003408.475803 z=408.000000
}
```

## Comparisons with other PROJ bindings

There are many existing bindings for PROJ. Generally speaking, these:

* Only transform one coordinate a time, making them extremely slow when
  transforming large number of coordinates.

* Are tied to a single geometry representation.

* Do not handle errors during transformation.

These existing bindings include:

* [`github.com/everystreet/go-proj`](https://github.com/everystreet/go-proj).

* [`github.com/go-spatial/proj`](https://github.com/go-spatial/proj) is an
  incomplete rewrite of PROJ4 in Go with the last commit on October 25, 2019.

* [`github.com/nextgis/go-proj`](https://github.com/nextgis/go-proj).

* [`github.com/omniscale/go-proj`](https://github.com/omniscale/go-proj).

* [`github.com/pebbe/go-proj-4`](https://github.com/pebbe/go-proj-4) has limited
  functionality with the last commit on February 20, 2021.

* [`github.com/xeonx/proj4`](https://github.com/xeonx/proj4) has limited
  functionality with a last commit on December 23, 2015.


## License

MIT

