package main

// Adder interface
type Adder interface {
	Add(a, b int) int
}

// X type is considered as an Adder interface
type X struct{}

// Add function
func (x *X) Add(a, b int) int {
	return a + b
}
