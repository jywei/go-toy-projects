package middleware

// Add is a variadic function that adds up numbers
func Add(nums ...int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// Chain holds the sum
type Chain struct {
	Sum int
}

// AddNext is chainable sum function
func (c *Chain) AddNext(num int) *Chain {
	c.Sum += num
	return c
}
