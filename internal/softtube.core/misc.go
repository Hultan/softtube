package core

import (
	"fmt"
	"os/user"
)

// Get the current users home directory
func getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
