package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hultan/crypto"
	log "github.com/hultan/softtube/internal/logger"
	core "github.com/hultan/softtube/internal/softtube.core"
	"github.com/hultan/softtube/internal/softtube.database"
)

var (
	logger          *log.Logger
	config          *core.Config
	backupCutOff    = time.Now().AddDate(0, 0, -10)
	thumbnailCutOff = time.Now().AddDate(0, 0, -7)
	db              *database.Database
)

const (
	errorOpenConfig      = 1
	errorOpenLog         = 2
	errorCleanBackup     = 3
	errorCleanThumbnails = 4
	errorDecryptPassword = 5
	errorOpenDatabase    = 6
)

func main() {
	// Init config file
	config = &core.Config{}
	err := config.Load("main")
	if err != nil {
		fmt.Println("Failed to open config file!")
		fmt.Println(err)
		os.Exit(errorOpenConfig)
	}

	// Open log file
	logger, err = log.NewStandardLogger(path.Join(config.ServerPaths.Log, config.Logs.Cleanup))
	if err != nil {
		fmt.Println("Failed to open log file!")
		fmt.Println(err)
		os.Exit(errorOpenLog)
	}
	defer logger.Close()

	// Start updating the softtube database
	logger.Info.Println()
	logger.Info.Println("----------------")
	logger.Info.Println("softtube.cleanup")
	logger.Info.Println("----------------")
	logger.Info.Println()

	// Decrypt the MySQL password
	conn := config.Connection
	password, err := crypto.Decrypt(conn.Password)
	if err != nil {
		logger.Error.Println("Failed to decrypt MySQL password!")
		logger.Error.Println(err)
		os.Exit(errorDecryptPassword)
	}

	// Create the database object, and get all subscriptions
	db = database.NewDatabase(conn.Server, conn.Port, conn.Database, conn.Username, password)
	err = db.Open()
	if err != nil {
		logger.Error.Println("Failed to open database!")
		logger.Error.Println(err)
		os.Exit(errorOpenDatabase)
	}

	defer db.Close()

	logger.Info.Println("Removing backups:")
	logger.Info.Println()

	err = cleanBackups()
	if err != nil {
		logger.Error.Println("Failed to cleanup backups!")
		logger.Error.Println(err)
		os.Exit(errorCleanBackup)
	}

	logger.Info.Println()
	logger.Info.Println("Removing thumbnails:")
	logger.Info.Println()

	err = cleanThumbnails(db)
	if err != nil {
		logger.Error.Println("Failed to cleanup thumbnails!")
		logger.Error.Println(err)
		os.Exit(errorCleanThumbnails)
	}
}

func cleanBackups() error {
	root := "/softtube/backup"
	files, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := path.Join(root, file.Name())
		info, err := os.Stat(fileName)
		if err != nil {
			// Ignore file
			logger.Warning.Printf("Failed to stat file '%s'\n", fileName)
			logger.Warning.Println(err)
			continue
		}
		if !info.IsDir() {
			modTime := info.ModTime()
			if modTime.Before(backupCutOff) {
				logger.Info.Printf("Removing old backup: '%s'\n", fileName)
				_ = os.Remove(fileName)
			}
		}
	}

	return nil
}

func cleanThumbnails(db *database.Database) error {
	root := "/softtube/thumbnails"
	files, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := path.Join(root, file.Name())

		info, err := os.Stat(fileName)
		if err != nil {
			// Ignore file
			logger.Warning.Printf("Failed to stat file '%s'\n", fileName)
			logger.Warning.Println(err)
			continue
		}
		if !info.IsDir() {
			modTime := info.ModTime()
			if modTime.Before(thumbnailCutOff) {
				videoId := FilenameWithoutExtension(file.Name())
				if canDeleteThumbnail(db, videoId) {
					logger.Info.Printf("Removing old thumbnail: '%s'\n", fileName)
					_ = os.Remove(fileName)
				}
			}
		}
	}

	return nil
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func canDeleteThumbnail(db *database.Database, videoId string) bool {
	status, err := db.Videos.GetStatus(videoId)
	if err != nil {
		logger.Warning.Printf("Failed to get status of video: '%s'\n", videoId)
		logger.Warning.Printf("Error: %s\n", err)
		logger.Warning.Println(err)
		return false
	}
	// Delete thumbnails for not downloaded videos and deleted videos
	return status == 0 || status == 4
}
