# go-proj

[![GoDoc](https://pkg.go.dev/badge/github.com/twpayne/go-proj/v11)](https://pkg.go.dev/github.com/twpayne/go-proj/v11)

Package go-proj provides an interface to [PROJ](https://proj.org).

## Features

* High performance bulk transformation of coordinates.
* Idiomatic Go API, including complete error handling.
* Supports PROJ versions 6 and upwards.
* Compatible with all geometry libraries.
* Convenience functions for handling coordinates as `[]float64`s.
* Automatically handles C memory management.
* Well tested.

## Install

```console
$ go get github.com/twpayne/go-proj/v11
```

You must also install the PROJ development headers and libraries. These are
typically in the package `libproj-dev` on Debian-like systems, `proj-devel` on
RedHat-like systems, and `proj` in Homebrew.

## Example

```go
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
```

## Comparisons with other PROJ bindings

There are many existing bindings for PROJ. Generally speaking, these:

* Only transform one coordinate a time, making them extremely slow when
  transforming large number of coordinates.

* Are tied to a single geometry representation.

* Do not handle errors during transformation.

* Are no longer maintained.

These existing bindings include:

* [`github.com/everystreet/go-proj`](https://github.com/everystreet/go-proj),
  latest commit December 11, 2021.

* [`github.com/go-spatial/proj`](https://github.com/go-spatial/proj) is an
  incomplete rewrite of PROJ4 in Go with the last commit on October 25, 2019.

* [`github.com/omniscale/go-proj`](https://github.com/omniscale/go-proj) has
  limited functionality.

* [`github.com/pebbe/go-proj-4`](https://github.com/pebbe/go-proj-4) has limited
  functionality with the last commit on February 20, 2021.

* [`github.com/xeonx/proj4`](https://github.com/xeonx/proj4) has limited
  functionality with a last commit on December 23, 2015.


## License

MIT
