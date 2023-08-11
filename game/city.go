package game

import (
	"github.com.bisoncorp.autostrade/game/utils"
	api "github.com.bisoncorp.autostrade/gameapi"
	"github.com.bisoncorp.autostrade/graph"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

type city struct {
	api.CityData
	propertyMu sync.RWMutex

	api.Runnable

	parentSimulation                   *simulation
	entryQueue                         *utils.UnboundedChan[vehicle]
	generationTicker, processingTicker *time.Ticker
	roadsOut, roadsIn                  []*road
	roadsMu                            sync.RWMutex
}

func newCity(data api.CityData, parentSimulation *simulation) *city {
	c := &city{
		CityData:         data,
		parentSimulation: parentSimulation,
		entryQueue:       utils.NewUnboundedChan[vehicle](),
		generationTicker: time.NewTicker(data.GenerationTime),
		processingTicker: time.NewTicker(data.ProcessingTime),
		roadsIn:          make([]*road, 0),
		roadsOut:         make([]*road, 0),
	}
	c.Runnable = utils.NewBaseRunnable(c)
	return c
}

func (c *city) enqueue(v *vehicle) {
	c.entryQueue.In() <- v
}
func (c *city) route(v *vehicle) {
	if v.trip.Arrived() {
		return
	}
	v.trip.Next()
	c.roadsMu.RLock()
	defer c.roadsMu.RUnlock()
	for _, r := range c.roadsOut {
		if r.Dst().Name() == v.trip.Current().Name() {
			r.route(v)
			return
		}
	}
}

func (c *city) addRoadIn(r *road) {
	c.roadsMu.Lock()
	defer c.roadsMu.Unlock()
	for _, r2 := range c.roadsIn {
		if r == r2 {
			return
		}
	}
	c.roadsIn = append(c.roadsIn, r)
}
func (c *city) remRoadIn(r *road) {
	c.roadsMu.Lock()
	defer c.roadsMu.Unlock()
	for index, r2 := range c.roadsIn {
		if r == r2 {
			c.roadsIn = append(c.roadsIn[:index], c.roadsIn[index+1:]...)
			return
		}
	}
}
func (c *city) addRoadOut(r *road) {
	c.roadsMu.Lock()
	defer c.roadsMu.Unlock()
	for _, r2 := range c.roadsOut {
		if r == r2 {
			return
		}
	}
	c.roadsOut = append(c.roadsOut, r)
}
func (c *city) remRoadOut(r *road) {
	c.roadsMu.Lock()
	defer c.roadsMu.Unlock()
	for index, r2 := range c.roadsOut {
		if r == r2 {
			c.roadsOut = append(c.roadsOut[:index], c.roadsOut[index+1:]...)
			return
		}
	}
}
func (c *city) generateVehicle() *vehicle {
	pSpeed := float64(80 + rand.Intn(500))
	v := newVehicle(api.VehicleData{
		Plate:          c.parentSimulation.generatePlate(),
		Color:          colorToRgba(c.Color()),
		PreferredSpeed: pSpeed,
	}, c.parentSimulation.generateTrip(c.Name(), pSpeed))
	return v
}

func (c *city) Update(uint64) {
	select {
	case <-c.generationTicker.C:
		v := c.generateVehicle()
		c.route(v)
	case <-c.processingTicker.C:
		select {
		case v := <-c.entryQueue.Out():
			c.route(v)
		default:
		}
	}
}

func (c *city) Name() string {
	return c.CityData.Name
}
func (c *city) Color() color.Color {
	c.propertyMu.RLock()
	defer c.propertyMu.RUnlock()
	return c.CityData.Color
}
func (c *city) SetColor(col color.Color) {
	c.propertyMu.Lock()
	defer c.propertyMu.Unlock()
	c.CityData.Color = colorToRgba(col)
}
func (c *city) Position() api.Position {
	return c.CityData.Pos
}
func (c *city) SetPosition(position api.Position) {
	c.propertyMu.Lock()
	defer c.propertyMu.Unlock()
	c.CityData.Pos = position
}
func (c *city) GenerationTime() time.Duration {
	c.propertyMu.RLock()
	defer c.propertyMu.RUnlock()
	return c.CityData.GenerationTime
}
func (c *city) SetGenerationTime(duration time.Duration) {
	c.propertyMu.Lock()
	defer c.propertyMu.Unlock()
	c.CityData.GenerationTime = duration
	c.generationTicker.Reset(duration)
}
func (c *city) ProcessingTime() time.Duration {
	c.propertyMu.RLock()
	defer c.propertyMu.RUnlock()
	return c.CityData.ProcessingTime
}
func (c *city) SetProcessingTime(duration time.Duration) {
	c.propertyMu.Lock()
	defer c.propertyMu.Unlock()
	c.CityData.ProcessingTime = duration
	c.processingTicker.Reset(duration)
}
func (c *city) RoadsIn() []api.Road {
	c.roadsMu.RLock()
	defer c.roadsMu.RUnlock()
	rs := make([]api.Road, len(c.roadsIn))
	for i := 0; i < len(c.roadsIn); i++ {
		rs[i] = c.roadsIn[i]
	}
	return rs
}
func (c *city) RoadsOut() []api.Road {
	c.roadsMu.RLock()
	defer c.roadsMu.RUnlock()
	rs := make([]api.Road, len(c.roadsOut))
	for i := 0; i < len(c.roadsOut); i++ {
		rs[i] = c.roadsIn[i]
	}
	return rs
}

func (c *city) Links() []graph.Link {
	links := make([]graph.Link, len(c.roadsOut))
	for i := range links {
		links[i] = c.roadsOut[i]
	}
	return links
}
