package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewReader(os.Stdin)
	input, _ := scanner.ReadString('\n')
	fmt.Println("Hello, World.")
	fmt.Println(input)
}
