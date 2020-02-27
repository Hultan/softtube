//
// youtube-dl
//
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	database "github.com/hultan/softtube/softtube.database"
)

type youtube struct {
}

// Get the duration of a youtube video
func (y youtube) getDuration(config *Config, videoID string) {
	// youtube-dl --get-duration -- '%s'
	command := fmt.Sprintf(videoDurationCommand, videoID)
	cmdOutput, err := exec.Command("/bin/bash", "-c", command).Output()
	// TODO : Fix error handling
	if err != nil {
		log.Fatal(err)
	}
	// Save duration in the database
	databaseVideo := database.VideosTable{Path: config.Paths.Database}
	databaseVideo.UpdateDuration(videoID, string(cmdOutput))
	// TODO : Fix logging
}

// Get the thumbnail of a youtube video
func (y youtube) getThumbnail(config *Config, videoID string) {
	// %s/%s.jpg
	path := fmt.Sprintf(thumbnailPath, config.Paths.Thumbnails, videoID)
	// youtube-dl --write-thumbnail --skip-download --no-overwrites -o '%s' -- '%s'
	command := fmt.Sprintf(thumbnailCommand, path, videoID)
	_, err := exec.Command("/bin/bash", "-c", command).Output()
	// TODO : Fix error handling
	if err != nil {
		log.Fatal(err)
	}
	// TODO : Fix logging
}

// Get the subscription RSS to a string.
func (y youtube) getSubscriptionRSS(channelID string) (string, error) {
	url := fmt.Sprintf(subscriptionRSSURL, channelID)
	// Get the xml from the URL
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Convert the response body to a string
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	xml := buf.String()

	return xml, nil
}
