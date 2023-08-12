package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	api "github.com.bisoncorp.autostrade/gameapi"
	"github.com.bisoncorp.autostrade/set"
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
	OnTapped        func(*fyne.PointEvent)

	MapSize fyne.Size

	SimulationHook api.Simulation
}

func NewMap() *Map {
	m := &Map{}
	m.ExtendBaseWidget(m)
	return m
}

func (m *Map) SetSize(Size fyne.Size) {
	m.MapSize = Size
	m.Refresh()
}

func (m *Map) Tapped(event *fyne.PointEvent) {
	m.callOnTapped(event)
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
func (m *Map) callOnTapped(event *fyne.PointEvent) {
	if m.OnTapped != nil {
		m.OnTapped(event)
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
	scaleSize := func(size fyne.Size) fyne.Size {
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
		road.Resize(scaleSize(road.MinSize()))
		road.Move(hook.Src().Position().ToPos32())
	}

	// Vehicles
	for _, object := range m.vehicles {
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
	needHover := set.New[string]()
	for _, object := range m.cities {
		if city := object.(*City); city.hover {
			needHover.Insert(city.hook.Name())
		}
	}
	for _, object := range m.vehicles {
		if vehicle := object.(*Vehicle); vehicle.hover {
			needHover.Insert(vehicle.hook.Plate())
		}
	}

	if m.wid.SimulationHook == nil {
		m.mutex.Unlock()
		return
	}

	citiesHooks := m.wid.SimulationHook.Cities()
	m.cities = make([]fyne.CanvasObject, len(citiesHooks))
	for i := 0; i < len(citiesHooks); i++ {
		city := NewCity(citiesHooks[i], m.wid.callOnCityTapped)
		if needHover.Has(citiesHooks[i].Name()) {
			city.hover = true
			city.Refresh()
		}
		m.cities[i] = city
	}

	roadsHooks := m.wid.SimulationHook.Roads()
	m.roads = make([]fyne.CanvasObject, len(roadsHooks))
	for i := 0; i < len(roadsHooks); i++ {
		m.roads[i] = NewRoad(roadsHooks[i])
	}

	vehiclesHooks := make([]api.Vehicle, 0, 4000)
	for _, hook := range roadsHooks {
		vehiclesHooks = append(vehiclesHooks, hook.Vehicles()...)
	}
	m.vehicles = make([]fyne.CanvasObject, len(vehiclesHooks))
	for i := 0; i < len(vehiclesHooks); i++ {
		vehicle := NewVehicle(vehiclesHooks[i], m.wid.callOnVehicleTapped)
		if needHover.Has(vehiclesHooks[i].Plate()) {
			vehicle.hover = true
			vehicle.Refresh()
		}
		m.vehicles[i] = vehicle
	}
	m.mutex.Unlock()
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
