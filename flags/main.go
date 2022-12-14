package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	n   = flag.Bool("n", false, "Omit trailing new line")
	sep = flag.String("s", " ", "Separator")
)

func main() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *sep))

	if !*n {
		fmt.Println()
	}
}
