package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

// CreateLogger creates a new logger that writes to the given filename
func CreateLogger(filename string) *log.Logger {
	// os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666 will create a new file if not existing and grant the permission
	file, err := os.OpenFile(filename+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

// Time runs the next function in the chain
func Time(logger *log.Logger, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		logger.Println(elapsed)
	})
}

// // Add is a variadic function that adds up numbers
// func Add(nums ...int) int {
// 	sum := 0
// 	for _, num := range nums {
// 		sum += num
// 	}
// 	return sum
// }
//
// // Chain holds the sum
// type Chain struct {
// 	Sum int
// }
//
// // AddNext is chainable sum function
// func (c *Chain) AddNext(num int) *Chain {
// 	c.Sum += num
// 	return c
// }
//
// // Finally is to return the final sum of the chain
// func (c *Chain) Finally(num int) int {
// 	return c.Sum + num
// }
//
// // Next runs the next function in the chainable
// func Next(next http.HandlerFunc) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Before")
// 		next.ServeHTTP(w, r)
// 		fmt.Println("After")
// 	})
// }
