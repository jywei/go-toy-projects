package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jywei/toy-projects/chat"
)

var index = template.Must(template.ParseFiles("./index.html"))

func home(w http.ResponseWriter, r *http.Request) {
	index.Execute(w, nil)
}

func main() {
	go chat.DefaultHub.Start()

	http.HandleFunc("/", home)
	http.HandleFunc("/ws", chat.WSHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
