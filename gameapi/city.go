package gameapi

import (
	"image/color"
	"time"
)

type CityData struct {
	Name           string
	Color          color.RGBA
	Pos            Position
	GenerationTime time.Duration
	ProcessingTime time.Duration
}

type City interface {
	// Name is the city name
	Name() string

	Colorable

	// Position of the city
	Position() Position

	// GenerationTime is the time for city to generate a vehicle
	GenerationTime() time.Duration
	// SetGenerationTime set generation time
	SetGenerationTime(time.Duration)

	// ProcessingTime is time for the city to poll element from the queue
	ProcessingTime() time.Duration
	// SetProcessingTime set consume time
	SetProcessingTime(time.Duration)

	Runnable
}
