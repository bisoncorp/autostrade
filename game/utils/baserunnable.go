package utils

import (
	"github.com.bisoncorp.autostrade/gameapi"
	"sync/atomic"
)

type updater interface {
	update()
}

type baseRunnable struct {
	stopCh  chan struct{}
	running atomic.Bool
	impl    updater
}

func NewBaseRunnable(impl updater) gameapi.Runnable {
	return &baseRunnable{stopCh: make(chan struct{}), impl: impl}
}

func (b *baseRunnable) Start() {
	shouldStart := b.running.CompareAndSwap(false, true)
	if fn := b.impl.update; shouldStart {
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
	if b.Running() {
		b.stopCh <- struct{}{}
		b.running.Store(false)
	}
}

func (b *baseRunnable) Running() bool {
	return b.running.Load()
}
