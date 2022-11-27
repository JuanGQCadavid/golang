package main

import "fmt"

func main() {
	//binaryOperations()
	//converting()
	//octaAndHexadecimal()
	runeOperations()
}

func runeOperations() {
	ascii := 'a'
	newLine := '\n'
	fmt.Printf("%d %[1]c %[1]q \n", ascii)
	fmt.Printf("%d %[1]q \n", newLine)
}

func octaAndHexadecimal() {
	o := 0666

	fmt.Printf("%d %[1]o %#[1]o \n", o) // the # adverb for %o or %x
}

func converting() {
	f := 3.141 // a float64
	i := int(f)

	fmt.Println(f, i)

	f = 1.99
	fmt.Println(int(f)) // truncates towards zero
}

func unaryAddictions() {
	x := -1
	fmt.Println(x)
	fmt.Println(+x)
	fmt.Println(-x)
}

func binaryOperations() {
	var x uint8 = 1<<1 | 1<<5
	var y uint8 = 1<<1 | 1<<2

	fmt.Printf("%08b - %v \n", x, x)                       // The set {1,5}
	fmt.Printf("%08b - %v \n", y, y)                       // The set {1,2}
	fmt.Printf("%08b - %v \n", x&y, x&y)                   // The interseption {1}
	fmt.Printf("%08b - %v \n", x|y, x|y)                   // The union {1,5,2}
	fmt.Printf("%08b - %v \n", x^y, x^y)                   // The symetric difference  {5,2}
	fmt.Printf("%08b - %v \n", x&^y, x&^y)                 // The difference  {5,2}
	fmt.Printf("%08b - %v <-> %08b - %v \n", y, y, ^y, ^y) // Y negation  {5,2}

	fmt.Println("Membership test")
	for i := uint(0); i < 8; i++ {
		if (x|y)&(1<<i) != 0 { // membership test
			fmt.Println(i) // "1", "5"
		}
	}
	fmt.Printf("%08b - %v <-> %08b - %v \n", x, x, x<<1, x<<1) // Move x one row -> {2, 6}
	fmt.Printf("%08b - %v <-> %08b - %v \n", x, x, x>>1, x>>1) // Move x one row -> {0, 4}

}
