package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/warthog618/gpio"
)

func setupGPIO(pins ...int) error {
	for _, pin := range pins {
		// Export the GPIO pin
		cmd := exec.Command("bash", "-c", fmt.Sprintf("echo %d > /sys/class/gpio/export", pin))
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to export GPIO pin %d: %w", pin, err)
		}

		// Set the GPIO pin direction to low
		cmd = exec.Command("bash", "-c", fmt.Sprintf("echo low > /sys/class/gpio/gpio%d/direction", pin))
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to set GPIO pin %d direction to low: %w", pin, err)
		}
	}

	return nil
}

func main() {
	err := setupGPIO(NorthGreen, NorthYellow, NorthRed, SouthGreen, SouthYellow, SouthRed, EastGreen, EastYellow, EastRed, WestGreen, WestYellow, WestRed, Button)
	if err != nil {
		panic(err)
	}

	err = gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()

	north := [3]*gpio.Pin{
		gpio.NewPin(NorthGreen),
		gpio.NewPin(NorthYellow),
		gpio.NewPin(NorthRed),
	}
	south := [3]*gpio.Pin{
		gpio.NewPin(SouthGreen),
		gpio.NewPin(SouthYellow),
		gpio.NewPin(SouthRed),
	}
	east := [3]*gpio.Pin{
		gpio.NewPin(EastGreen),
		gpio.NewPin(EastYellow),
		gpio.NewPin(EastRed),
	}
	west := [3]*gpio.Pin{
		gpio.NewPin(WestGreen),
		gpio.NewPin(WestYellow),
		gpio.NewPin(WestRed),
	}

	const (
		Green  = 0
		Yellow = 1
		Red    = 2
	)

	button := gpio.NewPin(Button)
	button.Input()

	button.Watch(gpio.EdgeRising, func(p *gpio.Pin) {
		// Simulate pedestrian crossing
		north[Yellow].High()
		south[Yellow].High()
		time.Sleep(2 * time.Second)
		north[Green].Low()
		south[Green].Low()
		north[Red].High()
		south[Red].High()
		east[Green].High()
		west[Green].High()
	})

	for {
		// North-South green, East-West red
		north[Green].High()
		south[Green].High()
		east[Red].High()
		west[Red].High()
		time.Sleep(5 * time.Second)

		// Transition to East-West green
		north[Yellow].High()
		south[Yellow].High()
		time.Sleep(2 * time.Second)
		north[Green].Low()
		south[Green].Low()
		north[Red].High()
		south[Red].High()
		east[Green].High()
		west[Green].High()
		time.Sleep(5 * time.Second)

		// Transition to North-South green
		east[Yellow].High()
		west[Yellow].High()
		time.Sleep(2 * time.Second)
		east[Green].Low()
		west[Green].Low()
		east[Red].High()
		west[Red].High()
	}
}
