package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	fmt.Println("Hello embedded world!")
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		rpio.Close()
		os.Exit(0)
	}()

	err := rpio.Open()
	if err != nil {
		panic(err)
	}

	n := rpio.Pin(NorthGreen)
	n.High()

	s := rpio.Pin(SouthRed)
	s.High()

	w := rpio.Pin(WestYellow)
	w.High()

	e := rpio.Pin(EastRed)
	e.High()

	ee := rpio.Pin(EastGreen)
	ee.High()

	button := rpio.Pin(Button)
	button.Input()

	buzzer := rpio.Pin(Buzzer)

	for {
		buzzer.Write(button.Read())
	}
}
