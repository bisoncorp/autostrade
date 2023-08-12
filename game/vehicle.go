package game

import (
	api "github.com.bisoncorp.autostrade/gameapi"
	"image/color"
	"sync"
)

type vehicle struct {
	api.VehicleData
	propertyMu sync.RWMutex

	trip        api.Trip
	currentRoad *road
}

func newVehicle(data api.VehicleData, trip api.Trip) *vehicle {
	return &vehicle{VehicleData: data, trip: trip}
}

func (v *vehicle) Plate() string {
	v.propertyMu.RLock()
	defer v.propertyMu.RUnlock()
	return v.VehicleData.Plate
}
func (v *vehicle) Color() color.Color {
	v.propertyMu.RLock()
	defer v.propertyMu.RUnlock()
	return v.VehicleData.Color
}
func (v *vehicle) SetColor(c color.Color) {
	v.propertyMu.Lock()
	defer v.propertyMu.Unlock()
	v.VehicleData.Color = colorToRgba(c)
}
func (v *vehicle) Progress() float64 {
	v.propertyMu.RLock()
	defer v.propertyMu.RUnlock()
	return v.VehicleData.Progress
}
func (v *vehicle) PreferredSpeed() float64 {
	v.propertyMu.RLock()
	defer v.propertyMu.RUnlock()
	return v.VehicleData.PreferredSpeed
}
func (v *vehicle) SetPreferredSpeed(f float64) {
	v.propertyMu.Lock()
	defer v.propertyMu.Unlock()
	v.VehicleData.PreferredSpeed = f
}
func (v *vehicle) Trip() api.Trip {
	return v.trip
}
func (v *vehicle) Road() api.Road {
	if v.currentRoad == nil {
		return nil
	}
	return v.currentRoad
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
