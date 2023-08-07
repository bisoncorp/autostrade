package utils

import (
	"github.com.bisoncorp.autostrade/gameapi"
	"sync/atomic"
)

type Ticker interface {
	Tick()
}

type baseRunnable struct {
	stopCh  chan struct{}
	running atomic.Bool
	impl    Ticker
}

func NewBaseRunnable(impl Ticker) gameapi.Runnable {
	return &baseRunnable{stopCh: make(chan struct{}), impl: impl}
}

func (b *baseRunnable) Start() {
	shouldStart := b.running.CompareAndSwap(false, true)
	if fn := b.impl.Tick; shouldStart {
		go func() {
			for {
				select {
				case <-b.stopCh:
					return
				default:
				}
				fn()
			}
		}()
	}
}

func (b *baseRunnable) Stop() {
	b.stopCh <- struct{}{}
	b.running.Store(false)
}

func (b *baseRunnable) Running() bool {
	return b.running.Load()
}
