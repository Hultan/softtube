package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"

	"github.com/hultan/softteam/framework"
	core "github.com/hultan/softtube/internal/softtube.core"
)

const (
	largeFileSize = 50000 // 50 kb is an acceptable size
	bestQuality   = 100   // JPEG quality
)

var (
	logger *framework.Logger
	config *core.Config
)

func main() {
	var err error

	fw := framework.NewFramework()
	logger, err = fw.Log.NewStandardLogger("/softtube/log/softtube.shrink.log")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to open log file : ", err)
		os.Exit(1)
	}

	logger.Info.Println("")
	logger.Info.Println("----------------------")
	logger.Info.Println("Softtube.shrink start!")
	logger.Info.Println("----------------------")
	logger.Info.Println("")

	// Init config file
	config = &core.Config{}
	err = config.Load("main")
	if err != nil {
		logger.Error.Println("Failed to load the SoftTube config : ")
		logger.Error.Println(err)
		logger.Error.Println("Exiting...")
		os.Exit(1)
	}

	var thumbnailPath = config.ServerPaths.Thumbnails

	// Find large files
	large, err := findLargeFiles(thumbnailPath)
	if err != nil {
		logger.Error.Println("Failed to find large files : ")
		logger.Error.Println(err)
		logger.Error.Println("Exiting...")
		os.Exit(1)
	}

	logger.Info.Printf("Found %d large files.\n", len(large))
	logger.Info.Println("")

	// Shrink files
	for _, fileName := range large {
		fullPath := path.Join(thumbnailPath, fileName)

		// Log successful shrinks, errors are logged in shrinkFile()
		err = shrinkFile(fullPath)
		if err == nil {
			logger.Info.Println("Successfully shrunk file : ", fullPath)
		}
	}
}

func shrinkFile(fullPath string) error {
	thumbnail, err := gdk.PixbufNewFromFile(fullPath)
	if err != nil {
		logger.Error.Println("Failed to load file : ", fullPath)
		return err
	}

	thumbnail, err = thumbnail.ScaleSimple(160, 90, gdk.INTERP_BILINEAR)
	if err != nil {
		logger.Error.Println("Failed to scale file : ", fullPath)
		return err
	}

	err = thumbnail.SaveJPEG(fullPath, bestQuality)
	if err != nil {
		logger.Error.Println("Failed to save file : ", fullPath)
		return err
	}

	return nil
}

func findLargeFiles(root string) ([]string, error) {
	var f []string
	err := filepath.Walk(
		root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && info.Size() > largeFileSize {
				f = append(f, info.Name())
			}
			return nil
		},
	)
	return f, err
}
