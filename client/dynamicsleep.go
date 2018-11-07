package client

import (
	"sync/atomic"
	"time"
)

type dynamicSleep struct {
	original         time.Duration
	maxIteration     uint64
	currentIteration uint64
}

func newDynamicSleep(original time.Duration, maxIteration uint64) *dynamicSleep {
	return &dynamicSleep{
		original:     original,
		maxIteration: maxIteration,
	}
}

func (d *dynamicSleep) sleepOriginal() <-chan time.Time {
	return time.After(d.original)
}

func (d *dynamicSleep) sleep() <-chan time.Time {
	i := atomic.LoadUint64(&d.currentIteration)
	return time.After(d.original * time.Duration(i/d.maxIteration))
}

func (d *dynamicSleep) increment() {
	atomic.AddUint64(&d.currentIteration, 1)
}

func (d *dynamicSleep) reset() {
	atomic.StoreUint64(&d.currentIteration, 0)
}
