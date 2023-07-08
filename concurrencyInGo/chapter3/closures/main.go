package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	salutation := "Hello"

	// The gorutine will be created within the same memory space that contains the
	// salutation as the context that it was created contains it, so it could change it
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "Welcome"

	}()

	wg.Wait()
	fmt.Println(salutation)

	for _, greetings := range []string{"Hello", "Greetings", "Good day"} {
		wg.Add(1)

		// Where it is so possible that the loop will finish its execution even before
		// the first gorutine gets CPU time, thus it is so possible that the output is the last
		// value that was assigned to the greetings variable, as golang keeps a track over the variables
		// pointers it knows that greetings is being used so it will put it on the heap to make it
		// visibale
		go func() {
			defer wg.Done()
			fmt.Println(greetings)
		}()
	}

	wg.Wait()

	for _, greetings := range []string{"Hello", "Greetings", "Good day"} {
		wg.Add(1)

		// Here we are passing teh value as a copy as a parameter, so we avoid the problem
		// of pointing to the same address on the heap.
		go func(greeting string) {
			defer wg.Done()
			fmt.Println(greeting)
		}(greetings)
	}
	wg.Wait()

}
