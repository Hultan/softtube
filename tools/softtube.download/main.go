package main

import (
	"fmt"
	"github.com/hultan/softtube/internal/softtube.database"
	"os"
	"os/exec"
	"path"
	"sync"

	core "github.com/hultan/softtube/internal/softtube.core"
)

var (
	logger *core.Logger
	config *core.Config
	db     *database.Database
)

func main() {
	// Load config file
	config = new(core.Config)
	err := config.Load("main")
	if err != nil {
		fmt.Println("ERROR (Open config) : ", err.Error())
		os.Exit(errorOpenConfig)
	}

	// Setup logging
	logger = core.NewLog(path.Join(config.ServerPaths.Log, config.Logs.Download))
	defer logger.Close()

	// Start updating the softtube database
	logger.LogStart("softtube download")
	defer logger.LogFinished("softtube download")

	// Decrypt the MySQL password
	conn := config.Connection
	crypto := core.Crypt{}
	password, err := crypto.Decrypt(conn.Password)
	if err != nil {
		logger.Log("Failed to decrypt MySQL password!")
		logger.LogError(err)
		panic(err)
	}

	// Create the database object, and get all subscriptions
	db = database.New(conn.Server, conn.Port, conn.Database, conn.Username, password)
	err = db.OpenDatabase()
	if err != nil {
		logger.Log("ERROR (Open database)")
		logger.LogError(err)
		os.Exit(errorOpenDatabase)
	}
	defer db.CloseDatabase()
	downloads, err := db.Download.GetAll()
	if err != nil {
		logger.Log(err.Error())
		panic(err)
	}

	var waitGroup sync.WaitGroup

	for i := 0; i < len(downloads); i++ {
		waitGroup.Add(1)
		go downloadVideo(downloads[i].ID, &waitGroup)
	}
	waitGroup.Wait()
}

// Download a youtube video
func downloadVideo(videoID string, wait *sync.WaitGroup) error {
	// Set video status as downloading
	err := db.Videos.UpdateStatus(videoID, constStatusDownloading)
	if err != nil {
		logger.Log("Failed to set video status to downloading before download!")
		logger.LogError(err)
		wait.Done()
		return err
	}

	// Set the video as downloaded in database
	// Delete it from the table download immediately to
	// avoid multiple download attempts (that can cause
	// crashes)
	err = db.Download.SetAsDownloaded(videoID)
	if err != nil {
		logger.Log("Failed to delete video from table download after download!")
		logger.LogError(err)
		wait.Done()
		return err
	}

	// Download the video
	command := fmt.Sprintf("%s -f 'bestvideo[height<=1080]+bestaudio/best[height<=1080]' --no-overwrites -o '%s/%%(id)s.%%(ext)s' -- '%s'", getYoutubePath(), config.ServerPaths.Videos, videoID)
	//command := fmt.Sprintf("%s -f best --no-overwrites -o '%s/%%(id)s.%%(ext)s' -- '%s'", getYoutubePath(), config.ServerPaths.Videos, videoID)
	fmt.Println(command)
	cmd := exec.Command("/bin/bash", "-c", command)
	// Wait for the command to be executed (video to be downloaded)
	err = cmd.Run()
	if err != nil {
		logger.Log("Failed to download video!")
		msg := fmt.Sprintf("Command : %s", command)
		logger.Log(msg)
		logger.LogError(err)
		wait.Done()
		return err
	}

	// Set video status as downloaded
	err = db.Videos.UpdateStatus(videoID, constStatusDownloaded)
	if err != nil {
		logger.Log("Failed to set video status to downloaded after download!")
		logger.LogError(err)
		wait.Done()
		return err
	}

	wait.Done()
	return nil
}

func getYoutubePath() string {
	return path.Join(config.ServerPaths.YoutubeDL, "youtube-dl")
}
