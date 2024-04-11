package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
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
	err := setupGPIO(NorthGreen, NorthYellow, NorthRed, SouthGreen, SouthYellow, SouthRed, EastGreen, EastYellow, EastRed, WestGreen, WestYellow, WestRed)
	if err != nil {
		panic(err)
	}
	defer cleanupGPIO(NorthGreen, NorthYellow, NorthRed, SouthGreen, SouthYellow, SouthRed, EastGreen, EastYellow, EastRed, WestGreen, WestYellow, WestRed, Button)

	err = gpio.Open()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	stopAutomatic := make(chan struct{})
	stopButton := make(chan struct{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		stopAutomatic <- struct{}{}
		stopButton <- struct{}{}
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

	stableState := make(chan bool)
	unbouncer := NewUnbouncer(stableState, 100*time.Millisecond)

	button.Watch(gpio.EdgeBoth, func(p *gpio.Pin) {
		fmt.Println("Button pressed!", button.Read())
		unbouncer.OnClicked(bool(button.Read()))
	})

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopButton:
				return
			case state := <-stableState:
				if !state {
					// Toggle all the lights
					fmt.Println("Button pressed stable!", state)

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
				}
			}

		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopAutomatic:
				return
			default:
				// North-South green, East-West red
				north[Yellow].Low()
				south[Yellow].Low()
				west[Yellow].Low()
				east[Yellow].Low()
				north[Red].Low()
				south[Red].Low()

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
			}
		}
	}()
	wg.Wait()
	fmt.Println("Goodbye!")
}
