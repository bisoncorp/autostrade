package utils

import (
	api "github.com.bisoncorp.autostrade/gameapi"
)

type SafeListener struct {
	api.Listener
}

func (s *SafeListener) CityAdded(city api.City) {
	if s.Listener == nil {
		return
	}
	s.Listener.CityAdded(city)
}

func (s *SafeListener) CityRemoved(city api.City) {
	if s.Listener == nil {
		return
	}
	s.Listener.CityRemoved(city)
}

func (s *SafeListener) RoadAdded(road api.Road) {
	if s.Listener == nil {
		return
	}
	s.Listener.RoadAdded(road)
}

func (s *SafeListener) RoadRemoved(road api.Road) {
	if s.Listener == nil {
		return
	}
	s.Listener.RoadRemoved(road)
}

func (s *SafeListener) VehicleSpawned(vehicle api.Vehicle) {
	if s.Listener == nil {
		return
	}
	s.Listener.VehicleSpawned(vehicle)
}

func (s *SafeListener) VehicleDespawned(vehicle api.Vehicle) {
	if s.Listener == nil {
		return
	}
	s.Listener.VehicleDespawned(vehicle)
}
