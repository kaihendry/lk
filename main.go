package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func in(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

var acceptedImageExt = []string{".png", ".gif", ".jpg", ".jpeg", ".webp"}
var images = []string{}
var dirPath = "."
var photosPerPage = 5

var chttp = http.NewServeMux()

func main() {

	flag.Parse()

	directory := flag.Arg(0)
	dirPath, _ = filepath.Abs(directory)

	fmt.Println("foo dirPath", dirPath)

	filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && in(acceptedImageExt, strings.ToLower(path.Ext(filePath))) {
			images = append(images, filePath)
		}
		return nil
	})

	fmt.Println(images)

	chttp.Handle("/", http.FileServer(http.Dir("/")))
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {

		r.ParseForm()

		pageStr := r.FormValue("page")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			fmt.Println("invalid param for page, zeroing")
			page = 0
		}

		fmt.Println("No. of images:", len(images))

		offset := page * photosPerPage
		limit := offset + photosPerPage

		if offset > len(images) {
			offset = len(images)
		}

		if limit > len(images) {
			limit = len(images)
		}

		prev := page - 1
		if prev < 0 {
			prev = 0
		}

		fmt.Println("Page:", page)

		tmplParams := struct {
			Photos []string
			Next   int
			Prev   int
		}{
			Photos: images[offset:limit],
			Next:   page + 1,
			Prev:   prev,
		}

		t, err := template.New("foo").Parse(`
		<a href="/?page={{ .Next }}">Next</a><a href="/?page={{ .Prev }}">Prev</a>
		{{ range .Photos }}<h1>{{ . }}</h1><p><img src={{ . }}></p>{{ end }}
		<a href="/?page={{ .Next }}">Next</a><a href="/?page={{ .Prev }}">Prev</a>
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t.Execute(w, tmplParams)
	}
}
