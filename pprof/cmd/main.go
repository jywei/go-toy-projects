package main

import (
	"net/http"
	_ "net/http/pprof"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe(":3000", nil)
}

// go to: localhost:3000/debug/pprof/
// Use https://github.com/wg/wrk to load test
