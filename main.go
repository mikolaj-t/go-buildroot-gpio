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

	north := TrafficLight{
		Green:  gpio.NewPin(NorthGreen),
		Yellow: gpio.NewPin(NorthYellow),
		Red:    gpio.NewPin(NorthRed),
	}
	south := TrafficLight{
		Green:  gpio.NewPin(SouthGreen),
		Yellow: gpio.NewPin(SouthYellow),
		Red:    gpio.NewPin(SouthRed),
	}
	east := TrafficLight{
		Green:  gpio.NewPin(EastGreen),
		Yellow: gpio.NewPin(EastYellow),
		Red:    gpio.NewPin(EastRed),
	}
	west := TrafficLight{
		Green:  gpio.NewPin(WestGreen),
		Yellow: gpio.NewPin(WestYellow),
		Red:    gpio.NewPin(WestRed),
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
					north.ChangeState(Yellow)
					south.ChangeState(Yellow)
					east.ChangeState(Yellow)
					west.ChangeState(Yellow)
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
				north.ChangeState(Green)
				south.ChangeState(Green)
				east.ChangeState(Red)
				west.ChangeState(Red)

				time.Sleep(5 * time.Second)

				// Transition to East-West green
				north.ChangeState(Yellow)
				south.ChangeState(Yellow)
				east.ChangeState(Yellow)
				west.ChangeState(Yellow)

				time.Sleep(2 * time.Second)

				// North-South red, East-West green
				north.ChangeState(Red)
				south.ChangeState(Red)
				east.ChangeState(Green)
				west.ChangeState(Green)

				time.Sleep(5 * time.Second)

				// Transition to North-South green
				north.ChangeState(Yellow)
				south.ChangeState(Yellow)
				east.ChangeState(Yellow)
				west.ChangeState(Yellow)

				time.Sleep(2 * time.Second)
			}
		}
	}()
	wg.Wait()
	fmt.Println("Goodbye!")
}
