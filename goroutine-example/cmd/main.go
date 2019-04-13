package main

import (
	"fmt"
	"net/http"
	"os"
)

var sites = []string{
	"https://github.com",
	"https://google.com",
	"https://stackoverflow.com",
	"https://facebook.com",
	"https://twitter.com",
	"https://golang.org",
	"https://forum.golangbridge.org",
	"https://packtpub.com/",
}

func get() {
	for _, s := range sites {
		res, _ := http.Get(s)
		fmt.Printf("%s %d\n", s, res.StatusCode)
		res.Body.Close()
	}
}

func getConcurrently() {
	ch := make(chan string)

	for _, s := range sites {
		go func(s string) {
			res, _ := http.Get(s)
			ch <- fmt.Sprintf("%s %d", s, res.StatusCode)
			res.Body.Close()
		}(s)
	}

	for range sites {
		fmt.Println(<-ch)
	}
}

// goroutine's differences from threads:
// 1. its stack in memory can grow and shrink as needed
// 2. threads are scheduled with the kernal, goroutine is scheduled via gorand.time scheduler
// 3. goroutines talk to each other via channels
// 4. channel can send and receive
func main() {
	switch os.Args[1] {
	case "seq":
		get()
	case "conc":
		getConcurrently()
	default:
		fmt.Println("Please choose `seq` or `conc`.")
	}
}
