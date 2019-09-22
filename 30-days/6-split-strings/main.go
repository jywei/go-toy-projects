package main

import (
	"bytes"
	"fmt"
)

func main() {
	// Enter your code here. Read input from STDIN. Print output to STDOUT
	var T int
	var input string

	fmt.Scan(&T)

	for i := 0; i < T; i++ {
		fmt.Scan(&input)

		odd := bytes.Buffer{}
		even := bytes.Buffer{}
		for j := 0; j < len(input); j++ {
			if j%2 == 0 {
				even.WriteString(string(input[j]))
			} else {
				odd.WriteString(string(input[j]))
			}
		}
		fmt.Println(even.String(), odd.String())
	}

}
