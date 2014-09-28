package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

	http.Handle("/o/", http.StripPrefix("/o/", http.FileServer(http.Dir("/"))))
	http.Handle("/t/", http.StripPrefix("/t/", (http.FileServer(http.Dir("/home/hendry/.cache/sxiv")))))
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("right", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func foo(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("foo").Parse(`{{ range . }}<a title={{ . }} href=/o/{{ . }}><img width=160 src="/t/{{ . }}.jpg"></a>{{ end }}`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, images)
}
