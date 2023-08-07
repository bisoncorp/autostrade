package gameapi

type RoadData struct {
	MaxSpeed float64
}

type Road interface {
	// MaxSpeed is the maximum speed of all vehicles on this road
	MaxSpeed() float64
	// SetMaxSpeed set maximum speed
	SetMaxSpeed(float64)

	Src() City
	Dst() City

	Vehicles() []Vehicle
	Runnable
}
