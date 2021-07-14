package main

import (
	"fmt"
	core "github.com/hultan/softtube/internal/softtube.core"
	"github.com/hultan/softtube/internal/softtube.database"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var (
	cutOff = time.Now().AddDate(0, 0, -30)
	config *core.Config
	db     *database.Database
)

func main() {
	//// Load config file
	//config = new(core.Config)
	//err := config.Load("main")
	//if err!=nil {
	//	fmt.Println(err.Error())
	//	os.Exit(1)
	//}
	//
	//// Decrypt the MySQL password
	//conn := config.Connection
	//crypt := core.Crypt{}
	//password, err := crypt.Decrypt(conn.Password)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Create the database object, and get all subscriptions
	//db = database.New(conn.Server, conn.Port, conn.Database, conn.Username, password)
	//err = db.OpenDatabase()
	//if err!=nil {
	//	fmt.Println(err.Error())
	//	os.Exit(1)
	//}
	//
	//defer db.CloseDatabase()

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
			continue
		}
		if !info.IsDir() {
			modTime := info.ModTime()
			if modTime.Before(cutOff) {
				fmt.Println("Removing old backup:", fileName)
				_ = os.Remove(fileName)
			}
		}
	}
}

func cleanThumbnails(db *database.Database) {
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
			continue
		}
		if !info.IsDir() {
			modTime := info.ModTime()
			if modTime.Before(cutOff) {
				videoId := FilenameWithoutExtension(file.Name())
				if !videoIsDownloadedAndNotWatched(db, videoId) {
					_ = os.Remove(fileName)
				}
			}
		}
	}
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func videoIsDownloadedAndNotWatched(db *database.Database, videoId string) bool {
	status, _ := db.Videos.GetStatus(videoId)
	return status == 2
}