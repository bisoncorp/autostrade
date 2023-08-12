package game

import (
	"fmt"
	"github.com.bisoncorp.autostrade/game/utils"
	api "github.com.bisoncorp.autostrade/gameapi"
	"github.com.bisoncorp.autostrade/graph"
	"github.com.bisoncorp.autostrade/graph/dijkstra"
	"math/rand"
	"sync"
	"sync/atomic"
)

type simulation struct {
	speed      float64
	propertyMu sync.RWMutex

	lastPlateCh <-chan api.Plate
	nextPlateCh <-chan api.Plate

	cities   []*city
	cityMap  map[string]int
	citiesMu sync.RWMutex

	roads   []*road
	roadMap map[string]int
	roadsMu sync.RWMutex

	listener *utils.SafeListener

	running atomic.Bool
}

func newSimulation(data api.SimulationData) *simulation {
	s := &simulation{
		speed:    data.Speed,
		cities:   make([]*city, 0),
		cityMap:  make(map[string]int),
		roads:    make([]*road, 0),
		roadMap:  make(map[string]int),
		listener: &utils.SafeListener{},
	}

	cityHook := make([]api.City, len(data.Cities))
	for i := 0; i < len(data.Cities); i++ {
		cityHook[i] = s.AddCity(data.Cities[i])
	}
	for i := 0; i < len(data.Roads); i++ {
		rd := data.Roads[i]
		s.AddOneWayRoad(cityHook[rd.SrcIndex], cityHook[rd.DstIndex], rd.RoadData)
	}
	s.startPlateGenerator(data.LastPlate)
	return s
}

func (s *simulation) cityIndex(name string) int {
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	return s.cityMap[name]
}
func (s *simulation) generateTrip(src string, _ float64) api.Trip {
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	srcIndex, dstIndex := s.cityIndex(src), -1
	maxIndex := len(s.cities)
	for {
		dstIndex = rand.Intn(maxIndex)
		if srcIndex != dstIndex {
			break
		}
	}
	s.roadsMu.RLock()
	path := dijkstra.ShortestPath(s, srcIndex, dstIndex)
	s.roadsMu.RUnlock()
	cities := make([]api.City, len(path))
	for i := 0; i < len(path); i++ {
		cities[i] = s.cities[path[i]]
	}
	roads := make([]api.Road, len(cities)-1)
	for i := 0; i < len(cities)-1; i++ {
		roads[i], _ = s.Road(cities[i].Name(), cities[i+1].Name())
	}
	return api.NewTrip(cities, roads)
}
func (s *simulation) generatePlate() string {
	return (<-s.nextPlateCh).String()
}
func (s *simulation) startPlateGenerator(initialPlate api.Plate) {
	nextPlateCh, lastPlateCh := make(chan api.Plate, 16), make(chan api.Plate)
	s.nextPlateCh, s.lastPlateCh = nextPlateCh, lastPlateCh
	go func() {
		plate := initialPlate
		for {
			select {
			case nextPlateCh <- plate:
				plate = plate.Next()
			case lastPlateCh <- plate:
			}
		}
	}()
}

func (s *simulation) AddCity(data api.CityData) api.City {
	s.citiesMu.Lock()
	defer s.citiesMu.Unlock()

	if _, exist := s.cityMap[data.Name]; exist {
		return nil
	}

	c := newCity(data, s)
	s.cityMap[data.Name] = len(s.cities)
	s.cities = append(s.cities, c)
	s.listener.CityAdded(c)
	return c
}
func (s *simulation) RemoveCity(c api.City) {
	s.citiesMu.Lock()
	defer s.citiesMu.Unlock()

	index, exist := s.cityMap[c.Name()]
	if !exist {
		return
	}
	city0 := s.cities[index]

	roadsIn := city0.RoadsIn()
	roadsOut := city0.RoadsOut()
	roads := append(roadsIn, roadsOut...)
	city0.Stop()
	for _, r := range roads {
		s.RemoveRoad(r)
	}

	for k, v := range s.cityMap {
		if v > index {
			s.cityMap[k] = v - 1
		}
	}
	s.cities = append(s.cities[:index], s.cities[index+1:]...)
	s.listener.CityRemoved(c)
}

func (s *simulation) AddRoad(a, b api.City, data api.RoadData) (atob api.Road, btoa api.Road) {
	atob = s.AddOneWayRoad(a, b, data)
	btoa = s.AddOneWayRoad(b, a, data)
	return
}
func (s *simulation) AddOneWayRoad(src, dst api.City, data api.RoadData) api.Road {
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	indexSrc, existSrc := s.cityMap[src.Name()]
	indexDst, existDst := s.cityMap[dst.Name()]
	if !existSrc || !existDst {
		return nil
	}
	src0, dst0 := s.cities[indexSrc], s.cities[indexDst]

	s.roadsMu.Lock()
	defer s.roadsMu.Unlock()

	name := roadName(src0.Name(), dst0.Name())
	_, existRoad := s.roadMap[name]
	if existRoad {
		return nil
	}

	r := newRoad(data, s, src0, dst0)
	src0.addRoadOut(r)
	dst0.addRoadIn(r)

	s.roadMap[name] = len(s.roads)
	s.roads = append(s.roads, r)
	s.listener.RoadAdded(r)
	return r
}
func (s *simulation) RemoveRoad(r api.Road) {
	s.roadsMu.Lock()
	defer s.roadsMu.Unlock()
	name := roadName(r.Src().Name(), r.Dst().Name())
	index, existRoad := s.roadMap[name]
	if !existRoad {
		return
	}

	r0 := s.roads[index]
	r0.Stop()
	r0.Src().(*city).remRoadOut(r0)
	r0.Dst().(*city).remRoadIn(r0)

	for k, v := range s.roadMap {
		if v > index {
			s.roadMap[k] = v - 1
		}
	}
	s.roads = append(s.roads[:index], s.roads[index+1:]...)
	s.listener.RoadRemoved(r)
}

func (s *simulation) City(name string) api.City {
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	if i, exist := s.cityMap[name]; exist {
		return s.cities[i]
	}
	return nil
}
func (s *simulation) Road(a, b string) (atob, btoa api.Road) {
	s.roadsMu.RLock()
	defer s.roadsMu.RUnlock()
	atobString, btoaString := roadName(a, b), roadName(b, a)
	if index, exist := s.roadMap[atobString]; exist {
		atob = s.roads[index]
	}
	if index, exist := s.roadMap[btoaString]; exist {
		btoa = s.roads[index]
	}
	return
}
func (s *simulation) Vehicle(plate string) api.Vehicle {
	s.roadsMu.RLock()
	defer s.roadsMu.RUnlock()
	for _, r := range s.roads {
		vehicles := r.Vehicles()
		for _, v := range vehicles {
			if v.Plate() == plate {
				return v
			}
		}
	}
	return nil
}

func (s *simulation) Cities() []api.City {
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	cities := make([]api.City, len(s.cities))
	for i := 0; i < len(s.cities); i++ {
		cities[i] = s.cities[i]
	}
	return cities
}
func (s *simulation) Roads() []api.Road {
	s.roadsMu.RLock()
	defer s.roadsMu.RUnlock()
	roads := make([]api.Road, len(s.roads))
	for i := 0; i < len(s.roads); i++ {
		roads[i] = s.roads[i]
	}
	return roads
}

func (s *simulation) PackData() api.SimulationData {
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	s.roadsMu.RLock()
	defer s.roadsMu.RUnlock()

	data := api.SimulationData{
		Speed:     s.Speed(),
		LastPlate: <-s.lastPlateCh,
		Cities:    make([]api.CityData, 0, len(s.cities)),
		Roads: make([]struct {
			api.RoadData
			SrcIndex, DstIndex int
		}, 0),
		Vehicles: make([]struct {
			api.VehicleData
			RoadIndex int
		}, 0),
	}

	for _, c := range s.cities {
		c.propertyMu.RLock()
		cd := c.CityData
		c.propertyMu.RUnlock()
		data.Cities = append(data.Cities, cd)
	}

	for _, r := range s.roads {
		srcName, dstName := r.Src().Name(), r.Dst().Name()
		srcIndex, dstIndex := s.cityMap[srcName], s.cityMap[dstName]
		r.propertyMu.RLock()
		rd := r.RoadData
		r.propertyMu.RUnlock()
		roadIndex := s.roadMap[roadName(srcName, dstName)]
		data.Roads = append(data.Roads, struct {
			api.RoadData
			SrcIndex, DstIndex int
		}{RoadData: rd, SrcIndex: srcIndex, DstIndex: dstIndex})
		r.vehiclesMu.RLock()
		vehicles := make([]*vehicle, len(r.vehicles))
		copy(vehicles, r.vehicles)
		r.vehiclesMu.RUnlock()
		for _, v := range vehicles {
			v.propertyMu.RLock()
			vd := v.VehicleData
			v.propertyMu.RUnlock()
			data.Vehicles = append(data.Vehicles, struct {
				api.VehicleData
				RoadIndex int
			}{VehicleData: vd, RoadIndex: roadIndex})
		}
	}

	return data
}

func (s *simulation) SetListener(listener api.Listener) {
	s.listener.Listener = listener
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

func (s *simulation) Start() {
	shouldStart := s.running.CompareAndSwap(false, true)
	if !shouldStart {
		return
	}
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	s.roadsMu.RLock()
	defer s.roadsMu.RUnlock()
	for _, c := range s.cities {
		c.Start()
	}
	for _, r := range s.roads {
		r.Start()
	}
}
func (s *simulation) Stop() {
	shouldStop := s.running.CompareAndSwap(true, false)
	if !shouldStop {
		return
	}
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	s.roadsMu.RLock()
	defer s.roadsMu.RUnlock()
	for _, c := range s.cities {
		c.Stop()
	}
	for _, r := range s.roads {
		r.Stop()
	}
}
func (s *simulation) Running() bool {
	return s.running.Load()
}

func (s *simulation) Nodes() []graph.Node {
	s.citiesMu.RLock()
	defer s.citiesMu.RUnlock()
	nodes := make([]graph.Node, len(s.cities))
	for i := range nodes {
		nodes[i] = s.cities[i]
	}
	return nodes
}

func roadName(a, b string) string {
	return fmt.Sprintf("%s-%s", a, b)
}
