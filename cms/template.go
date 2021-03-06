package cms

import (
	"html/template"
	"os"
	"time"
)

var tmplPath = os.Getenv("GOPATH") + "/src/github.com/jywei/toy-projects/cms/templates"

// Tmpl is a reference to all of our templates
// ParseGlob would return a template and error, and Must will do the eror checking
var Tmpl = template.Must(template.ParseGlob(tmplPath))

// Page is the struct used for each webpage
type Page struct {
	ID      int
	Title   string
	Content string
	Posts   []*Post
}

// Post is the struct used for each blog post
type Post struct {
	ID            int
	Title         string
	Content       string
	DatePublished time.Time
	Comments      []*Comment
}

// Comment is the struct used for each comment
type Comment struct {
	ID            int
	PostID        int
	Author        string
	Comment       string
	DatePublished time.Time
}
