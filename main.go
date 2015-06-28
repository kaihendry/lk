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
	"time"

	"github.com/skratchdot/open-golang/open"
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
var dirThumbs = fmt.Sprintf("%s%s", os.Getenv("HOME"), "/.cache/lk")
var dirPath = "."
var gitVersion string
var showVersionFlag = flag.Bool("version", false, "Show version")

func main() {

	flag.Parse()

	if *showVersionFlag {
		fmt.Println("lk", gitVersion, "https://github.com/kaihendry/lk")
		os.Exit(0)
	}

	directory := flag.Arg(0)
	dirPath, _ = filepath.Abs(directory)

	filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && in(acceptedImageExt, strings.ToLower(path.Ext(filePath))) && !strings.HasPrefix(path.Base(filePath), ".") {
			log.Printf("Appending %s", filePath)
			images = append(images, filePath)
		}
		return nil
	})

	imgLength := len(images)

	start := time.Now()
	thumbsGenerated := 0
	for i, filePath := range images {
		thumbnail := fmt.Sprintf("%s%s.jpg", dirThumbs, filePath)
		if _, err := os.Stat(thumbnail); os.IsNotExist(err) {
			fmt.Printf("%3.f%% %s\n", ((float64(i)+1)/float64(imgLength))*100, thumbnail)
			genthumb(filePath, thumbnail)
			thumbsGenerated++
		}
	}
	elapsed := time.Since(start)
	if thumbsGenerated > 0 {
		log.Printf("Generating %d thumbs took %s", thumbsGenerated, elapsed)
	}

	// Don't allow path under dirPath to be viewed
	http.Handle("/o/", http.StripPrefix(path.Join("/o", dirPath), http.FileServer(http.Dir(dirPath))))
	http.Handle("/t/", http.StripPrefix(path.Join("/t", dirPath), http.FileServer(http.Dir(path.Join(dirThumbs, dirPath)))))
	http.HandleFunc("/favicon.ico", http.NotFound)

	http.HandleFunc("/", lk)
	fmt.Println("lk is serving", dirPath, "from http://0.0.0.0:3000")
	open.Start("http://0.0.0.0:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func lk(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("foo").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
</head>
<body>
{{ range . }}<a title="{{ . }}" href="/o{{ . }}">
<img alt="" width=230 height=230 src="/t{{ . }}.jpg">
</a>
{{ end }}
<p>By <a href=https://github.com/kaihendry/lk>lk</a></p>
</body>
</html>`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(w, images)

	log.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.UserAgent())

}
