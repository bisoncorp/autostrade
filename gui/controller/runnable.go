package controller

import (
	"github.com.bisoncorp.autostrade/gameapi"
	"sync"
)

type RunnableEventType int

const (
	Started RunnableEventType = iota
	Stopped
)

type RunnableCallback func(RunnableEventType)

type RunnableController struct {
	RunnableObject gameapi.Runnable
	callbacks      []RunnableCallback
	callbacksMu    sync.Mutex
}

func NewRunnableController(runnableObject gameapi.Runnable) *RunnableController {
	return &RunnableController{RunnableObject: runnableObject, callbacks: make([]RunnableCallback, 0)}
}

func (r *RunnableController) AddCallback(fn RunnableCallback) {
	r.callbacksMu.Lock()
	defer r.callbacksMu.Unlock()
	r.callbacks = append(r.callbacks, fn)
}

func (r *RunnableController) Start() {
	// Potential Deadlock if the called functions recall start
	r.callbacksMu.Lock()
	defer r.callbacksMu.Unlock()
	if r.RunnableObject == nil {
		return
	}
	r.RunnableObject.Start()
	r.callAll(Started)
}

func (r *RunnableController) Stop() {
	r.callbacksMu.Lock()
	defer r.callbacksMu.Unlock()
	if r.RunnableObject == nil {
		return
	}
	r.RunnableObject.Stop()
	r.callAll(Stopped)
}

func (r *RunnableController) Running() bool {
	r.callbacksMu.Lock()
	defer r.callbacksMu.Unlock()
	if r.RunnableObject != nil {
		return r.RunnableObject.Running()
	}
	return false
}

func (r *RunnableController) callAll(et RunnableEventType) {
	for _, fn := range r.callbacks {
		if fn != nil {
			fn(et)
		}
	}
}
