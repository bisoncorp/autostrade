package gameapi

import (
	"encoding/json"
	"image/color"
	"math"
	"os"
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

func WriteSimulationData(data SimulationData, path string) {
	file, err := os.Create(path)
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	err = enc.Encode(data)
	if err != nil {
		panic(err)
	}
}

func ReadSimulationData(path string) SimulationData {
	file, err := os.Open(path)
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}
	enc := json.NewDecoder(file)
	data := SimulationData{}
	err = enc.Decode(&data)
	if err != nil {
		panic(err)
	}
	return data
}
