package main

import (
	"flag"
	"fmt"
	"html/template"
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

var chttp = http.NewServeMux()

func main() {

	flag.Parse()

	directory := flag.Arg(0)
	dirPath, _ = filepath.Abs(directory)

	fmt.Println("foo dirPath", dirPath)

	chttp.Handle("/", http.FileServer(http.Dir("/")))
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {

		fmt.Println("foo dirPath", dirPath)

		filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
			if err == nil && in(acceptedImageExt, strings.ToLower(path.Ext(filePath))) {
				images = append(images, filePath)
			}
			return nil
		})

		fmt.Println(images)

		t, err := template.New("foo").Parse(`{{ range . }}<h1>{{ . }}</h1><p><img src={{ . }}></p>{{ end }}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t.Execute(w, images)
	}
}
