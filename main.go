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

var acceptedImageExt = []string{".jpg", ".jpeg"}
var images = []string{}
var dirThumbs = fmt.Sprintf("%s%s", os.Getenv("HOME"), "/.cache/sxiv")
var dirPath = "."

func main() {

	flag.Parse()

	directory := flag.Arg(0)
	dirPath, _ = filepath.Abs(directory)

	fmt.Println("lk is serving", dirPath, "from http://0.0.0.0:3000")

	filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && in(acceptedImageExt, strings.ToLower(path.Ext(filePath))) {
			thumbnail := fmt.Sprintf("%s%s.jpg", dirThumbs, filePath)
			if _, err := os.Stat(thumbnail); os.IsNotExist(err) {
				fmt.Println("Missing thumbnail:", thumbnail)
				genthumb(filePath, thumbnail)
			}
			images = append(images, filePath)
		}
		return nil
	})

	// Don't allow path under dirPath to be viewed
	// http://www.reddit.com/r/golang/comments/2l59wk/web_based_jpg_viewer_for_sharing_images_on_a_lan/clrpbyo
	http.Handle("/o/", http.StripPrefix(path.Join("/o", dirPath), http.FileServer(http.Dir(dirPath))))
	http.Handle("/t/", http.StripPrefix(path.Join("/t", dirPath), http.FileServer(http.Dir(path.Join(dirThumbs, dirPath)))))

	http.HandleFunc("/", lk)
	http.ListenAndServe(":3000", nil)
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("right", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func lk(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("foo").Parse(`{{ range . }}<a title={{ . }} href=/o{{ . }}>
<img width=160 src="/t{{ . }}.jpg">
</a>
{{ end }}`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, images)
}
