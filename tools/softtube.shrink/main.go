package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"

	core "github.com/hultan/softtube/internal/softtube.core"
)

const (
	largeFileSize = 50000		// 50 kb is an acceptable size
	bestQuality = 100			// JPEG quality
)

var (
	config *core.Config
)

func main() {
	// Init config file
	config = new(core.Config)
	err := config.Load("main")
	if err != nil {
		fmt.Println("ERROR (Open config) : ", err.Error())
		os.Exit(1)
	}

	var thumbnailPath = config.ServerPaths.Thumbnails
	// var thumbnailPath = "/softtube/thumb/"

	// Find large files
	large, err := findLargeFiles(thumbnailPath)
	if err != nil {
		log.Fatal(err)
	}

	// Shrink files
	for _, fileName := range large {
		fullPath := fmt.Sprintf("%s/%s",thumbnailPath , fileName)

		log.Println("Shrinking file : ", fullPath)

		err = shrinkFile(fullPath)
		if err != nil {
			log.Println(err)
		}
	}
}

func shrinkFile(fullPath string) error {
	thumbnail, err := gdk.PixbufNewFromFile(fullPath)
	if err != nil {
		return err
	}

	thumbnail, err = thumbnail.ScaleSimple(160, 90, gdk.INTERP_BILINEAR)
	if err != nil {
		return err
	}

	err = thumbnail.SaveJPEG(fullPath, bestQuality)
	if err != nil {
		return err
	}

	return nil
}

func findLargeFiles(root string) ([]string, error) {
	var f []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Size() > largeFileSize {
			f = append(f, info.Name())
		}
		return nil
	})
	return f, err
}
