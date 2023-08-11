package game

import api "github.com.bisoncorp.autostrade/gameapi"

func New() api.Simulation {
	return NewFromData(api.SimulationData{LastPlate: api.FirstPlate})
}

func NewFromData(data api.SimulationData) api.Simulation {
	sim := newSimulation(data)
	return sim
}
