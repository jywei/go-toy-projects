package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func hanoi(n int, a string, b string, c string) {
	if n == 1 {
		fmt.Println(a, "->", c)
	} else {
		hanoi(n-1, a, c, b)
		hanoi(1, a, b, c)
		hanoi(n-1, b, a, c)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	num, _ := strconv.Atoi(scanner.Text())
	hanoi(num, "A", "B", "C")
}
