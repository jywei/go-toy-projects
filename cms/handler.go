package cms

import (
	"net/http"
	"strings"
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
			p := &Page{
				Title:   title,
				Content: content,
			}
			_, err := CreatePage(p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				// return, otherwise the func will continue to execute
				return
			}
			Tmpl.ExecuteTemplate(w, "page", p)
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

// ServePage serves a page based on the route matched. This will match any URL
// beginning with /page
func ServePage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/page/")

	if path == "" {
		http.NotFound(w, r)
		return
	}

	p := &Page{
		Title:   strings.ToTitle(path),
		Content: "Here is my page",
	}

	Tmpl.ExecuteTemplate(w, "page", p)
}

// ServePost serves a post
func ServePost(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/post/")

	if path == "" {
		http.NotFound(w, r)
		return
	}

	p := &Post{
		Title:   strings.ToTitle(path),
		Content: "Here is my page",
		Comments: []*Comment{
			&Comment{
				Author:        "Roy Wei",
				Comment:       "Looks great!",
				DatePublished: time.Now(),
			},
		},
	}

	Tmpl.ExecuteTemplate(w, "post", p)
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
