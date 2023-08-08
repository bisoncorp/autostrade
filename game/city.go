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

	generateTrip                       func(src string, vSpeed float64) []string
	entryQueue                         *utils.UnboundedChan[vehicle]
	generationTicker, processingTicker *time.Ticker
	roads                              []*road
}

func (c *city) Links() []graph.Link {
	links := make([]graph.Link, len(c.roads))
	for i := range links {
		links[i] = c.roads[i]
	}
	return links
}

func newCity(data api.CityData, generateTrip func(string, float64) []string) *city {
	c := &city{
		CityData:         data,
		generateTrip:     generateTrip,
		entryQueue:       utils.NewUnboundedChan[vehicle](),
		generationTicker: time.NewTicker(data.GenerationTime),
		processingTicker: time.NewTicker(data.ProcessingTime),
	}
	c.Runnable = utils.NewBaseRunnable(c)
	return c
}

func (c *city) enqueue(v *vehicle) {
	c.entryQueue.In() <- v
}

func (c *city) route(v *vehicle) {
	if len(v.trip) == 0 {
		return
	}
	next := v.trip[0]
	v.trip = v.trip[1:]
	for _, r := range c.roads {
		if dst := r.dst; dst.Name() == next {
			r.enqueue(v)
		}
	}
}

func (c *city) Tick() {
	select {
	case <-c.generationTicker.C:
		pSpeed := float64(80 + rand.Intn(500))
		v := newVehicle(api.VehicleData{
			Plate:          <-plateCh,
			Color:          colorToRgba(c.Color()),
			PreferredSpeed: pSpeed,
		}, c.generateTrip(c.Name(), pSpeed))
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
	c.propertyMu.RLock()
	defer c.propertyMu.RUnlock()
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
