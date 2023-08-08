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
	callAll(r.callbacks, Started)
}

func (r *RunnableController) Stop() {
	r.callbacksMu.Lock()
	defer r.callbacksMu.Unlock()
	if r.RunnableObject == nil {
		return
	}
	r.RunnableObject.Stop()
	callAll(r.callbacks, Stopped)
}

func (r *RunnableController) Running() bool {
	r.callbacksMu.Lock()
	defer r.callbacksMu.Unlock()
	if r.RunnableObject != nil {
		return r.RunnableObject.Running()
	}
	return false
}

func callAll(fns []RunnableCallback, et RunnableEventType) {
	for _, fn := range fns {
		if fn != nil {
			fn(et)
		}
	}
}
