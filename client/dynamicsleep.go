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

func (d *dynamicSleep) sleepOriginal() {
	time.Sleep(d.original)
}

func (d *dynamicSleep) sleepReset() {
	atomic.StoreUint64(&d.currentIteration, 0)
	time.Sleep(d.original)
}

func (d *dynamicSleep) sleepIncrement() {
	i := atomic.LoadUint64(&d.currentIteration)
	atomic.AddUint64(&d.currentIteration, 1)
	time.Sleep(d.original * time.Duration(1+i/d.maxIteration))
}
