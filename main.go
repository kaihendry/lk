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

	"github.com/pyk/byten"
	"github.com/skratchdot/open-golang/open"
)

type media struct {
	filename string
	f        os.FileInfo
}

func in(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

var acceptedImageExt = []string{".jpg", ".jpeg", ".mp4", ".png"}
var dirThumbs = fmt.Sprintf("%s%s", os.Getenv("HOME"), "/.cache/lk")
var dirPath = "."
var gitVersion string
var showVersionFlag = flag.Bool("version", false, "Show version")
var port = flag.Int("port", 0, "listen port")

func hostname() string {
	hostname, _ := os.Hostname()
	// If hostname does not have dots (i.e. not fully qualified), then return zeroconf address for LAN browsing
	if strings.Split(hostname, ".")[0] == hostname {
		return hostname + ".local"
	}
	return hostname
}

func main() {

	flag.Parse()

	if *showVersionFlag {
		fmt.Println("https://github.com/kaihendry/lk", gitVersion)
		os.Exit(0)
	}

	directory := flag.Arg(0)
	dirPath, _ = filepath.Abs(directory)

	// Getting rid of /../ etc
	dirPath = path.Clean(dirPath)

	// Don't allow path under dirPath to be viewed
	http.Handle("/o/", http.StripPrefix(path.Join("/o", dirPath), http.FileServer(http.Dir(dirPath))))
	http.HandleFunc("/favicon.ico", http.NotFound)

	http.HandleFunc("/", lk)
	http.HandleFunc("/t/", thumb)

	// http://stackoverflow.com/a/33985208/4534
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Panic(err)
	}

	if a, ok := ln.Addr().(*net.TCPAddr); ok {
		host := fmt.Sprintf("http://%s:%d", hostname(), a.Port)
		log.Println("Serving from", host)
		open.Start(host)
	}
	if err := http.Serve(ln, nil); err != nil {
		log.Panic(err)
	}

}

func thumb(w http.ResponseWriter, r *http.Request) {

	// Path cleaning
	requestedPath := path.Clean(r.URL.Path[2:])

	// Make sure you can't go under the dirPath
	if !strings.HasPrefix(requestedPath, dirPath) {
		http.NotFound(w, r)
		return
	}

	thumbPath := filepath.Join(dirThumbs, requestedPath)
	if _, err := os.Stat(thumbPath); err != nil {
		log.Println("THUMB:", thumbPath, "does not exist")
		srcPath := requestedPath
		if _, err := os.Stat(srcPath); err != nil {
			log.Println("ORIGINAL", srcPath, "does not exist")
			http.NotFound(w, r)
			return
		}

		log.Println("Must generate thumb for", srcPath)
		err := genthumb(srcPath, thumbPath)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		log.Println("Created thumb", thumbPath)
	}
	log.Println("Serving thumb", thumbPath)
	http.ServeFile(w, r, thumbPath)
}

func findmedia(m *[]media) func(filename string, f os.FileInfo, err error) error {
	return func(filename string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// log.Printf("Visited: %s\n", filename)
		if !f.IsDir() && in(acceptedImageExt, strings.ToLower(path.Ext(filename))) {
			// log.Printf("Adding: %s\n%+v\n", filename, f)
			*m = append(*m, media{filename, f})
		}
		return nil
	}
}

func lk(w http.ResponseWriter, r *http.Request) {

	srcPath := filepath.Join(dirPath, r.URL.Path)

	var m []media
	err := filepath.Walk(srcPath, findmedia(&m))
	if err != nil {
		log.Println(err)
	}

	t := template.New("medialist")

	template.Must(t.Funcs(template.FuncMap{"markMedia": markMedia}).Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<style>
body { font-family: "Lucida Sans Unicode", "Lucida Grande", sans-serif; font-size: 2vw; margin: 0 }
/* https://codepen.io/dudleystorey/pen/Kgofa */
.media {
	padding: .5vw;
	flex-flow: row wrap;
	display: flex;
}
.media div * {
	width: 100%;
	height: auto;
}
.media div {
	flex: auto;
	width: 230px;
	margin: .5vw;
}
@media screen and (max-width: 400px) {
	.media div { margin: 0; }
	.media { padding: 0; }

}
</style>
</head>
<body>
<section class=media>
{{ range .Media }}<div>{{ . | markMedia }}</div>
{{ end }}
</section>
<p>By <a href=https://github.com/kaihendry/lk>lk {{ .Version }}</a></p>
</body>
</html>`))

	data := struct {
		Media   []media
		Version string
	}{
		m,
		gitVersion,
	}

	t.Execute(w, data)

	log.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.UserAgent())

}

func markMedia(m media) template.HTML {
	switch strings.ToLower(path.Ext(m.filename)) {
	case ".jpg":
		s := fmt.Sprintf("<a title=\"%s\" href=\"/o%s\"><img alt=\"\" width=230 height=230 src=\"/t%s\"></a>", m.filename, m.filename, m.filename)
		return template.HTML(s)
	case ".png":
		s := fmt.Sprintf("<a title=\"%s\" href=\"/o%s\"><img alt=\"\" width=230 height=230 src=\"/o%s\"></a>", m.filename, m.filename, m.filename)
		return template.HTML(s)
	case ".mp4":
		s := fmt.Sprintf("<video controls title=\"%s\" width=230 height=230 src=\"/o%s\"></video>", byten.Size(m.f.Size())+" "+m.filename, m.filename)
		return template.HTML(s)
	default:
		return template.HTML(m.f.Name())
	}
}
