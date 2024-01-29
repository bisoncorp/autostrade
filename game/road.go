package game

import (
	"github.com/bisoncorp/autostrade/game/utils"
	api "github.com/bisoncorp/autostrade/gameapi"
	"math"
	"sync"
	"time"
)

type road struct {
	api.RoadData
	propertyMu sync.RWMutex

	api.Runnable

	sim    *simulation
	ticker *time.Ticker

	vehicleQueue *utils.UnboundedChan[vehicle]

	vehicles   []*vehicle
	vehiclesMu sync.RWMutex

	src, dst     *city
	lastCallTime time.Time
}

func newRoad(data api.RoadData, sim *simulation, src, dst *city) *road {
	r := &road{
		RoadData:     data,
		sim:          sim,
		ticker:       time.NewTicker(time.Second / 60),
		vehicleQueue: utils.NewUnboundedChan[vehicle](),
		vehicles:     make([]*vehicle, 0, 1<<5),
		dst:          dst,
		src:          src,
	}
	r.Runnable = utils.NewBaseRunnable(r)
	return r
}

func (r *road) Update(iteration uint64) {
	select {
	case <-r.ticker.C:
		if iteration == 0 {
			r.lastCallTime = time.Now()
		}
		timeElapsed := time.Now().Sub(r.lastCallTime)
		r.moveVehicles(timeElapsed.Hours() * r.sim.Speed())
		r.lastCallTime = time.Now()
	case v := <-r.vehicleQueue.Out():
		r.vehiclesMu.Lock()
		r.vehicles = append(r.vehicles, v)
		r.vehiclesMu.Unlock()
	}
}

func (r *road) moveVehicles(timeElapsed float64) {
	maxSpeed := r.MaxSpeed()
	distance := api.Distance(r.src.Position(), r.dst.Position())

	r.vehiclesMu.RLock()
	vehicles := make([]*vehicle, len(r.vehicles))
	copy(vehicles, r.vehicles)
	r.vehiclesMu.RUnlock()

	wg := sync.WaitGroup{}
	wg.Add(len(vehicles))
	for i, v := range vehicles {
		go func(v *vehicle, i int) {
			sr := moveVehicle(v, timeElapsed, maxSpeed, distance)
			if sr {
				r.dst.enqueue(v)
				vehicles[i] = nil
			}
			wg.Done()
		}(v, i)
	}
	wg.Wait()

	nv := make([]*vehicle, 0, len(vehicles))
	for _, v := range vehicles {
		if v != nil {
			nv = append(nv, v)
		}
	}
	r.vehiclesMu.Lock()
	r.vehicles = nv
	r.vehiclesMu.Unlock()
}
func moveVehicle(v *vehicle, timeElapsed float64, maxSpeed float64, distance float64) (shouldRoute bool) {
	v.propertyMu.RLock()
	vd := v.VehicleData
	v.propertyMu.RUnlock()

	speed := math.Min(vd.PreferredSpeed, maxSpeed)
	vd.Progress += speed * timeElapsed / distance
	shouldRoute = vd.Progress >= 1

	v.propertyMu.Lock()
	v.VehicleData.Progress = vd.Progress
	v.propertyMu.Unlock()
	return
}
func (r *road) route(v *vehicle) {
	v.propertyMu.Lock()
	v.VehicleData.Progress = 0
	v.propertyMu.Unlock()
	r.vehicleQueue.In() <- v
}

func (r *road) MaxSpeed() float64 {
	r.propertyMu.RLock()
	defer r.propertyMu.RUnlock()
	return r.RoadData.MaxSpeed
}
func (r *road) SetMaxSpeed(f float64) {
	r.propertyMu.Lock()
	defer r.propertyMu.Unlock()
	r.RoadData.MaxSpeed = f
}
func (r *road) Vehicles() []api.Vehicle {
	r.vehiclesMu.RLock()
	vehicles := make([]*vehicle, len(r.vehicles))
	copy(vehicles, r.vehicles)
	r.vehiclesMu.RUnlock()
	vs := make([]api.Vehicle, len(vehicles))
	for i, v := range vehicles {
		vs[i] = v
	}
	return vs
}
func (r *road) Src() api.City {
	return r.src
}
func (r *road) Dst() api.City {
	return r.dst
}

func (r *road) NodeIndex() int {
	return r.sim.cityIndex(r.dst.Name())
}
func (r *road) Weight() int {
	distance := api.Distance(r.src.Position(), r.dst.Position())
	tripTime := distance / r.MaxSpeed()
	r.vehiclesMu.RLock()
	vNum := float64(len(r.vehicles))
	r.vehiclesMu.RUnlock()
	vehiclePerKm := vNum / (distance)
	tripTime *= vehiclePerKm
	return int(tripTime)
}
