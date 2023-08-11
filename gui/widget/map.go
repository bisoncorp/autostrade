package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	api "github.com.bisoncorp.autostrade/gameapi"
	"sync"
	"time"
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

	OnCityTapped    func(data api.CityData)
	OnVehicleTapped func(data api.VehicleData)
	OnDataRequired  func() api.SimulationData
	OnTapped        func(event *fyne.PointEvent)

	simulationTicker *time.Ticker

	data   api.SimulationData
	dataMu sync.RWMutex // protection for multithreading refresh operation
}

func (m *Map) SetData(data api.SimulationData) {
	m.dataMu.Lock()
	m.data = data
	m.dataMu.Unlock()
	m.Refresh()
}
func (m *Map) Start() {
	m.simulationTicker.Reset(time.Second / 60)
}
func (m *Map) Stop() {
	m.simulationTicker.Stop()
}

func (m *Map) run() {
	go func() {
		for {
			<-m.simulationTicker.C
			if m.OnDataRequired != nil {
				m.SetData(m.OnDataRequired())
			}
		}
	}()
}
func (m *Map) callOnCityTapped(data api.CityData) {
	if m.OnCityTapped != nil {
		m.OnCityTapped(data)
	}
}
func (m *Map) callOnVehicleTapped(data api.VehicleData) {
	if m.OnVehicleTapped != nil {
		m.OnVehicleTapped(data)
	}
}

func (m *Map) Tapped(event *fyne.PointEvent) {
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

func NewMap() *Map {
	m := &Map{
		simulationTicker: time.NewTicker(time.Second / 60),
	}
	m.ExtendBaseWidget(m)
	m.run()
	return m
}

type mapRenderer struct {
	wid      *Map
	cities   []fyne.CanvasObject
	roads    []fyne.CanvasObject
	vehicles []fyne.CanvasObject
}

func (m *mapRenderer) Destroy() {}
func (m *mapRenderer) Layout(size fyne.Size) {
	m.wid.dataMu.Lock()
	defer m.wid.dataMu.Unlock()
	data := m.wid.data
	m.cities = refreshCityObjects(m.cities, data.Cities, func(data api.CityData) {
		m.wid.callOnCityTapped(data)
	})
	m.roads = refreshRoadObjects(m.roads, data.Roads, data.Cities)
	m.vehicles = refreshVehicleObjects(m.vehicles, data.Vehicles, func(data api.VehicleData) {
		m.wid.callOnVehicleTapped(data)
	})
	for i, city := range m.cities {
		city.Resize(city.MinSize())
		city.Move(scale(data.Cities[i].Pos.ToPos32(), scaleFactor))
		centerObject(city)
	}
	for i, vehicle := range m.vehicles {
		roadIndex := data.Vehicles[i].RoadIndex
		road := data.Roads[roadIndex]
		srcPos := data.Cities[road.SrcIndex].Pos
		dstPos := data.Cities[road.DstIndex].Pos
		progress := data.Vehicles[i].Progress
		actualPos := fyne.Position(api.Lerp(srcPos, dstPos, progress).ToPos32())
		vehicle.Resize(vehicle.MinSize())
		vehicle.Move(scale(actualPos, scaleFactor))
		centerObject(vehicle)
	}
}
func (m *mapRenderer) MinSize() fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, city := range m.cities {
		ms := city.MinSize()
		w, h := ms.Components()
		x, y := city.Position().Components()
		ms = fyne.NewSize(w+x, h+y)
		minSize = minSize.Max(ms)
	}
	return minSize
}
func (m *mapRenderer) Objects() []fyne.CanvasObject {
	objs := make([]fyne.CanvasObject, 0, len(m.roads)+len(m.cities)+len(m.vehicles))
	objs = append(objs, m.roads...)
	objs = append(objs, m.vehicles...)
	objs = append(objs, m.cities...)
	return objs
}
func (m *mapRenderer) Refresh() {
	m.Layout(fyne.Size{})
	canvas.Refresh(m.wid)
}

func refreshCityObjects(
	drawableObjects []fyne.CanvasObject,
	citiesData []api.CityData,
	onTapped func(data api.CityData),
) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(citiesData))
	for i := 0; i < len(objects); i++ {
		if i < len(drawableObjects) {
			objects[i] = drawableObjects[i]
		} else {
			objects[i] = NewCity(citiesData[i], onTapped)
		}
		objects[i].(*City).SetData(citiesData[i])
	}
	return objects
}
func refreshRoadObjects(
	drawableObjects []fyne.CanvasObject,
	roadsData []struct {
		api.RoadData
		SrcIndex, DstIndex int
	},
	citiesData []api.CityData,
) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(roadsData))
	for i := 0; i < len(objects); i++ {
		if i < len(drawableObjects) {
			objects[i] = drawableObjects[i]
		} else {
			objects[i] = NewRoad()
		}
		objects[i].(*Road).SetData(roadsData[i].RoadData, citiesData[roadsData[i].SrcIndex], citiesData[roadsData[i].DstIndex])
	}
	return objects
}
func refreshVehicleObjects(
	vehicles []fyne.CanvasObject,
	vehicleData []struct {
		api.VehicleData
		RoadIndex int
	},
	onTapped func(api.VehicleData),
) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(vehicleData))
	for i := 0; i < len(vehicleData); i++ {
		if i < len(vehicles) {
			objects[i] = vehicles[i]
		} else {
			objects[i] = NewVehicle(onTapped)
		}
		objects[i].(*Vehicle).SetData(vehicleData[i].VehicleData)
	}
	return objects
}

func centerObject(object fyne.CanvasObject) {
	x, y := object.Position().Components()
	w, h := object.Size().Components()
	object.Move(fyne.NewPos(x-w/2, y-h/2))
}

const scaleFactor = 1

func scale(v struct{ X, Y float32 }, scale float32) struct{ X, Y float32 } {
	s := scale
	return struct{ X, Y float32 }{X: v.X * s, Y: v.Y * s}
}
func calcScaleFactor(size, normalSize fyne.Size) float32 {
	sw := size.Width / normalSize.Width
	sh := size.Height / normalSize.Height
	if sw > sh {
		return sw
	}
	return sh
}
