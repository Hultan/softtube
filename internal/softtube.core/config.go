package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	ServerPaths struct {
		Log        string `json:"log"`
		Backup     string `json:"backup"`
		YoutubeDL  string `json:"youtube-dl"`
		Thumbnails string `json:"thumbnails"`
		Videos     string `json:"videos"`
	} `json:"server-paths"`
	ClientPaths struct {
		Log        string `json:"log"`
		Thumbnails string `json:"thumbnails"`
		Videos     string `json:"videos"`
	} `json:"client-paths"`
	Logs struct {
		Backup   string `json:"backup"`
		Update   string `json:"update"`
		Download string `json:"download"`
		SoftTube string `json:"softtube"`
		Cleanup  string `json:"cleanup"`
		Shrink   string `json:"shrink"`
	} `json:"logs"`
	Intervals struct {
		High   int `json:"high"`
		Medium int `json:"medium"`
		Low    int `json:"low"`
	} `json:"intervals"`
}

// Load the SoftTube configuration file
func (config *Config) Load(mode string) error {
	// Get the path to the config file
	configPath := config.getConfigPath(mode)

	// Make sure the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		errorMessage := fmt.Sprintf("settings file is missing (%s)", constConfigPath)
		return errors.New(errorMessage)
	}

	// Open the config file
	configFile, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to open config")
	}

	// Parse the JSON document
	jsonParser := json.NewDecoder(configFile)
	_ = jsonParser.Decode(&config)

	_ = configFile.Close()

	return nil
}

// Save the SoftTube configuration file
func (config *Config) Save(mode string) {
	// Get the path to the config file
	configPath := config.getConfigPath(mode)

	// Open the config file
	configFile, err := os.OpenFile(configPath, os.O_TRUNC|os.O_WRONLY, 0644)

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Create JSON from the config object
	data, err := json.MarshalIndent(config, "", "\t")

	// Handle errors
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Write the data
	_, _ = configFile.Write(data)

	_ = configFile.Close()
}

// Get the path to the config file
// Mode = "test" returns test config path
// otherwise returns the normal config path
func (config *Config) getConfigPath(mode string) string {
	home := getHomeDirectory()

	var configPath string
	if strings.ToLower(mode) == "test" {
		configPath = constConfigPathTest
	} else {
		configPath = constConfigPath
	}

	return path.Join(home, configPath)
}
