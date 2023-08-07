package game

import (
	api "github.com.bisoncorp.autostrade/gameapi"
	"github.com.bisoncorp.autostrade/graph"
	"github.com.bisoncorp.autostrade/graph/dijkstra"
	"math/rand"
	"sync"
)

type simulation struct {
	speed      float64
	propertyMu sync.RWMutex
	apiMu      sync.Mutex
	cities     []*city
	cityMap    map[string]int
}

func newSimulation(speed float64) *simulation {
	return &simulation{
		speed:   speed,
		cities:  make([]*city, 0),
		cityMap: make(map[string]int),
	}
}

func (s *simulation) Nodes() []graph.Node {
	nodes := make([]graph.Node, len(s.cities))
	for i := range nodes {
		nodes[i] = s.cities[i]
	}
	return nodes
}

func (s *simulation) generateTrip(cityName string, speed float64) []string {
	srcIndex, dstIndex := s.cityMap[cityName], 0
	for {
		dstIndex = rand.Intn(len(s.cities))
		if srcIndex != dstIndex {
			break
		}
	}
	path := dijkstra.ShortestPath(s, srcIndex, dstIndex)[1:]
	pathString := make([]string, len(path))
	for i, index := range path {
		pathString[i] = s.cities[index].Name()
	}
	return pathString
}

func (s *simulation) cityIndex(name string) int {
	return s.cityMap[name]
}

func (s *simulation) AddCity(data api.CityData) api.City {
	s.apiMu.Lock()
	defer s.apiMu.Unlock()
	if _, exist := s.cityMap[data.Name]; exist {
		return nil
	}
	c := newCity(data, func(cityName string, vSpeed float64) []string { return s.generateTrip(cityName, vSpeed) })
	s.cityMap[data.Name] = len(s.cities)
	s.cities = append(s.cities, c)
	return c
}

func (s *simulation) RemoveCity(city api.City) {
	s.apiMu.Lock()
	defer s.apiMu.Unlock()
	index, exist := s.cityMap[city.Name()]
	if !exist {
		return
	}
	city0 := s.cities[index]
	city0.Stop()
	for _, r := range city0.roads {
		r.Stop()
	}
	for _, c := range s.cities {
		if c.Name() == city0.Name() {
			continue
		}
		for _, r := range c.roads {
			if r.dst.Name() == city0.Name() {
				s.RemoveRoad(r)
			}
		}
	}
	for k, v := range s.cityMap {
		if v > index {
			s.cityMap[k] = v - 1
		}
	}
	s.cities = append(s.cities[:index], s.cities[index+1:]...)
}

func (s *simulation) AddRoad(a, b api.City, data api.RoadData) (atob api.Road, btoa api.Road) {
	atob = s.AddOneWayRoad(a, b, data)
	btoa = s.AddOneWayRoad(b, a, data)
	return
}

func (s *simulation) AddOneWayRoad(src, dst api.City, data api.RoadData) api.Road {
	s.apiMu.Lock()
	defer s.apiMu.Unlock()
	src0, ok := src.(*city)
	if !ok {
		return nil
	}
	dst0, ok := dst.(*city)
	if !ok {
		return nil
	}
	for _, r := range src0.roads {
		if r.dst.Name() == dst0.Name() {
			return nil
		}
	}
	r := newRoad(data, func(name string) int { return s.cityIndex(name) }, func() float64 { return s.Speed() }, src0, dst0)

	if src0.Running() {
		src0.Stop()
		defer src0.Start()
	}
	src0.roads = append(src0.roads, r)
	return r
}

func (s *simulation) RemoveRoad(r api.Road) {
	s.apiMu.Lock()
	defer s.apiMu.Unlock()
	r0, ok := r.(*road)
	if !ok {
		return
	}

	src := r0.src
	r0.Stop()
	if src.Running() {
		src.Stop()
		defer src.Start()
	}
	for i, cr := range src.roads {
		if cr.dst.Name() == r0.dst.Name() {
			src.roads = append(src.roads[:i], src.roads[i+1:]...)
			break
		}
	}
}

func (s *simulation) Speed() float64 {
	s.propertyMu.RLock()
	defer s.propertyMu.RUnlock()
	return s.speed
}
func (s *simulation) SetSpeed(speed float64) {
	s.propertyMu.Lock()
	defer s.propertyMu.Unlock()
	s.speed = speed
}

func (s *simulation) PackData() api.SimulationData {
	s.apiMu.Lock()
	defer s.apiMu.Unlock()
	data := api.SimulationData{
		Speed:  s.Speed(),
		Cities: make([]api.CityData, 0, len(s.cities)),
		Roads: make([]struct {
			api.RoadData
			SrcIndex, DstIndex int
		}, 0),
		Vehicles: make([]struct {
			api.VehicleData
			RoadIndex int
		}, 0),
	}

	cityMap := make(map[string]int)
	for _, c := range s.cities {
		c.propertyMu.RLock()
		cd := c.CityData
		c.propertyMu.RUnlock()
		cityMap[cd.Name] = len(data.Cities)
		data.Cities = append(data.Cities, cd)
	}

	for _, c := range s.cities {
		for _, r := range c.roads {
			srcName := c.Name()
			dstName := r.dst.Name()
			r.propertyMu.RLock()
			rd := r.RoadData
			r.propertyMu.RUnlock()
			index := len(data.Roads)
			data.Roads = append(data.Roads, struct {
				api.RoadData
				SrcIndex, DstIndex int
			}{RoadData: rd, SrcIndex: cityMap[srcName], DstIndex: cityMap[dstName]})
			for _, v := range r.vehicles {
				v.propertyMu.RLock()
				vd := v.VehicleData
				v.propertyMu.RUnlock()
				data.Vehicles = append(data.Vehicles, struct {
					api.VehicleData
					RoadIndex int
				}{VehicleData: vd, RoadIndex: index})
			}
		}
	}
	return data
}

func (s *simulation) Start() {
	s.apiMu.Lock()
	defer s.apiMu.Unlock()
	for _, c := range s.cities {
		for _, r := range c.roads {
			r.Start()
		}
		c.Start()
	}
}
func (s *simulation) Stop() {
	s.apiMu.Lock()
	defer s.apiMu.Unlock()
	for _, c := range s.cities {
		c.Stop()
		for _, r := range c.roads {
			r.Stop()
		}
	}
}
func (s *simulation) Running() bool {
	//TODO implement
	panic("not implemented")
}
