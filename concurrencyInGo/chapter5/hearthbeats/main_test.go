package main

import (
	"fmt"
	"testing"
	"time"
)

func DoWork(
	done <-chan interface{},
	pulseInterval time.Duration,
	nums ...int,
) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	workStream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(workStream)

		time.Sleep(2 * time.Second)

		pulse := time.Tick(pulseInterval)

	numLoop:
		for _, n := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case workStream <- n:
					continue numLoop
				}
			}
		}
	}()

	return heartbeat, workStream
}

func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 4, 5}
	const timeout = 2 * time.Second
	heartbeat, results := DoWork(done, timeout/2, intSlice...)

	// Here we are waiting for the firt pulse indicating that the work is already started
	<-heartbeat
	i := 0
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}

		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; expected != r {
				t.Errorf("Index %v; expected %v, but received %v,", i, expected, r)
			}

		case <-time.After(timeout):
			t.Fatal("Test timed out")
			return
		}
		i++
	}
}
