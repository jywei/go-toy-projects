package main

import (
	"fmt"

	"github.com/jywei/toy-projects/middleware"
)

func main() {
	sum := middleware.Add(1, 2, 3)
	fmt.Println(sum)
}
