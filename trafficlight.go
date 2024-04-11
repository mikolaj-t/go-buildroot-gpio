package main

import (
	"github.com/warthog618/gpio"
)

type State uint8

const (
	Red State = iota
	Yellow
	Green
)

type TrafficLight struct {
	Green  *gpio.Pin
	Yellow *gpio.Pin
	Red    *gpio.Pin
	state  State
}

func (t *TrafficLight) ChangeState(state State) {
	switch state {
	case Red:
		t.Red.High()
		t.Yellow.Low()
		t.Green.Low()
	case Yellow:
		t.Yellow.High()
		t.Green.Low()
		switch t.state {
		case Red:
			t.Red.High()
		default:
			t.Red.Low()
		}
	case Green:
		t.Red.Low()
		t.Yellow.Low()
		t.Green.High()
	}
	t.state = state
}

func (t *TrafficLight) ChangeToOpposite(state State) {
	switch state {
	case Red:
		t.ChangeState(Green)
	case Green:
		t.ChangeState(Red)
	}
}
