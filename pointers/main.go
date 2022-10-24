package main

import "fmt"

func main() {
	v := 1
	inc(&v)
	fmt.Println("v:", v)

	v5 := newInt()
	fmt.Println("v5:", *v5)

	fmt.Println("newInt() == newInt(): ", newInt() == newInt())
	aliases()
}
func aliases() {
	a := 1
	p := &a // p points to a, p is an alias for a

	fmt.Println("p:", *p, "a:", a)
	*p = 5
	fmt.Println("p:", *p, "a:", a)

}
func inc(p *int) int {
	*p++
	return *p
}

func newInt() *int {
	v := 5
	return &v
}
