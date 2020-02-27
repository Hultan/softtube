package main

import (
	"fmt"
	"log"
	"os"
)

// Log : Log file
type Log struct {
	Path string
	File *os.File
}

func createAndOpenLog(path string) Log {
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

func (l Log) logStart(config *Config) {
	log.Println("------------------")
	log.Println("- update started -")
	log.Println("------------------")
	log.Println("")
	log.Println("Settings:")
	log.Println("---------")
	log.Println("Log path : " + config.Paths.Log)
	log.Println("Database path : " + config.Paths.Database)
	log.Println("Videos path : " + config.Paths.Videos)
	log.Println("Thumbnails path : " + config.Paths.Thumbnails)
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
