package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println(runtime.NumCPU())

	done := make(chan interface{})

	counter := func(done <-chan interface{}) <-chan int {
		wg := sync.WaitGroup{}
		result := make(chan int, 1)
		wg.Add(1)
		go func(done <-chan interface{}, result chan<- int) {
			wg.Done()
			counter := 0
			for {
				select {
				case <-done:
					result <- counter
					return
				default:
					counter++
				}
			}
		}(done, result)

		wg.Wait()
		return result
	}

	result := counter(done)
	fmt.Println("Going into sleep")
	time.Sleep(5 * time.Second)
	close(done)

	fmt.Println(<-result)

}
