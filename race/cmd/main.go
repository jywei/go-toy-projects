package main

import (
	"fmt"
	"math/rand"
	"time"
)

// N is N
var N = 1

func main() {
	rand.Seed(time.Now().UnixNano())
	go mutateAtRandom()
	mutateAtRandom()

	fmt.Println(N == 1, N)
}

func mutateAtRandom() {
	v := time.Duration(rand.Intn(3))
	time.Sleep(v * time.Second)
	N = int(v)
}
