package main

import "fmt"

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func main() {
	fmt.Printf("%T %[1]v\n", Sunday)
	fmt.Printf("%T %[1]v\n", Monday)
	fmt.Printf("%T %[1]v\n", Tuesday)
	fmt.Printf("%T %[1]v\n", Wednesday)
	fmt.Printf("%T %[1]v\n", Thursday)
	fmt.Printf("%T %[1]v\n", Friday)
	fmt.Printf("%T %[1]v\n", Saturday)
}
