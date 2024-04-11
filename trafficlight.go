package main

import "github.com/warthog618/gpio"

type TrafficLight struct {
	GreenLight  *gpio.Pin
	YellowLight *gpio.Pin
	RedLight    *gpio.Pin
}
