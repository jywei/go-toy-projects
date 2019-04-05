package cms

import (
	"net/http"
	"time"
)

// HandleNew handles preview NewCatalogService
func HandleNew(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		Tmpl.ExecuteTemplate(w, "new", nil)

	case "POST":
		title := req.FormValue("title")
		content := req.FormValue("content")
		contentType := req.FormValue("contentType")
		req.ParseForm()

		if contentType == "page" {
			Tmpl.ExecuteTemplate(w, "page", &Page{
				Title:   title,
				Content: content,
			})
			return
		}

		if contentType == "post" {
			Tmpl.ExecuteTemplate(w, "post", &Post{
				Title:   title,
				Content: content,
			})
			return
		}
	default:
		http.Error(w, "Method not supported: "+req.Method, http.StatusMethodNotAllowed)
	}
}

func ServeIndex(w http.ResponseWriter, req *http.Request) {
	p := &Page{
		Title:   "Go Projects CMS",
		Content: "Welcome to our home page!",
		Posts: []*Post{
			&Post{
				Title:         "Hello, World!",
				Content:       "Hello world! Thanks for coming to the site.",
				DatePublished: time.Now(),
			},
			&Post{
				Title:         "A Post with Comments",
				Content:       "Here's a controversial post. It's sure to attract comments.",
				DatePublished: time.Now().Add(-time.Hour),
				Comments: []*Comment{
					&Comment{
						Author:        "Roy Wei",
						Comment:       "Nevermind, I guess I just commented on my own post...",
						DatePublished: time.Now().Add(-time.Hour / 2),
					},
				},
			},
		},
	}

	Tmpl.ExecuteTemplate(w, "page", p)
}
