# go-proj

[![GoDoc](https://godoc.org/github.com/twpayne/go-proj?status.svg)](https://godoc.org/github.com/twpayne/go-proj)

Package go-proj provides an interface to [PROJ](https://proj.org).

## Features

* High performance bulk transformation of coordinates.
* Idiomatic Go API, including complete error handling.
* Compatible with all geometry libraries.
* Well tested.

## Example

```go
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
```

## Comparisons with related software

There are many existing bindings for PROJ. Generally speaking, these:

* Only transform one coordinate a time, making them extremely slow when
  transforming large number of coordinates.

* Are tied to a single geometry representation.

* Do not handle errors during transformation.

* Use an older version of the PROJ library.

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

