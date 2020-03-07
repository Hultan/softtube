package core

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Logger : Log file
type Logger struct {
	Path string
	File *os.File
}

// NewLog : Creates and opens the log file
func NewLog(path string) Logger {
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
	l := Logger{path, nil}
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

// Close the log file
func (l Logger) Close() {
	l.File.Close()
}

// LogStart : Write beginning of log message
func (l Logger) LogStart(config *Config, application string) {
	msg := fmt.Sprintf("- %s started -", application)

	log.Println(strings.Repeat("-", len(msg)))
	log.Println(msg)
	log.Println(strings.Repeat("-", len(msg)))
	log.Println("")
	log.Println("Settings:")
	log.Println("---------")
	log.Println("CONNECTION:")
	log.Println("	Server 			: ", config.Connection.Server)
	log.Println("	Port 			: ", config.Connection.Port)
	log.Println("	Database 		: ", config.Connection.Database)
	log.Println("	Username 		: ", config.Connection.Username)
	log.Println("BACKUP:")
	log.Println("	Path			: ", config.Backup.Path)
	log.Println("UPDATE:")
	log.Println("	Log path 		: ", config.Update.Log)
	log.Println("	Youtube-dl path	: ", config.Update.YoutubeDL)
	log.Println("	Thumbnails path	: ", config.Update.Thumbnails)
	log.Println("CLIENT:")
	log.Println("	Log path 		: ", config.Client.Log)
	log.Println("	Videos path		: ", config.Client.Videos)
	log.Println("INTERVALS:")
	log.Println("	High 			: ", config.Intervals.High)
	log.Println("	Medium 			: ", config.Intervals.Medium)
	log.Println("	Low 			: ", config.Intervals.Low)
	log.Println("---------")
	log.Println("")
}

// LogFinished : Write end of log message
func (l Logger) LogFinished(application string) {
	log.Println("")
	log.Println("-------------------")
	log.Println("- update finished -")
	log.Println("-------------------")
	log.Println("")
}

// Log : Simple log function
func (l Logger) Log(text string) {
	log.Println(text)
}

// LogError : Logs an error
func (l Logger) LogError(err error) {
	log.Println(err.Error())
}

// LogFormat : Logs and formats string
func (l Logger) LogFormat(v ...interface{}) {
	log.Println(v...)
}

// Failed to open the log file
func (l Logger) failedOpenLogFile(err error) {
	fmt.Println(err)
	fmt.Println("Failed to open SoftTube log file.")
	os.Exit(0)
}

func getLogFileSize(path string) (int64, error) {
	if _, err := os.Stat(path); err == nil {
		// log file exists
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

	return 0, nil
}
