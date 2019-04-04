package main

import (
	"os"

	"github.com/jywei/toy-projects/cms"
)

func main() {
	// p is a reference of cms.Page
	p := &cms.Page{
		Title:   "Hello, world!",
		Content: "This is the body of our webpage",
	}

	cms.Tmpl.ExecuteTemplate(os.Stdout, "index", p)
}
