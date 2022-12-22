package main

import (
	"fmt"
	"strconv"
)

func main() {
	x := 1234

	y := fmt.Sprintf("%d", x)

	fmt.Println(y, strconv.Itoa(x))
	fmt.Println(strconv.FormatInt(int64(x), 10))
	fmt.Println(strconv.FormatInt(int64(x), 2))

	newX, _ := strconv.ParseInt("12345", 10, 64)
	newY, _ := strconv.Atoi("12345")

	fmt.Println(newX, newY)

}
