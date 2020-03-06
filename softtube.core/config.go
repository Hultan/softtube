package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
)

// Config for the SoftTube program
type Config struct {
	Connection struct {
		Server   string `json:"server"`
		Port     int    `json:"port"`
		Database string `json:"database"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"connection"`
	Paths struct {
		Backup     string `json:"backup"`
		Log        string `json:"log"`
		YoutubeDL  string `json:"youtube-dl"`
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
func (config *Config) Load(mode string) error {
	// Get the path to the config file
	path := getConfigPath(mode)

	// Make sure the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		errorMessage := fmt.Sprintf("settings file is missing (%s)", constConfigPath)
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
func (config *Config) Save(mode string) {
	// Get the path to the config file
	path := getConfigPath(mode)

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

// Get path to the config file
// Mode = "test" returns test config path
// otherwise returns normal config path
func getConfigPath(mode string) string {
	home := getHomeDirectory()

	var configPath string
	if strings.ToLower(mode) == "test" {
		configPath = constConfigPathTest
	} else {
		configPath = constConfigPath
	}

	return path.Join(home, configPath)
}

// Get current users home directory
func getHomeDirectory() string {
	u, err := user.Current()
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get user home directory : %s", err)
		panic(errorMessage)
	}
	return u.HomeDir
}
