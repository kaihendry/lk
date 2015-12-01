package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
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

		if err != nil {
			return err
		}

		// Skip any dot files
		if strings.HasPrefix(filepath.Base(filePath), ".") {
			// fmt.Println("Skipping", filePath)
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		// Only append jpg images
		if in(acceptedImageExt, strings.ToLower(path.Ext(filePath))) {
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
			// TODO: make this spawn simultaneous jobs
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

	// http://stackoverflow.com/a/33985208/4534
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Panic(err)
	}

	hostname, _ := os.Hostname()
	if a, ok := ln.Addr().(*net.TCPAddr); ok {
		host := fmt.Sprintf("http://%s:%d", hostname, a.Port)
		fmt.Println("Serving from", host)
		open.Start(host)
	}
	if err := http.Serve(ln, nil); err != nil {
		log.Panic(err)
	}

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
