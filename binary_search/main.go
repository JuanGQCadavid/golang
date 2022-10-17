package main

import (
	"log"
)

func main() {
	id := 10000
	// elements := []int{
	// 	1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	// }

	elements := generateArray(id)
	log.Printf("Elements: %+v\n", elements)

	search(id-1, elements)
}

func generateArray(size int) []int {
	elements := make([]int, size)

	for i := 0; i < size; i++ {
		elements[i] = i
	}

	return elements
}

func search(id int, elements []int) int {
	low, high := 0, len(elements)
	steps := 0

	for low <= high {
		steps++
		mid := (high + low) / 2
		guess := elements[mid]
		log.Println("High: ", high, " Low: ", low, " Med: ", mid, " Guess:", guess, "Steps:", steps)

		if guess == id {
			log.Println("Founded at ", mid, "with steps:", steps)
			return mid
		} else if guess > id {
			log.Println("Down!")
			high = mid - 1
		} else {
			log.Println("Up!")
			low = mid + 1
		}

	}

	log.Println("Not Founded.")
	return -1
}
