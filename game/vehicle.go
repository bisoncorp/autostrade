package game

import (
	api "github.com.bisoncorp.autostrade/gameapi"
	"image/color"
	"sync"
)

type vehicle struct {
	api.VehicleData
	propertyMu sync.RWMutex

	trip []string
}

func newVehicle(data api.VehicleData, trip []string) *vehicle {
	return &vehicle{VehicleData: data, trip: trip}
}

func (v *vehicle) Plate() string {
	v.propertyMu.RLock()
	defer v.propertyMu.RUnlock()
	return v.VehicleData.Plate
}
func (v *vehicle) Color() color.RGBA {
	v.propertyMu.RLock()
	defer v.propertyMu.RUnlock()
	return v.VehicleData.Color
}
func (v *vehicle) SetColor(rgba color.RGBA) {
	v.propertyMu.Lock()
	defer v.propertyMu.Unlock()
	v.VehicleData.Color = rgba
}
func (v *vehicle) Progress() float64 {
	v.propertyMu.RLock()
	defer v.propertyMu.RUnlock()
	return v.VehicleData.Progress
}
func (v *vehicle) TotalProgress() float64 {
	v.propertyMu.RLock()
	defer v.propertyMu.RUnlock()
	return v.VehicleData.TotalProgress
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
