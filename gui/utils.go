package gui

import (
	"fmt"
	"image/color"
)

type speedString float64

func (s speedString) String() string {
	return fmt.Sprint("Simulation speed: ", int(s))
}

func colorToRgba(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}
