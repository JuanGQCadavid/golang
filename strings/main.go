package main

import "fmt"

func main() {
	s := "Hello, world\n"
	fmt.Println(s[0], string(s[0]))
	fmt.Println(len(s))
	fmt.Println(string(s[0:5]))
	fmt.Println(string(s[5:]), string(s[:5]))
	fmt.Println(string(s[:]))

	d := "\a Hu?"

	fmt.Println(d)

	const GoUsage = `Go usage:
	
	Usage:
		blha blah blah
	`

}
