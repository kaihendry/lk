package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nfnt/resize"
)

func genthumb(src string, dst string) (err error) {

	dir, _ := filepath.Split(dst)
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return
	}

	path, err := exec.LookPath("vipsthumbnail")
	if err == nil {
		out, err := exec.Command(path, "-s", "460x460", "-o", dst, src).CombinedOutput()
		if err != nil {
			fmt.Printf("The output is %s\n", out)
			log.Fatal(err)
		}
		return err
	}

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

	out, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	return

}
