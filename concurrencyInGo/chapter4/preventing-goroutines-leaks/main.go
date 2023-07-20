package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)

		go func() {
			defer fmt.Println("newRandStream clouse exited.")
			defer close(randStream)

			for {
				select {
				case <-done:
					return
				case randStream <- rand.Int():

				}

			}
		}()
		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)

	fmt.Println("3 random ints:")

	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

	close(done)

	time.Sleep(1 * time.Second)

}
func sendingStrings() {
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork Exited")
			defer close(terminated)

			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	messages := make(chan string)
	terminate := doWork(done, messages)

	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine")
		close(done)
	}()

loop:
	for {
		select {
		case <-terminate:
			fmt.Println("Exiting loop")
			break loop
		case messages <- "Still Alive":
		}
	}

	fmt.Println("Done.")

}
func simpleCase() {
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork Exited")
			defer close(terminated)

			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminate := doWork(done, nil)

	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine")
		close(done)
	}()

	<-terminate
	fmt.Println("Done.")
}
