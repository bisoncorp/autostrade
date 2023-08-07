package gameapi

import (
	"image/color"
	"math"
)

type Runnable interface {
	// Start controls when run process start
	Start()
	// Stop controls when run process stop
	Stop()
	// Running returns current state
	Running() bool
}

type Colorable interface {
	Color() color.RGBA
	SetColor(color.RGBA)
}

type Position struct {
	X, Y float64
}

func (p Position) ToPos32() struct{ X, Y float32 } {
	return struct{ X, Y float32 }{
		X: float32(p.X),
		Y: float32(p.Y),
	}
}

// Lerp linear interpolation
func Lerp(p1, p2 Position, l float64) Position {
	return Position{
		X: p1.X + (p2.X-p1.X)*l,
		Y: p1.Y + (p2.Y-p1.Y)*l,
	}
}

func Distance(p1, p2 Position) float64 {
	dx, dy := p2.X-p1.X, p2.Y-p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}
