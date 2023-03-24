package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	Path    string
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	logFile *os.File
}

const constMaxLogFileSize = 500000

func NewStandardLogger(path string) (*Logger, error) {
	logger := &Logger{Path: path}

	logFile, err := logger.getLogFile()
	if err != nil {
		return nil, err
	}

	logger.initLogging(io.Discard, logFile, logFile, logFile)
	logger.logFile = logFile

	return logger, nil
}

func (l *Logger) Close() {
	_ = l.logFile.Close()
}

func (l *Logger) getLogFile() (*os.File, error) {
	// Get log file size
	size, _ := l.getLogFileSize(l.Path)
	// Check if the log file is too large
	if size > constMaxLogFileSize {
		// Remove old bak file
		err := os.Remove(l.Path + ".bak")
		if err != nil {
			return nil, err
		}
		// Rename log file to log file.bak
		err = os.Rename(l.Path, l.Path+".bak")
		if err != nil {
			return nil, err
		}
	}
	// Open the log file
	f, err := os.OpenFile(l.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// Handle errors
	if err != nil {
		return nil, err
	}

	// Return log object
	return f, nil
}

func (l *Logger) getLogFileSize(path string) (int64, error) {
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

func (l *Logger) initLogging(traceHandle, infoHandle, warningHandle, errorHandle io.Writer) {
	l.Trace = log.New(
		traceHandle,
		"[TRACE] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	l.Info = log.New(
		infoHandle,
		"[INFO]  ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	l.Warning = log.New(
		warningHandle,
		"[WARN]  ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	l.Error = log.New(
		errorHandle,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}
