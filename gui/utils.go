package gui

import (
	"fmt"
	"image/color"
	"math/rand"
)

type speedString float64

func (s speedString) String() string {
	return fmt.Sprint("Simulation speed: ", int(s))
}

func randomColor() color.Color {
	return color.NRGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
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
