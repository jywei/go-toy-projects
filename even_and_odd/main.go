package main

import (
	"fmt"
	"strconv"
)

func main() {

	numbers := []int{}

	for i := 0; i <= 10; i++ {
		numbers = append(numbers, i)
	}

	for _, num := range numbers {
		stringNum := strconv.Itoa(num)
		if num%2 == 0 {
			fmt.Println(stringNum + " is even")
		} else {
			fmt.Println(stringNum + " is odd")
		}
	}
}
