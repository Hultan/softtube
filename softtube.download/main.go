package main

import (
	"fmt"
	"os/exec"
	"path"
	"sync"

	core "github.com/hultan/softtube/softtube.core"
)

var (
	logger core.Logger
	config *core.Config
	db     core.Database
)

func main() {
	// Load config file
	config = new(core.Config)
	config.Load("main")

	// Setup logging
	logger = core.NewLog(path.Join(config.ServerPaths.Log, config.Logs.Download))
	defer logger.Close()

	// Start updating the softtube database
	logger.LogStart(config, "softtube download")
	defer logger.LogFinished("softtube download")

	// Decrypt the MySQL password
	conn := config.Connection
	crypt := core.Crypt{}
	password, err := crypt.Decrypt(conn.Password)
	if err != nil {
		logger.Log("Failed to decrypt MySQL password!")
		logger.LogError(err)
		panic(err)
	}

	// Create the database object, and get all subscriptions
	db = core.New(conn.Server, conn.Port, conn.Database, conn.Username, password)
	db.OpenDatabase()
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
	command := fmt.Sprintf("%s --no-overwrites -o '%s/%%(id)s.%%(ext)s' -- '%s'", getYoutubePath(), config.ServerPaths.Videos, videoID)
	fmt.Println(command)
	cmd := exec.Command("/bin/bash", "-c", command)
	// Wait for the command to be executed (video to be downloaded)
	err := cmd.Run()
	if err != nil {
		logger.Log("Failed to download video!")
		msg := fmt.Sprintf("Command : %s", command)
		logger.Log(msg)
		logger.LogError(err)
		wait.Done()
		return err
	}
	// Set the video as downloaded in database
	err = db.Download.SetAsDownloaded(videoID)
	if err != nil {
		logger.Log("Failed to set as downloaded after download!")
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
