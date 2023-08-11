package gameapi

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
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
	Color() color.Color
	SetColor(color.Color)
}

type Speedable interface {
	Speed() float64
	SetSpeed(float64)
}

var FirstPlate = Plate{
	A: 'A',
	B: 'A',
	C: 'A',
	D: 'A',
	N: 0,
}

type Plate struct {
	A, B, C, D rune
	N          int
}

func (p Plate) String() string {
	return fmt.Sprintf("%c%c%03d%c%c", p.A, p.B, p.N, p.C, p.D)
}

func (p Plate) Next() Plate {
	incrementRune := func(r rune) (rune, bool) {
		r++
		if r == 'Z'+1 {
			return 'A', true
		}
		return r, false
	}
	incrementNext := false
	p.N++
	if p.N == 1000 {
		p.N = 0
		incrementNext = true
	}
	if incrementNext {
		p.D, incrementNext = incrementRune(p.D)
	}
	if incrementNext {
		p.C, incrementNext = incrementRune(p.C)
	}
	if incrementNext {
		p.B, incrementNext = incrementRune(p.B)
	}
	if incrementNext {
		p.A, _ = incrementRune(p.A)
	}
	return p
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

func WriteSimulationData(data SimulationData, writer io.Writer) {
	enc := json.NewEncoder(writer)
	enc.SetIndent("", "\t")
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
}

func ReadSimulationData(reader io.Reader) SimulationData {
	enc := json.NewDecoder(reader)
	data := SimulationData{}
	err := enc.Decode(&data)
	if err != nil {
		panic(err)
	}
	return data
}
