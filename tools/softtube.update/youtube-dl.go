//
// youtube-dl
//
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/hultan/softteam/log"
)

type youtube struct {
}

// Get the duration of a youtube video
func (y youtube) getDuration(videoID string, logger *log.Logger) error {

	output, err := y.getDurationInternal(videoID,logger)
	if err!=nil {
		logger.Log("Failed to get duration, trying again in 5 seconds!")
		time.Sleep(5 * time.Second)
		output, err = y.getDurationInternal(videoID,logger)
		if err!=nil {
			logger.Log("Failed to get duration, trying again in 30 seconds!")
			time.Sleep(30 * time.Second)
			output, err = y.getDurationInternal(videoID, logger)
			if err != nil {
				logger.Log("DURATION FAILED : ")
				logger.Log("DURATION OUTPUT (start) : ")
				logger.Log(output)
				logger.Log("DURATION OUTPUT (end) ")
				logger.Log("DURATION ERROR (start) : ")
				logger.LogError(err)
				logger.Log("DURATION ERROR (end) ")
				return err
			}
		}
	}

	// Is it a live streaming event?
	duration := string(output)
	if strings.HasPrefix(duration, "ERROR: Premieres") {
		duration = "LIVE"
	}

	// Save duration in the database
	err= db.Videos.UpdateDuration(videoID, strings.Trim(duration, " \n"))
	if err!=nil {
		logger.Log(err.Error())
	}
	return nil
}

// Get the duration of a youtube video
func (y youtube) getDurationInternal(videoID string, logger *log.Logger) (string, error) {
	// youtube-dl --get-duration -- '%s'
	command := fmt.Sprintf(constVideoDurationCommand, y.getYoutubePath(), videoID)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), nil
}

// Get the thumbnail of a youtube video
func (y youtube) getThumbnail(videoID, thumbnailPath string, logger *log.Logger) error {

	output, err := y.getThumbnailInternal(videoID, thumbnailPath)
	if err!=nil {
		logger.Log("Failed to download thumbnail, trying again in 5 seconds!")
		time.Sleep(5 * time.Second)
		output, err = y.getThumbnailInternal(videoID, thumbnailPath)
		if err!=nil {
			logger.Log("Failed to download thumbnail, trying again in 30 seconds!")
			time.Sleep(30 * time.Second)
			output, err = y.getThumbnailInternal(videoID, thumbnailPath)
			if err != nil {
				logger.Log("THUMBNAIL OUTPUT (start) : ")
				logger.Log(output)
				logger.Log("THUMBNAIL OUTPUT (end) ")
				logger.Log("THUMBNAIL ERROR (start) : ")
				logger.LogError(err)
				logger.Log("THUMBNAIL ERROR (end) ")
				return err
			}
		}
	}
	return nil
}

// Get the thumbnail of a youtube video
func (y youtube) getThumbnailInternal(videoID, thumbnailPath string) (string, error) {
	// %s/%s.jpg
	thumbPath := fmt.Sprintf(constThumbnailLocation, thumbnailPath, videoID)

	// Don't download thumbnail if it already exists
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		// youtube-dl --write-thumbnail --skip-download --no-overwrites -o '%s' -- '%s'
		command := fmt.Sprintf(constThumbnailCommand, y.getYoutubePath(), thumbPath, videoID)
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return string(output), err
		}
	}
	return "", nil
}

// Get the subscription RSS to a string.
func (y youtube) getSubscriptionRSS(channelID string) (string, error) {
	url := fmt.Sprintf(constSubscriptionRSSURL, channelID)
	// Get the xml from the URL
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}


	// Convert the response body to a string
	buf := new(bytes.Buffer)
	_,err = buf.ReadFrom(response.Body)
	if err!=nil {
		logger.LogError(err)
	}
	xml := buf.String()

	err = response.Body.Close()
	if err!=nil {
		logger.LogError(err)
	}

	return xml, nil
}

func (y youtube) getYoutubePath() string {
	return path.Join(config.ServerPaths.YoutubeDL, "youtube-dl")
}
