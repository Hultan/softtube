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
	"regexp"
	"strings"
	"time"
)

type youtube struct {
}

// Get the duration of a YouTube video
func (y youtube) getDuration(videoID string) error {
	backoff := []time.Duration{5 * time.Second, 30 * time.Second}

	for i := 0; i <= len(backoff); i++ {
		duration, err := y.getDurationInternal(videoID)

		if err == nil && duration != "" {
			y.updateDuration(videoID, duration)
			return nil
		}

		logger.Warning.Printf("Attempt %d: Failed to get duration for video %s. Error: %v, Output: %s\n",
			i+1, videoID, err, duration)

		// Special case: if we know it's unavailable, don't retry
		if strings.Contains(duration, "Video unavailable") {
			return err
		}

		// If this was the last attempt, save empty duration and return
		if i == len(backoff) {
			y.updateDuration(videoID, "")
			return err
		}

		// Otherwise, wait and retry
		time.Sleep(backoff[i])
	}

	return nil // Should not reach here
}

// Get the duration of a YouTube video
func (y youtube) getDurationInternal(videoId string) (string, error) {
	if videoId == "" {
		return "", fmt.Errorf("getDurationInternal: videoID cannot be empty")
	}

	command := fmt.Sprintf(constVideoDurationCommand, y.getYoutubePath(), videoId)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.CombinedOutput()
	// We check if it is a live/error/member event before checking the
	// error here, because checking the duration of a live/error/member
	// event DOES fail and returns an error.
	durationString := string(output)
	if isLive(durationString) {
		return "LIVE", nil
	} else if isError(durationString) {
		return "ERROR", nil
	} else if isMember(durationString) {
		return "MEMBER", nil
	}

	// Trying to avoid ERROR: fragment 1 not found, unable to continue
	// Problem on YouTube:s or yt-dlp:s side?
	// Note that the fragment error string is in the durationString (output)
	// and not in the error (err).
	if err != nil && !strings.Contains(durationString, "fragment") {
		return durationString, fmt.Errorf("getDurationInternal: get duration failed: %s", err)
	}

	// This regex matches:
	// - hh:mm:ss
	// - mm:ss
	// - ss (only if it appears at the end, and optionally surrounded by parentheses or whitespace)
	re := regexp.MustCompile(`(?m)(\d{1,2}:)?\d{1,2}(:\d{2})?$`)

	if match := re.FindString(durationString); match != "" {
		return match, nil
	}
	return "", fmt.Errorf("getDurationInternal: failed to extract duration from '%s'", durationString)
}

func (y youtube) updateDuration(videoID, duration string) {
	// Save duration in the database
	err := db.Videos.UpdateDuration(videoID, strings.Trim(duration, " \n"))
	if err != nil {
		logger.Error.Printf("Failed to update duration (%s)!\n", videoID)
		logger.Error.Println(err)
	}
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
		outputString := string(output)

		if err != nil && !strings.Contains(outputString, "fragment") {
			return outputString, err
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
	// Get the XML from the URL
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

func isLive(duration string) bool {
	upper := strings.ToUpper(duration)
	if strings.Contains(upper, "PREMIERES") {
		return true
	}
	if strings.Contains(upper, "LIVE EVENT") {
		return true
	}
	if upper == "" {
		return true
	}

	return false
}

func isError(duration string) bool {
	if strings.Contains(duration, "uploader has not") {
		return true
	}
	if strings.Contains(duration, "This video has been removed by the uploader") {
		return true
	}

	return false
}

func isMember(duration string) bool {
	if strings.Contains(duration, "Join this channel") {
		return true
	}

	return false
}
