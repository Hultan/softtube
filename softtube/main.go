package main

import (
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	playVideo()
}

// Get the duration of a youtube video
func playVideo() {
	// youtube-dl --get-duration -- '%s'
	command := "smplayer smb://192.168.1.3/softtube/test.mkv"
	cmd := exec.Command("/bin/bash", "-c", command)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
}
