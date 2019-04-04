package cms

import "html/template"

// Tmpl is an exported variable
var Tmpl = template.Must(template.ParseGlob("../templates/*"))

type Page struct {
	Title   string
	Content string
}
