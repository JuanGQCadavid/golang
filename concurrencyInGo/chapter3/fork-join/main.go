package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	wg.Add(1)
	go sayHello() // Fork point
	fmt.Println("Sup!")

	wg.Add(1)
	go func() {
		fmt.Println("Ey bro!")
		wg.Done()
	}() // Fork point

	wg.Add(1)
	anOther := func() {
		fmt.Println("Tell me")
		wg.Done()
	}

	go anOther() // Fork point

	wg.Wait() // joint point

}

func sayHello() {
	fmt.Println("Hello")
	wg.Done()
}
