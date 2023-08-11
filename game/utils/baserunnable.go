package utils

import (
	"github.com.bisoncorp.autostrade/gameapi"
	"sync/atomic"
)

type Updater interface {
	Update(iteration uint64)
}

type baseRunnable struct {
	stopCh  chan struct{}
	running atomic.Bool
	impl    Updater
}

func NewBaseRunnable(impl Updater) gameapi.Runnable {
	return &baseRunnable{stopCh: make(chan struct{}), impl: impl}
}

func (b *baseRunnable) Start() {
	shouldStart := b.running.CompareAndSwap(false, true)
	if !shouldStart {
		return
	}
	go func() {
		iteration := uint64(0)
		for {
			select {
			case <-b.stopCh:
				return
			default:
			}
			b.impl.Update(iteration)
			iteration++
		}
	}()
}

func (b *baseRunnable) Stop() {
	shouldStop := b.running.CompareAndSwap(true, false)
	if shouldStop {
		b.stopCh <- struct{}{}
	}
}

func (b *baseRunnable) Running() bool {
	return b.running.Load()
}
