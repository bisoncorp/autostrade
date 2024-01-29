package controller

import (
	"github.com/bisoncorp/autostrade/gameapi"
	"sync"
)

type SpeedChangedCallback func(float64)

type SpeedableController struct {
	SpeedableObject gameapi.Speedable
	callbacks       []SpeedChangedCallback
	callbacksMu     sync.Mutex
}

func NewSpeedableController(speedableObject gameapi.Speedable) *SpeedableController {
	return &SpeedableController{SpeedableObject: speedableObject, callbacks: make([]SpeedChangedCallback, 0)}
}

func (s *SpeedableController) AddCallback(fn SpeedChangedCallback) {
	s.callbacksMu.Lock()
	defer s.callbacksMu.Unlock()
	s.callbacks = append(s.callbacks, fn)
}

func (s *SpeedableController) Speed() float64 {
	s.callbacksMu.Lock()
	defer s.callbacksMu.Unlock()
	return s.SpeedableObject.Speed()
}

func (s *SpeedableController) SetSpeed(val float64) {
	s.callbacksMu.Lock()
	defer s.callbacksMu.Unlock()
	s.SpeedableObject.SetSpeed(val)
	s.callAll(val)
}

func (s *SpeedableController) callAll(val float64) {
	for _, fn := range s.callbacks {
		if fn != nil {
			fn(val)
		}
	}
}
