package main

import (
	"fmt"
	"net/http"

	"github.com/jywei/toy-projects/middleware"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Executing...")
	w.Write([]byte("Hello"))
}

func main() {
	sum := middleware.Add(1, 2, 3)
	fmt.Println(sum)

	chain := &middleware.Chain{0}
	sum2 := chain.AddNext(1).AddNext(2).AddNext(3).Finally(0)
	fmt.Println(sum2)

	http.Handle("/", middleware.Next(hello))
	http.ListenAndServe(":3000", nil)
}
