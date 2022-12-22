package main

import (
	"fmt"
	"time"
)

const (
	FiveSeconds time.Duration = 5 * time.Second
	TwoSeconds                = 2 * time.Second
)

const (
	a = 1
	b
	c = 2
	d
)

func main() {
	fmt.Printf("%T %[1]v\n", FiveSeconds)
	fmt.Printf("%T %[1]v\n", TwoSeconds)

	fmt.Println(a, b, c, d)
}
