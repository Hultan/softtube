package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
)

// Config for the SoftTube program
type Config struct {
	Paths struct {
		Log        string `json:"log"`
		Database   string `json:"database"`
		Thumbnails string `json:"thumbnails"`
		Videos     string `json:"videos"`
	} `json:"paths"`
	Intervals struct {
		High   int `json:"high"`
		Medium int `json:"medium"`
		Low    int `json:"low"`
	} `json:"intervals"`
}

// Load : Loads a SoftTube configuration file
func (config *Config) Load() error {
	// Get the path to the config file
	path := getConfigPath()

	// Make sure the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errorMessage := fmt.Sprintf("settings file is missing (%s)", configPath)
		return errors.New(errorMessage)
	}

	// Open config file
	configFile, err := os.Open(path)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()

	// Parse the JSON document
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return nil
}

// Save : Saves a SoftTube configuration file
func (config *Config) Save() {
	// Get the path to the config file
	path := getConfigPath()

	// Open config file
	configFile, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0644)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer configFile.Close()

	// Create JSON from config object
	data, err := json.MarshalIndent(config, "", "\t")

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Write the data
	configFile.Write(data)
}

func getConfigPath() string {
	home := getHomeDirectory()
	return path.Join(home, configPath)
}

func getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
