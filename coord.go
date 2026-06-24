package proj

import "math"

type Coord [4]float64

func NewCoord(x, y, z, m float64) Coord {
	return Coord{x, y, z, m}
}

func (c *Coord) DegToRad() Coord {
	return Coord{math.Pi * c[0] / 180, math.Pi * c[1] / 180, c[2], c[3]}
}

func (c *Coord) RadToDeg() Coord {
	return Coord{180 * c[0] / math.Pi, 180 * c[1] / math.Pi, c[2], c[3]}
}

func (c *Coord) X() float64 { return c[0] }

func (c *Coord) Y() float64 { return c[1] }

func (c *Coord) Z() float64 { return c[2] }

func (c *Coord) M() float64 { return c[3] }
