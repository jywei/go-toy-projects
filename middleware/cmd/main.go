package main

import (
	"fmt"

	"github.com/jywei/toy-projects/middleware"
)

func main() {
	sum := middleware.Add(1, 2, 3)
	fmt.Println(sum)

	chain := &middleware.Chain{0}
	sum2 := chain.AddNext(1).AddNext(2).AddNext(3)
	fmt.Println(sum2.Sum)
}
