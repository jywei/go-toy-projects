package main

import (
	"flag"
	"net/http"
	"os"
)

func main() {
	var dir string
	port := flag.String("port", "3000", "port to serve HTTP on")
	// it's actually a pointer to the string
	path := flag.String("path", "", "path to serve")
	flag.Parse()

	// Need to dereference the pointer
	if *path == "" {
		dir, _ = os.Getwd()
	} else {
		dir = *path
	}

	http.ListenAndServe(":"+*port, http.FileServer(http.Dir(dir)))
}

// go run fs.go -port=5151 -path=site
