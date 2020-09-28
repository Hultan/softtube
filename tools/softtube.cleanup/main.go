package main

import (
	"fmt"
	crypt "github.com/hultan/softteam-tools/pkg/crypt"
	core "github.com/hultan/softtube/internal/softtube.core"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var (
	cutOff = time.Now().AddDate(0, 0, -14)
	config *core.Config
	db     *core.Database
)

func main() {
	// Load config file
	config = new(core.Config)
	config.Load("main")

	// Decrypt the MySQL password
	conn := config.Connection
	crypt := crypt.Crypt{}
	password, err := crypt.Decrypt(conn.Password)
	if err != nil {
		panic(err)
	}

	// Create the database object, and get all subscriptions
	db = core.New(conn.Server, conn.Port, conn.Database, conn.Username, password)
	db.OpenDatabase()
	defer db.CloseDatabase()

	cleanBackups()
	cleanThumbnails(db)
}

func cleanBackups() {
	root := "/softtube/backup"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fileName := path.Join(root, file.Name())
		info, err := os.Stat(fileName)
		if err != nil {
			// Ignore file
		}
		if !info.IsDir() {
			time := info.ModTime()
			if time.Before(cutOff) {
				fmt.Println("Removing old backup:", fileName)
				os.Remove(fileName)
			}
		}
	}
}

func cleanThumbnails(db *core.Database) {
	root := "/softtube/thumbnails"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}


	for _, file := range files {
		fileName := path.Join(root, file.Name())

		info, err := os.Stat(fileName)
		if err != nil {
			// Ignore file
		}
		if !info.IsDir() {
			time := info.ModTime()
			if time.Before(cutOff) {
				videoId := FilenameWithoutExtension(file.Name())
				if !videoIsDownloadedAndNotWatched(db, videoId) {
					os.Remove(fileName)
				}
			}
		}
	}
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func videoIsDownloadedAndNotWatched(db *core.Database, videoId string) bool {
	status, _ := db.Videos.GetStatus(videoId)
	return status == 2
}