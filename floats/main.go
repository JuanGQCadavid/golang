package main

import (
	"fmt"
	"math"
)

func main() {

	// Float 32 -> six decimal digis of precission
	// float64 -> 15 digits of precision
	var f float32 = 16777216 // 1<<24
	fmt.Println(f == f+1)

	f = 0.1
	f = .1
	f = 6.34e-1

	fmt.Println(f < math.MaxFloat32)

	var z float64

	fmt.Println(z, -z, 1/z, -1/z, z/z)

	fmt.Println(math.IsNaN(z / z))
	fmt.Println(math.IsNaN(1 / z))
	fmt.Println(math.IsInf(1/z, 1))
	fmt.Println((z / z) == math.NaN())

}
