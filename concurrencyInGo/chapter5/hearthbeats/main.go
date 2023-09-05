package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

}

func heartbeatwithresults() {
	doWork := func(
		done <-chan interface{},
	) (<-chan interface{}, <-chan int) {
		heartbeat := make(chan interface{}, 1)
		workStream := make(chan int)

		go func() {
			defer close(heartbeat)
			defer close(workStream)

			for i := 0; i < 10; i++ {
				// First we ensure of sending the heartbeat before doing the work
				select {
				case heartbeat <- struct{}{}:
				default:
				}

				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()

		return heartbeat, workStream
	}

	done := make(chan interface{})
	defer close(done)

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}

		case r, ok := <-results:
			if ok {
				fmt.Printf("Reuslt: %v\n", r)
			} else {
				return
			}

		case <-time.After(timeout):
			fmt.Println("Worker gorutine is not healthy!")
			return
		}
	}
}

func withTimeout() {
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)

		go func() {
			// defer close(heartbeat)
			// defer close(results)

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(pulseInterval * 2)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()

		return heartbeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("Reuslt: %v\n", r)
		case <-time.After(timeout):
			fmt.Println("Worker gorutine is not healthy!")
			return
		}
	}
}

func lineal() {
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)

		go func() {
			defer close(heartbeat)
			defer close(results)

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(pulseInterval * 2)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()

		return heartbeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("Reuslt: %v\n", r)
		}
	}

}
