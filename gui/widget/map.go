package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	api "github.com.bisoncorp.autostrade/gameapi"
	"image/color"
	"sync"
)

const (
	vehicleDimension = 10
	cityDimension    = 20
	roadDimension    = 2
	MapWidth         = 300
	MapHeight        = 300
)

type Map struct {
	widget.BaseWidget

	OnCityTapped    func(api.City)
	OnVehicleTapped func(api.Vehicle)
	OnTapped        func(fyne.Position)

	MapSize fyne.Size

	cities          map[api.City]*City
	citiesMu        sync.RWMutex
	citiesChanged   bool
	roads           map[api.Road]*Road
	roadsMu         sync.RWMutex
	roadsChanged    bool
	vehicles        map[api.Vehicle]*Vehicle
	vehiclesMu      sync.RWMutex
	vehiclesChanged bool
	simulationHook  api.Simulation
}

func NewMap(sim api.Simulation) *Map {
	m := &Map{
		cities:         make(map[api.City]*City),
		roads:          make(map[api.Road]*Road),
		vehicles:       make(map[api.Vehicle]*Vehicle),
		simulationHook: sim,
	}
	m.ExtendBaseWidget(m)
	cities := sim.Cities()
	for _, city := range cities {
		m.CityAdded(city)
	}
	roads := sim.Roads()
	for _, road := range roads {
		m.RoadAdded(road)
		vehicles := road.Vehicles()
		for _, vehicle := range vehicles {
			m.VehicleSpawned(vehicle)
		}
	}
	return m
}

func (m *Map) ColorRoads(roads []api.Road, col color.Color) {
	m.roadsMu.RLock()
	defer m.roadsMu.RUnlock()
	for _, road := range roads {
		obj, exist := m.roads[road]
		if !exist {
			continue
		}
		obj.line.StrokeColor = col
	}
}

func (m *Map) CityAdded(city api.City) {
	m.citiesMu.Lock()
	m.citiesChanged = true
	m.cities[city] = NewCity(city, m.callOnCityTapped)
	m.citiesMu.Unlock()
	m.Refresh()
}
func (m *Map) CityRemoved(city api.City) {
	m.citiesMu.Lock()
	m.citiesChanged = true
	delete(m.cities, city)
	m.citiesMu.Unlock()
	m.Refresh()
}

func (m *Map) RoadAdded(road api.Road) {
	m.roadsMu.Lock()
	m.roadsChanged = true
	m.roads[road] = NewRoad(road)
	m.roadsMu.Unlock()
	m.Refresh()
}
func (m *Map) RoadRemoved(road api.Road) {
	m.roadsMu.Lock()
	m.roadsChanged = true
	delete(m.roads, road)
	m.roadsMu.Unlock()
	m.Refresh()
}

func (m *Map) VehicleSpawned(vehicle api.Vehicle) {
	m.vehiclesMu.Lock()
	m.vehiclesChanged = true
	m.vehicles[vehicle] = NewVehicle(vehicle, m.callOnVehicleTapped)
	m.vehiclesMu.Unlock()
	m.Refresh()
}
func (m *Map) VehicleDespawned(vehicle api.Vehicle) {
	m.vehiclesMu.Lock()
	m.vehiclesChanged = true
	delete(m.vehicles, vehicle)
	m.vehiclesMu.Unlock()
	m.Refresh()
}

func (m *Map) SetSize(Size fyne.Size) {
	m.MapSize = Size
	m.Refresh()
}
func (m *Map) Tapped(event *fyne.PointEvent) {
	size := m.Size()
	minSize := m.MinSize()
	scaleX, scaleY := size.Width/minSize.Width, size.Height/minSize.Height
	descalePosition := func(position fyne.Position) fyne.Position {
		return fyne.NewPos(position.X/scaleX, position.Y/scaleY)
	}
	m.callOnTapped(descalePosition(event.Position))
}
func (m *Map) MinSize() fyne.Size {
	return m.MapSize
}

func (m *Map) callOnCityTapped(hook api.City) {
	if m.OnCityTapped != nil {
		m.OnCityTapped(hook)
	}
}
func (m *Map) callOnVehicleTapped(hook api.Vehicle) {
	if m.OnVehicleTapped != nil {
		m.OnVehicleTapped(hook)
	}
}
func (m *Map) callOnTapped(position fyne.Position) {
	if m.OnTapped != nil {
		m.OnTapped(position)
	}
}

func (m *Map) CreateRenderer() fyne.WidgetRenderer {
	mr := &mapRenderer{
		wid: m,
	}
	mr.Refresh()
	return mr
}

type mapRenderer struct {
	wid      *Map
	cities   []fyne.CanvasObject
	roads    []fyne.CanvasObject
	vehicles []fyne.CanvasObject
	mutex    sync.RWMutex
}

func (m *mapRenderer) Destroy() {}
func (m *mapRenderer) Layout(size fyne.Size) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	minSize := m.wid.MinSize()
	scaleX, scaleY := size.Width/minSize.Width, size.Height/minSize.Height
	scalePosition := func(position fyne.Position) fyne.Position {
		return fyne.NewPos(position.X*scaleX, position.Y*scaleY)
	}
	_ = func(size fyne.Size) fyne.Size {
		return fyne.NewSize(size.Width*scaleX, size.Height*scaleY)
	}
	citySize := cityMinSize()
	vehicleSize := vehicleMinSize()

	// Cities
	for _, object := range m.cities {
		city := object.(*City)
		hook := city.hook
		city.Resize(citySize)
		city.Move(scalePosition(hook.Position().ToPos32()))
		centerObject(city)
	}

	// Roads
	for _, object := range m.roads {
		road := object.(*Road)
		hook := road.hook
		road.line.Position1 = scalePosition(hook.Src().Position().ToPos32())
		road.line.Position2 = scalePosition(hook.Dst().Position().ToPos32())
	}

	// Vehicles
	for _, object := range m.vehicles {
		if object.(*Vehicle).hook.Road() == nil {
			continue
		}
		vehicle := object.(*Vehicle)
		hook := vehicle.hook
		vehicle.Resize(vehicleSize)
		pos := api.Lerp(hook.Road().Src().Position(), hook.Road().Dst().Position(), hook.Progress()).ToPos32()
		vehicle.Move(scalePosition(pos))
		centerObject(vehicle)
	}

}
func (m *mapRenderer) MinSize() fyne.Size {
	return m.wid.MapSize
}
func (m *mapRenderer) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 0, len(m.roads)+len(m.cities)+len(m.vehicles))
	objs = append(objs, m.roads...)
	objs = append(objs, m.vehicles...)
	objs = append(objs, m.cities...)
	return objs
}
func (m *mapRenderer) Refresh() {
	m.mutex.Lock()

	m.wid.citiesMu.RLock()
	if m.wid.citiesChanged {
		m.cities = make([]fyne.CanvasObject, 0, len(m.wid.cities))
		for _, object := range m.wid.cities {
			m.cities = append(m.cities, object)
		}
		m.wid.citiesChanged = false
	}
	m.wid.citiesMu.RUnlock()

	m.wid.roadsMu.RLock()
	if m.wid.roadsChanged {
		m.roads = make([]fyne.CanvasObject, 0, len(m.wid.roads))
		for _, object := range m.wid.roads {
			m.roads = append(m.roads, object)
		}
		m.wid.roadsChanged = false
	}
	m.wid.roadsMu.RUnlock()

	m.wid.vehiclesMu.RLock()
	if m.wid.vehiclesChanged {
		m.vehicles = make([]fyne.CanvasObject, 0, len(m.wid.vehicles))
		for _, object := range m.wid.vehicles {
			m.vehicles = append(m.vehicles, object)
		}
		m.wid.vehiclesChanged = false
	}
	m.wid.vehiclesMu.RUnlock()

	m.mutex.Unlock()

	for _, city := range m.cities {
		city.Refresh()
	}
	for _, road := range m.roads {
		road.Refresh()
	}
	for _, vehicle := range m.vehicles {
		vehicle.Refresh()
	}

	m.Layout(m.wid.Size())
}

type centerable interface {
	fyne.CanvasObject
	center() fyne.Position
}

func centerObject(object centerable) {
	x, y := object.Position().Components()
	cx, cy := object.center().Components()
	object.Move(fyne.NewPos(x-cx, y-cy))
}
