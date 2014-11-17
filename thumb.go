package main

import (
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func genthumb(src string, dst string) (err error) {

	file, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Thumbnail(460, 460, img, resize.NearestNeighbor)

	dir, _ := filepath.Split(dst)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return
	}
	out, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	return

}
