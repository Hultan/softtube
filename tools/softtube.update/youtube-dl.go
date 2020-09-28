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

	"github.com/hultan/softteam/log"
)

type youtube struct {
}

// Get the duration of a youtube video
func (y youtube) getDuration(videoID string, logger *log.Logger) error {
	// youtube-dl --get-duration -- '%s'
	command := fmt.Sprintf(constVideoDurationCommand, y.getYoutubePath(), videoID)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("Command : %s", command)
		logger.Log(msg)
		msg = fmt.Sprintf("Output : %s", output)
		logger.Log(msg)
		logger.LogError(err)
		return err
	}
	// Save duration in the database
	err= db.Videos.UpdateDuration(videoID, strings.Trim(string(output), " \n"))
	if err!=nil {
		logger.Log(err.Error())
	}
	return nil
}

// Get the thumbnail of a youtube video
func (y youtube) getThumbnail(videoID, thumbnailPath string, logger *log.Logger) error {
	// %s/%s.jpg
	thumbPath := fmt.Sprintf(constThumbnailLocation, thumbnailPath, videoID)

	// Don't download thumbnail if it already exists
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		// youtube-dl --write-thumbnail --skip-download --no-overwrites -o '%s' -- '%s'
		command := fmt.Sprintf(constThumbnailCommand, y.getYoutubePath(), thumbPath, videoID)
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			msg := fmt.Sprintf("Command : %s", command)
			logger.Log(msg)
			msg = fmt.Sprintf("Output : %s", output)
			logger.Log(msg)
			logger.LogError(err)
			return err
		}
	}
	return nil
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
