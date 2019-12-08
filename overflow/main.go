package main

import (
	"fmt"
	"math"
)

func main() {
	var n int8 = math.MaxInt8

	fmt.Println("max int8    :", n)   // 127
	fmt.Println("max int8 + 1:", n+1) // -128

	n = math.MinInt8
	fmt.Println("min int8    :", n)   // -128
	fmt.Println("min int8 - 1:", n-1) // 127

	var un uint8
	fmt.Println("min uint8    :", un)   // 0
	fmt.Println("min uint8 - 1:", un-1) // 255

	f := float32(math.MaxFloat32)
	fmt.Println("max float    :", f*1.2)
}
