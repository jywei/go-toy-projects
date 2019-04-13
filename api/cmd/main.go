package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jywei/toy-projects/api"
)

func main() {
	// Create images directory
	os.Mkdir("images", 0777)

	http.HandleFunc("/", api.Doc)
	http.HandleFunc("/image/", api.ShowImage)
	http.HandleFunc("/newpage", api.CreatePage)
	http.HandleFunc("/pages", api.AllPages)
	http.HandleFunc("/pages/", api.GetPage)
	http.HandleFunc("/upload", api.UploadImage)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
