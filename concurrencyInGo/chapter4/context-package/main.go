package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := cprintGretting(ctx); err != nil {
			fmt.Println("Error: ", err.Error())
			cancel()
			return
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := cprintFarewell(ctx); err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
	}()

	wg.Wait()

}

// With context

func cprintGretting(ctx context.Context) error {
	greeting, err := cgenGreeting(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s world! \n", greeting)
	return nil
}

func cprintFarewell(ctx context.Context) error {
	farewell, err := cgenFarewell(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s world! \n", farewell)
	return nil
}

func cgenGreeting(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := clocale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func cgenFarewell(ctx context.Context) (string, error) {
	switch locale, err := clocale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func clocale(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}

// Without Context

func withoutContext() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := printGretting(done); err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
	}()

	wg.Wait()
}

func printGretting(done <-chan interface{}) error {
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}

	fmt.Printf("%s world! \n", greeting)
	return nil
}

func printFarewell(done <-chan interface{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}

	fmt.Printf("%s world! \n", farewell)
	return nil
}

func genGreeting(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(done <-chan interface{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func locale(done <-chan interface{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("Canceled")
	case <-time.After(1 * time.Minute):
	}
	return "EN/US", nil
}
