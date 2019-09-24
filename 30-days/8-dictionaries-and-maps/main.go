package main

import (
	"fmt"
)

func main() {
	var n int
	m := make(map[string]int)

	fmt.Scan(&n)

	var name string
	var num int
	for i := 0; i < n; i++ {
		fmt.Scanf("%v %v", &name, &num)
		m[name] = num
	}

	var query string
	for {
		_, err := fmt.Scan(&query)
		if err != nil {
			break
		}
		if value, ok := m[query]; ok {
			fmt.Printf("%s=%d\n", query, value)
		} else {
			fmt.Println("Not found")
		}
	}

}
