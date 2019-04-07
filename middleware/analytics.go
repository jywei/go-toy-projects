package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ErrInvalidID    = errors.New("Invalid ID")
	ErrInvalidEmail = errors.New("Invalid Email")
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

// Recover will recover from any panicking goroutine
func Recover(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				switch err {
				case ErrInvalidEmail:
					http.Error(w, ErrInvalidEmail.Error(), http.StatusUnauthorized)
				case ErrInvalidID:
					http.Error(w, ErrInvalidID.Error(), http.StatusUnauthorized)
				default:
					http.Error(w, "Unknown error, recovered from panic", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// PassContext is used to pass values between middleware
type PassContext func(ctx context.Context, w http.ResponseWriter, r *http.Request)

// ServeHTTP satisfies the http.Handler interface
func (fn PassContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// use ctx as a hash
	ctx := context.WithValue(context.Background(), "foo", "bar")
	// fn method will execute Http.handler
	fn(ctx, w, r)
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
