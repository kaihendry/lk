package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

var chttp = http.NewServeMux()

func main() {
	chttp.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {

		images, _ := filepath.Glob("./a/long/path/*.jpg")

		t, err := template.New("foo").Parse(`{{ range . }}<h1>{{ . }}</h1><p><img src=./{{ . }}></p>{{ end }}`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t.Execute(w, images)
	}
}
