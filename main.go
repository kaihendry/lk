package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", foo)

	// We don't want this since the path to image is missing its path
	// http.Handle("/p/", http.StripPrefix("/p/", http.FileServer(http.Dir("./a/long/path/"))))

	// Wonder how to prevent file listings upon /p/
	http.Handle("/p/", http.StripPrefix("/p/", http.FileServer(http.Dir("."))))
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {

	images, _ := filepath.Glob("./a/long/path/*.jpg")

	t, err := template.New("foo").Parse(`{{ range . }}<h1>{{ . }}</h1><p><img src=/p/{{ . }}></p>{{ end }}`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, images)
}
