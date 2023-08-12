package gameapi

type SimulationData struct {
	Speed     float64
	LastPlate Plate
	Cities    []CityData
	Roads     []struct {
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
	Road(a, b string) (atob, btoa Road)
	Vehicle(plate string) Vehicle

	Cities() []City
	Roads() []Road

	PackData() SimulationData

	Speedable
	Runnable
}
