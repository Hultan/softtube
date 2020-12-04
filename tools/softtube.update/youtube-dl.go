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

	for i := 0; i < 3; i++ {
		duration, err := y.getDurationInternal(videoID)
		if y.checkDuration(videoID, duration) {
			return nil
		}
		if err != nil {
			switch i {
			case 0:
				logger.Log("Failed to get duration (" + videoID + "), trying again in 5 seconds!")
				time.Sleep(5 * time.Second)
				continue
			case 1:
				logger.Log("Failed to get duration (" + videoID + "), trying again in 30 seconds!")
				time.Sleep(30 * time.Second)
				continue
			case 2:
				logger.Log("DURATION OUTPUT (" + videoID + ") : ")
				logger.Log(duration)
				logger.Log("DURATION OUTPUT (end) ")
				// Save duration in the database
				y.updateDuration(videoID, "", logger)
				return err
			}
		}

		// Success, save duration in the database
		y.updateDuration(videoID, duration, logger)
	}

	return nil
}

func (y youtube) checkDuration(videoID, duration string) bool {
	if duration == "0" || strings.HasPrefix(duration, "ERROR: Premieres") || strings.HasPrefix(duration, "ERROR: This live event") {
		// Is it a live streaming event?
		// Save duration in the database
		y.updateDuration(videoID, "LIVE", logger)
		return true
	}
	return false
}

func (y youtube) updateDuration(videoID, duration string, logger *log.Logger) {
	// Save duration in the database
	err := db.Videos.UpdateDuration(videoID, strings.Trim(duration, " \n"))
	if err != nil {
		logger.Log("UPDATE DURATION ERROR (" + videoID + ") : ")
		logger.Log(err.Error())
		logger.Log("UPDATE DURATION ERROR (end) : ")
	}
}

// Get the duration of a youtube video
func (y youtube) getDurationInternal(videoID string) (string, error) {
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

	for i:=0;i<3;i++ {
		output, err := y.getThumbnailInternal(videoID, thumbnailPath)
		if err!=nil {
			switch i {
			case 0:
				logger.Log("Failed to download thumbnail (" + videoID + "), trying again in 5 seconds!")
				time.Sleep(5 * time.Second)
				continue
			case 1:
				logger.Log("Failed to download thumbnail (" + videoID + "), trying again in 30 seconds!")
				time.Sleep(30 * time.Second)
				continue
			case 2:
				logger.Log("THUMBNAIL OUTPUT (" + videoID + ") : ")
				logger.Log(output)
				logger.Log("THUMBNAIL OUTPUT (end) ")
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
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		logger.LogError(err)
	}
	xml := buf.String()

	err = response.Body.Close()
	if err != nil {
		logger.LogError(err)
	}

	return xml, nil
}

func (y youtube) getYoutubePath() string {
	return path.Join(config.ServerPaths.YoutubeDL, "youtube-dl")
}
