package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jywei/toy-projects/cms"
)

// Doc lists all the routes for our API
func Doc(w http.ResponseWriter, r *http.Request) {
	data := (map[string]string{
		"all_pages_url":   "/pages",
		"page_url":        "/pages/{id}",
		"create_page_url": "/newpage",
	})
	writeJSON(w, data)
}

// AllPages return all the pages
func AllPages(w http.ResponseWriter, r *http.Request) {
	data, err := cms.GetPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, data)
}

// CreatePage creates a new post or pages
func CreatePage(w http.ResponseWriter, r *http.Request) {
	page := new(cms.Page)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Unmarshal(data, page)
	id, err := cms.CreatePage(page)
	if err != nil {
		errJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]int{
		"user_id": id,
	})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		errJSON(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(resJSON)
}

func errJSON(w http.ResponseWriter, err string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte("{\n\terror: " + err + "\n}\n"))
}

// GetPage gets a single page from the API
func GetPage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimLeft(r.URL.Path, "/pages/")
	data, err := cms.GetPage(id)
	if err != nil {
		errJSON(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, data)
}
