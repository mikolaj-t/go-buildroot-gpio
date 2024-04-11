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
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo %d > /sys/class/gpio/export", pin))
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to export GPIO pin %d: %w", pin, err)
		}

		// Set the GPIO pin direction to low
		cmd = exec.Command("sh", "-c", fmt.Sprintf("echo low > /sys/class/gpio/gpio%d/direction", pin))
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to set GPIO pin %d direction to low: %w", pin, err)
		}
	}

	return nil
}

func cleanupGPIO(pins ...int) error {
	for _, pin := range pins {
		// Unexport the GPIO pin
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo %d > /sys/class/gpio/unexport", pin))
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to unexport GPIO pin %d: %w", pin, err)
		}
	}

	return nil
}

func main() {
	err := setupGPIO(NorthGreen, NorthYellow, NorthRed, SouthGreen, SouthYellow, SouthRed, EastGreen, EastYellow, EastRed, WestGreen, WestYellow, WestRed, Button)
	if err != nil {
		panic(err)
	}
	defer cleanupGPIO(NorthGreen, NorthYellow, NorthRed, SouthGreen, SouthYellow, SouthRed, EastGreen, EastYellow, EastRed, WestGreen, WestYellow, WestRed, Button)

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
		// Toggle all the lights

		north[Green].Toggle()
		north[Yellow].Toggle()
		north[Red].Toggle()

		south[Green].Toggle()
		south[Yellow].Toggle()
		south[Red].Toggle()

		east[Green].Toggle()
		east[Yellow].Toggle()
		east[Red].Toggle()

		west[Green].Toggle()
		west[Yellow].Toggle()
		west[Red].Toggle()
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
		west[Yellow].High()
		east[Yellow].High()

		north[Green].Low()
		south[Green].Low()

		time.Sleep(2 * time.Second)

		// North-South red, East-West green
		north[Yellow].Low()
		south[Yellow].Low()
		north[Red].High()
		south[Red].High()

		west[Red].Low()
		east[Red].Low()
		west[Yellow].Low()
		east[Yellow].Low()

		west[Green].High()
		east[Green].High()

		time.Sleep(5 * time.Second)

		// Transition to North-South green
		north[Yellow].High()
		south[Yellow].High()
		west[Yellow].High()
		east[Yellow].High()

		west[Green].Low()
		east[Green].Low()

		time.Sleep(2 * time.Second)
		north[Yellow].Low()
		south[Yellow].Low()
		west[Yellow].Low()
		east[Yellow].Low()
	}
}
