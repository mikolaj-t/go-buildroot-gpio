package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/warthog618/gpio"
)

func main() {
	fmt.Println("Hello embedded world!")
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		os.Exit(0)
	}()

	err := gpio.Open()
	if err != nil {
		panic(err)
	}

	n := gpio.NewPin(NorthGreen)
	n.High()

	s := gpio.NewPin(SouthRed)
	s.High()

	w := gpio.NewPin(WestYellow)
	w.High()

	e := gpio.NewPin(EastRed)
	e.High()

	ee := gpio.NewPin(EastGreen)
	ee.High()

	button := gpio.NewPin(Button)
	button.Input()

	button.Watch(gpio.EdgeRising, func(p *gpio.Pin) {
		n.Low()
		s.High()
		w.Low()
		e.High()
		ee.Low()
	})
}
