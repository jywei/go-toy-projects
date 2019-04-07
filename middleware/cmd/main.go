package main

import (
	"net/http"

	"github.com/jywei/toy-projects/middleware"
)

func hello(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Executing...")
	w.Write([]byte("Hello"))
}

// func tricky() string {
// 	defer log.Println("String 2")
// 	return "String 1" // Will print 2 1
// }
//
// func lastInFirstOut() {
// 	for i := 0; i < 4; i++ {
// 		defer log.Println(i)
// 	} // Will print 3 2 1 0
// }

func panicker(w http.ResponseWriter, r *http.Request) {
	panic("Wahhhh")
}

func main() {
	// sum := middleware.Add(1, 2, 3)
	// fmt.Println(sum)
	//
	// chain := &middleware.Chain{0}
	// sum2 := chain.AddNext(1).AddNext(2).AddNext(3).Finally(0)
	// fmt.Println(sum2)

	// http.Handle("/", middleware.Next(hello))
	logger := middleware.CreateLogger("section4")
	http.Handle("/", middleware.Time(logger, hello))
	http.Handle("/panic", middleware.Recover(panicker))
	http.ListenAndServe(":3000", nil)
}
