package gameapi

type SimulationData struct {
	Speed  float64
	Cities []CityData
	Roads  []struct {
		RoadData
		SrcIndex, DstIndex int
	}
	Vehicles []struct {
		VehicleData
		RoadIndex int
	}
}

type Simulation interface {
	AddCity(CityData) City
	RemoveCity(City)

	AddRoad(a, b City, data RoadData) (atob Road, btoa Road)
	AddOneWayRoad(src, dst City, data RoadData) Road
	RemoveRoad(Road)

	City(name string) City
	Vehicle(plate string) Vehicle

	PackData() SimulationData

	Speedable
	Runnable
}
