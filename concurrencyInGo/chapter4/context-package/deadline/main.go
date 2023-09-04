package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreeting(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printFarewell(ctx); err != nil {
			fmt.Printf("cannot print greeting: %v", err)
			cancel()
		}
	}()

	wg.Wait()
}

func printGreeting(ctx context.Context) error {

	greeting, err := genGreeting(ctx)

	if err != nil {
		return err
	}
	
	fmt.Print()
}


func printFarewall(ctx context.Context) error {

}
