package main

import (
	"fmt"
	"reflect"
)

type a struct {
	B string
	C int
}

func main() {
	x := &a{
		B: "B",
		C: 1,
	}
	y := &a{
		B: "B",
		C: 1,
	}
	fmt.Println(x == y)
	fmt.Println(&x, &y)

	fmt.Println(reflect.DeepEqual(x, y))
}
