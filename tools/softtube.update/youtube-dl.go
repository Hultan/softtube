// youtube-dl (yt-dlp)
package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

type youtube struct {
}

// Get the duration of a YouTube video
func (y youtube) getDuration(videoID string) error {

	for i := 0; i < 3; i++ {
		duration, err := y.getDurationInternal(videoID)
		if err != nil || duration == "" {
			switch i {
			case 0:
				logger.Warning.Printf("Failed to get duration (%s) : Error = (%s)\n", videoID, err)
				logger.Warning.Printf("Failed to get duration (%s) : Output = (%s)\n", videoID, duration)
				if strings.Contains(duration, "Video unavailable") {
					return err
				}
				time.Sleep(5 * time.Second)
				continue
			case 1:
				time.Sleep(30 * time.Second)
				continue
			case 2:
				logger.Warning.Printf("Failed to get duration (%s)\n", videoID)
				// Save duration in the database
				y.updateDuration(videoID, "")
				return err
			}
		}

		// Success, save duration in the database
		y.updateDuration(videoID, duration)
	}

	return nil
}

func (y youtube) updateDuration(videoID, duration string) {
	// Save duration in the database
	err := db.Videos.UpdateDuration(videoID, strings.Trim(duration, " \n"))
	if err != nil {
		logger.Error.Printf("Failed to update duration (%s)!\n", videoID)
		logger.Error.Println(err)
	}
}

func (y youtube) isLiveEvent(duration string) bool {
	if duration == "0" {
		return true
	}
	if strings.HasPrefix(duration, "ERROR: Premieres") {
		return true
	}
	if strings.HasPrefix(duration, "ERROR: This live event") {
		return true
	}
	return false
}

// Get the duration of a YouTube video
func (y youtube) getDurationInternal(videoId string) (string, error) {
	if videoId == "" {
		return "", errors.New("videoId cannot be null")
	}

	command := fmt.Sprintf(constVideoDurationCommand, y.getYoutubePath(), videoId)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.CombinedOutput()
	// We check if it is a live event before checking the
	// error here, because checking duration of a live event
	// DOES fail and returns an error.
	if y.isLiveEvent(string(output)) {
		return "LIVE", nil
	}
	if err != nil {
		return string(output), err
	}
	return y.getLastRow(string(output)), nil
}

func (y youtube) getLastRow(text string) string {
	var nonEmptyRows []string

	rows := strings.Split(text, "\n")
	for _, row := range rows {
		row = strings.Trim(row, " \t\n")
		if row != "" {
			nonEmptyRows = append(nonEmptyRows, row)
		}
	}

	if len(nonEmptyRows) == 0 {
		return ""
	}

	return nonEmptyRows[len(nonEmptyRows)-1]
}

// Get the thumbnail of a YouTube video
func (y youtube) getThumbnail(videoId string) error {
	for i := 0; i < 3; i++ {
		output, err := y.getThumbnailInternal(videoId)
		if err != nil {
			switch i {
			case 0:
				logger.Warning.Printf("Failed to download thumbnail (%s) : Error = (%s)\n", videoId, err)
				logger.Warning.Printf("Failed to download thumbnail (%s) : Output = (%s)\n", videoId, output)
				if strings.Contains(output, "Video unavailable") {
					return err
				}
				time.Sleep(5 * time.Second)
				continue
			case 1:
				time.Sleep(30 * time.Second)
				continue
			case 2:
				logger.Warning.Printf("Failed to download thumbnail (%s)\n", videoId)
				return err
			}
		}
	}
	return nil
}

// Get the thumbnail of a YouTube video
func (y youtube) getThumbnailInternal(videoId string) (string, error) {
	if videoId == "" {
		return "", errors.New("videoId cannot be null")
	}

	thumbPath := y.getThumbnailPath(videoId)

	// Don't download thumbnail if it already exists
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		command := fmt.Sprintf(constThumbnailCommand, y.getYoutubePath(), thumbPath, videoId)
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return string(output), err
		}
	}

	return "", nil
}

// Get the subscription RSS to a string.
func (y youtube) getSubscriptionRSS(channelId string) (string, error) {
	if channelId == "" {
		return "", errors.New("channelId cannot be null")
	}

	url := y.getRSSFeedURL(channelId)
	// Get the xml from the URL
	response, err := http.Get(url)
	if err != nil {
		logger.Error.Printf("Failed get url = '%s'!\n", url)
		logger.Error.Println(err)
		return "", err
	}

	// Convert the response body to a string
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		logger.Error.Println("Failed to read from response!")
		logger.Error.Println(err)
		return "", err
	}
	xml := buf.String()

	// Close the response object
	err = response.Body.Close()
	if err != nil {
		logger.Error.Println("Failed to close response!")
		logger.Error.Println(err)
		return "", err
	}

	return xml, nil
}

func (y youtube) getYoutubePath() string {
	return path.Join(config.ServerPaths.YoutubeDL, "yt-dlp")
}

func (y youtube) getRSSFeedURL(channelId string) string {
	return fmt.Sprintf(constSubscriptionRSSURL, channelId)
}

func (y youtube) getThumbnailPath(videoId string) string {
	return fmt.Sprintf(constThumbnailLocation, config.ServerPaths.Thumbnails, videoId)
}
