package main

import (
	"net/http"

	"github.com/jywei/toy-projects/cms"
)

func main() {
	http.HandleFunc("/", cms.ServeIndex)
	http.HandleFunc("/new", cms.HandleNew)
	http.ListenAndServe(":3000", nil)
}
