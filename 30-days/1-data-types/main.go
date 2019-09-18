package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var i uint64 = 4
	var d float64 = 4.0
	var s string = "HackerRank "

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	ii, err := strconv.ParseUint(scanner.Text(), 10, 64)
	if err != nil {
		fmt.Println("Error reading integer")
	}
	fmt.Println(ii + i)
	scanner.Scan()
	ff, err := strconv.ParseFloat(scanner.Text(), 64)
	if err != nil {
		fmt.Println("Error reading float")
	}
	fmt.Printf("%.1f\n", ff+d)
	scanner.Scan()
	fmt.Println(s + scanner.Text())
}
