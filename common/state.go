package common

import "sync/atomic"

const (
	stateOpen = iota
	stateClosed
)

type State struct {
	value uint32
}

func (s *State) IsOpen() bool {
	var state uint32
	atomic.LoadUint32(&state)
	return state == stateOpen
}

func (s *State) IsClosed() bool {
	var state uint32
	atomic.LoadUint32(&state)
	return state == stateClosed
}

func (s *State) SetClosed() {
	atomic.StoreUint32(&s.value, stateClosed)
}
