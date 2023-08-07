package game

import api "github.com.bisoncorp.autostrade/gameapi"

func New() api.Simulation {
	return newSimulation(60)
}

func NewFromData(data api.SimulationData) api.Simulation {
	sim := New()
	sim.SetSpeed(data.Speed)
	cityHook := make([]api.City, len(data.Cities))
	for i, cityData := range data.Cities {
		cityHook[i] = sim.AddCity(cityData)
	}
	for _, r := range data.Roads {
		sim.AddOneWayRoad(cityHook[r.SrcIndex], cityHook[r.DstIndex], r.RoadData)
	}
	return sim
}
