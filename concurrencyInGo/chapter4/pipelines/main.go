package main

import "fmt"

func main() {
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
