package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/3d0c/imgproc"
)

func genthumb(src string, dst string) (err error) {

	fmt.Println("Resizing:", src, dst)

	base := imgproc.NewSource(src)
	if base == nil {
		log.Fatal(base)
	}

	target := &imgproc.Options{
		Base:    base,
		Scale:   imgproc.NewScale("800"),
		Method:  3,
		Format:  "jpg",
		Quality: 100,
	}

	base = imgproc.NewSource(imgproc.Proc(target))

	// Crop 100x100 pixel from center
	target = &imgproc.Options{
		Base:    base,
		Crop:    imgproc.NewRoi("center,400,400"),
		Method:  3,
		Format:  "jpg",
		Quality: 80,
	}

	t := imgproc.Proc(target)

	if t != nil {
		dir, _ := filepath.Split(dst)
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return
		}
	}

	err = ioutil.WriteFile(dst, t, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return

}
