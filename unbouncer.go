package main

import (
	"fmt"
	"time"
)

type Unbouncer struct {
	stableOutput chan<- bool
	lastState    *bool
	timer        *time.Timer
	tbmax        time.Duration
}

func NewUnbouncer(stableOutput chan<- bool, tbmax time.Duration) *Unbouncer {
	return &Unbouncer{
		stableOutput: stableOutput,
		tbmax:        tbmax,
	}
}

func (u *Unbouncer) OnClicked(state bool) {
	if u.lastState == nil {
		fmt.Printf("Button pressed first time! %t\n", state)
		u.setState(state)
	} else if state != *u.lastState {
		fmt.Printf("Button pressed new state! %t -> %t\n", *u.lastState, state)
		u.setState(state)
	}
}

func (u *Unbouncer) setState(state bool) {
	u.lastState = &state
	if u.timer != nil {
		u.timer.Stop()
	}
	u.timer = time.AfterFunc(u.tbmax, func() {
		u.stableOutput <- state
		u.lastState = nil
	})
}
