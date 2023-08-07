package gameapi

import "image/color"

type VehicleData struct {
	Plate          string
	Color          color.RGBA
	Progress       float64
	TotalProgress  float64
	PreferredSpeed float64
}

type Vehicle interface {
	// Plate is the license plate of vehicle, it is unique
	Plate() string

	Colorable

	// Progress is the progress on the current road. Interval [0, 1)
	Progress() float64
	// TotalProgress is the progress from Src() to Dst(). Interval [0, 1)
	TotalProgress() float64

	// PreferredSpeed of the Vehicle
	PreferredSpeed() float64
	// SetPreferredSpeed set speed of the vehicle, the speed is capped to Road().MaxSpeed()
	SetPreferredSpeed(float64)
}
