package api

import (
	"io"
	"net/http"
	"os"
	"strings"
)

var here = os.Getenv("GOPATH") + "/src/github.com/jywei/toy-projects/api/images/"

// UploadImage allows uploading an image
func UploadImage(w http.ResponseWriter, r *http.Request) {
	// 1. get the image data and header from the request
	// Content-Encoding: multipart/form-data
	file, header, err := r.FormFile("image")
	if err != nil {
		errJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 2. create a new file in image directory with os.Create
	out, err := os.Create(here + header.Filename)
	if err != nil {
		errJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 3. Write the image to the file
	_, err = io.Copy(out, file)
	if err != nil {
		errJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. response with json
	writeJSON(w, map[string]string{
		"filename": header.Filename,
	})
}

// ShowImage shows the image based on the filename found in the path
func ShowImage(w http.ResponseWriter, r *http.Request) {
	// 1. Get the file name from the url
	name := strings.TrimLeft(r.URL.Path, "/image/")
	// 2. Open the file
	file, err := os.Open(here + name)
	if err != nil {
		errJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := pool.Get()
	defer pool.Put(buf)
	// 3. Move the data from our file to our bufferpool
	_, err = io.Copy(buf, file)
	if err != nil {
		errJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Set the content-type and write the image to our ResponseWriter
	w.Header().Set("Content-Type", "image/jpeg")
	buf.WriteTo(w)
}
