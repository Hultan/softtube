package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/hultan/crypto"
	log "github.com/hultan/softtube/internal/logger"
	"github.com/hultan/softtube/internal/softtube.database"

	core "github.com/hultan/softtube/internal/softtube.core"
)

const constDownloadCommand = "%s -f 'bestvideo[height<=1080]+bestaudio/best[height<=1080]' --no-overwrites -o '%s/%%(id)s.%%(ext)s' -- '%s'"

var (
	logger *log.Logger
	config *core.Config
	db     *database.Database
)

func main() {
	// Init config file
	config = new(core.Config)
	err := config.Load("main")
	if err != nil {
		fmt.Println("Failed to open config file!")
		fmt.Println(err)
		os.Exit(errorOpenConfig)
	}

	// Setup logging
	logger, err = log.NewStandardLogger(path.Join(config.ServerPaths.Log, config.Logs.Download))
	if err != nil {
		fmt.Println("Failed to open log file!")
		fmt.Println(err)
		os.Exit(errorOpenLogFile)
	}
	defer logger.Close()

	// Start updating the softtube database
	logger.Info.Println()
	logger.Info.Println("-----------------")
	logger.Info.Println("softtube.download")
	logger.Info.Println("-----------------")
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
	downloads, err := db.Download.GetAll()
	if err != nil {
		logger.Error.Println("Failed to get all downloads!")
		logger.Error.Println(err)
		os.Exit(errorDownload)
	}

	var waitGroup sync.WaitGroup

	for i := 0; i < len(downloads); i++ {
		waitGroup.Add(1)
		go downloadVideo(downloads[i].ID, &waitGroup)
	}
	waitGroup.Wait()
}

// Download a YouTube video
func downloadVideo(videoID string, wait *sync.WaitGroup) {
	defer wait.Done()

	// Set video status as downloading
	err := db.Videos.UpdateStatus(videoID, constStatusDownloading)
	if err != nil {
		logger.Error.Println("Failed to set video status to downloading before download!")
		logger.Error.Println(err)
		return
	}

	// Set the video as downloaded in database
	// Delete it from the table download immediately to
	// avoid multiple download attempts (that can cause
	// crashes)
	err = db.Download.Delete(videoID)
	if err != nil {
		logger.Error.Println("Failed to delete video from table download after download!")
		logger.Error.Println(err)
		return
	}

	// Create download command
	command := fmt.Sprintf(constDownloadCommand, getYoutubePath(), config.ServerPaths.Videos, videoID)
	cmd := exec.Command("/bin/bash", "-c", command)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	// Wait for the command to be executed (video to be downloaded)
	err = cmd.Run()
	if err != nil && !strings.Contains(stderr.String(), "fragment") {
		logger.Error.Println("Failed to download video!")
		msg := fmt.Sprintf("Command : %s", command)
		logger.Error.Println(msg)
		logger.Error.Println(err)
		logger.Error.Println("STDERR: \n" + stderr.String())
		return
	}

	// Set video status as downloaded
	err = db.Videos.UpdateStatus(videoID, constStatusDownloaded)
	if err != nil {
		logger.Error.Println("Failed to set video status to downloaded after download!")
		logger.Error.Println(err)
		return
	}

	return
}

func getYoutubePath() string {
	return path.Join(config.ServerPaths.YoutubeDL, "yt-dlp")
}
