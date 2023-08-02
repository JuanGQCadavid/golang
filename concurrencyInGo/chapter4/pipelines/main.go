package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)

		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()

	return takeStream
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)
		for {
			for _, value := range values {
				select {
				case <-done:
					return
				case valueStream <- value:
				}
			}
		}
	}()

	return valueStream
}
func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()

	return valueStream
}

func toInt(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)

		for value := range valueStream {
			select {
			case <-done:
				return
			case intStream <- value.(int):
			}
		}
	}()
	return intStream
}

func primeFounder(done <-chan interface{}, valueStream <-chan int) <-chan interface{} {
	intStream := make(chan interface{})
	go func() {
		defer close(intStream)

		for value := range valueStream {

			isPrime := true

			for i := 1; i <= value; i++ {
				if value%i == 0 {
					if i == 1 || i == value {
						continue
					}
					isPrime = false
					//break
				}
			}

			if isPrime {
				select {
				case <-done:
					return
				case intStream <- value:
				}
			}

		}
	}()
	return intStream
}

func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	joinChannel := make(chan interface{})
	var wg sync.WaitGroup

	wg.Add(len(channels))

	for _, channel := range channels {
		go func(ch <-chan interface{}) {
			defer wg.Done()
			for value := range ch {
				select {
				case <-done:
					return
				case joinChannel <- value:
				}
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		defer close(joinChannel)
	}()

	return joinChannel
}

func orDone(done, c <-chan interface{}) <-chan interface{} {
	orStream := make(chan interface{})

	go func() {
		defer close(orStream)
		for value := range c {
			select {
			case <-done:
				return
			case orStream <- value:
			}
		}
	}()

	return orStream
}

func orDone2(done, c <-chan interface{}) <-chan interface{} {
	orStream := make(chan interface{})

	go func() {
		defer close(orStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if ok == false {
					return
				}
				select {
				case orStream <- v:
				case <-done:
				}
			}
		}
	}()

	return orStream
}

// This function will send one channel data to two channels concurrently
func tee(done, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer close(out1)
		defer close(out2)

		for v := range orDone(done, in) {
			var out1, out2 = out1, out2 // Here we are shadowing the global variables

			for i := 0; i < 2; i++ { // As slect is randon we need to ensure that both of them are being called by sending a value
				// thus we iterate over the select twice
				select {
				case <-done:
					break
				case out1 <- v:
					out1 = nil // And once the channel already has the value then we put it as nil to avoid resending the value
					// Thus leaving the space to out2
				case out2 <- v:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func main() {
	done := make(chan interface{})
	//defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 10))

	// go func() {
	// 	time.Sleep(3 * time.Second)
	// 	close(done)
	// }()

	for val1 := range out1 {
		fmt.Printf("Out1: %v, out2: %v\n", val1, <-out2)
		close(done)

	}
}

func fanOutFanInt() {
	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} {
		return rand.Intn(50000000)
	}

	randIntStream := toInt(done, repeatFn(done, rand))

	numFinders := runtime.NumCPU()
	finder := make([]<-chan interface{}, numFinders)
	for i := 0; i < numFinders; i++ {
		finder[i] = primeFounder(done, randIntStream)
	}

	fmt.Println("Primes")

	start := time.Now()
	for v := range take(done, fanIn(done, finder...), 50) {
		fmt.Printf("\t%d\n", v)
	}
	fmt.Printf("It tooks: %v", time.Since(start))
}

func longPrimeFounder() {
	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} {
		return rand.Intn(50000000)
	}

	randIntStream := toInt(done, repeatFn(done, rand))

	fmt.Println("Primes")

	start := time.Now()
	for v := range take(done, primeFounder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", v)
	}
	fmt.Printf("It tooks: %v", time.Since(start))
}

func transformingToType() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})

		go func() {
			defer close(valueStream)
			for {
				for _, value := range values {
					select {
					case <-done:
						return
					case valueStream <- value:
					}
				}
			}
		}()

		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()

		return takeStream
	}

	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)

			for value := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- value.(string):
				}
			}
		}()
		return stringStream
	}

	done := make(chan interface{})
	defer close(done)

	// intStream := generator(done, 1, 2, 3, 4, 5, 6)
	// pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	var msg string
	for v := range toString(done, take(done, repeat(done, "hi", "how", "are", "you"), 10)) {
		msg += v + " "
	}

	fmt.Println(msg)
}

func repeatFnTake() {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})

		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()

		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()

		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} {
		return rand.Int()
	}
	// intStream := generator(done, 1, 2, 3, 4, 5, 6)
	// pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(v)
	}
}
func repeatTake() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})

		go func() {
			defer close(valueStream)
			for {
				for _, value := range values {
					select {
					case <-done:
						return
					case valueStream <- value:
					}
				}
			}
		}()

		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})

		go func() {
			defer close(takeStream)

			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()

		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	// intStream := generator(done, 1, 2, 3, 4, 5, 6)
	// pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range take(done, repeat(done, 1), 10) {
		fmt.Println(v)
	}
}

func streamChannelsPipelines() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int)

		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()

		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)

		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()

		return multipliedStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)

		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()

		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5, 6)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}
func streamSimplePipelines() {
	multiply := func(value, multiplier int) int {
		return value * multiplier
	}

	add := func(value, addtivie int) int {
		return value + addtivie
	}

	ints := []int{1, 2, 3, 4, 5}

	for _, v := range ints {

		// Here for every number we are instatianting three functions
		// this takes time and memory
		fmt.Println(multiply(add(multiply(v, 2), 1), 2))
	}
}

func batchSimplePipelines() {

	// This is a pipeline stage as
	// It consumes and returns the same type  ([]int)
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}

		return multipliedValues
	}

	// As this stage is taking a bunch of data and do
	// some operations on it ( receive int [], produce int[])
	// then it is call batch processing

	// One element at the time is named stream processing
	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}

		return addedValues
	}

	ints := []int{1, 2, 3, 4, 5}

	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}
