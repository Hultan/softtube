package main

import (
	"fmt"
	"log"
	"os"

	core "github.com/hultan/softtube/softtube.core"
)

// Log : Log file
type Log struct {
	Path string
	File *os.File
}

func getLogFileSize(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	fInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return fInfo.Size(), nil

}
func createAndOpenLog(path string) Log {
	// Get log file size
	size, _ := getLogFileSize(path)
	// Check if the log file is too large
	if size > constMaxLogFileSize {
		// Remove old bak file
		os.Remove(path + ".bak")
		// Rename log file to log file.bak
		os.Rename(path, path+".bak")
	}
	// Create a new log object
	l := Log{path, nil}
	// Open the log file
	f, err := os.OpenFile(l.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// Handle errors
	if err != nil {
		l.failedOpenLogFile(err)
	}
	// Set file
	l.File = f
	// Set log output
	log.SetOutput(f)
	// Return log object
	return l
}

// Failed to open the log file
func (l Log) failedOpenLogFile(err error) {
	fmt.Println(err)
	fmt.Println("Failed to open SoftTube log file.")
	os.Exit(0)
}

// Close the log file
func (l Log) close() {
	l.File.Close()
}

func (l Log) logStart(config *core.Config) {
	log.Println("------------------")
	log.Println("- update started -")
	log.Println("------------------")
	log.Println("")
	log.Println("Settings:")
	log.Println("---------")
	log.Println("CONNECTION:")
	log.Println("	Server 			: ", config.Connection.Server)
	log.Println("	Port 			: ", config.Connection.Port)
	log.Println("	Database 		: ", config.Connection.Database)
	log.Println("	Username 		: ", config.Connection.Username)
	log.Println("PATHS:")
	log.Println("	Backup path		: ", config.Paths.Backup)
	log.Println("	Log path 		: ", config.Paths.Log)
	log.Println("	Youtube-dl path		: ", config.Paths.YoutubeDL)
	log.Println("	Videos path		: ", config.Paths.Videos)
	log.Println("	Thumbnails path		: ", config.Paths.Thumbnails)
	log.Println("INTERVALS:")
	log.Println("	High 			: ", config.Intervals.High)
	log.Println("	Medium 			: ", config.Intervals.Medium)
	log.Println("	Low 			: ", config.Intervals.Low)
	log.Println("---------")
	log.Println("")
}

func (l Log) logFinished() {
	log.Println("")
	log.Println("-------------------")
	log.Println("- update finished -")
	log.Println("-------------------")
	log.Println("")
}

func (l Log) log(text string) {
	log.Println(text)
}

func (l Log) logError(err error) {
	log.Println(err.Error())
}

func (l Log) logFormat(v ...interface{}) {
	log.Println(v...)
}
