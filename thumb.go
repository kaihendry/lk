package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	img "github.com/3d0c/imgproc"
)

func genthumb(src string, dst string) (err error) {

	fmt.Println("Resizing:", src, dst)

	base := img.NewSource(src)
	if base == nil {
		log.Fatal(base)
	}

	target := &img.Options{
		Base:    base,
		Scale:   img.NewScale("160x"),
		Method:  3,
		Quality: 80,
	}

	target.Format = "jpg"

	t := img.Proc(target)

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
