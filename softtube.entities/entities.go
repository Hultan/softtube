package entities

import "time"

// Version : Represents the SoftTube database version
type Version struct {
	Major int
}

// Log : Represents a SoftTube log entry
type Log struct {
	LogID      int
	LogMessage string
	LogType    int
}

// Subscription : Represents a YouTube subscription in SoftTube
type Subscription struct {
	ID          string
	Name        string
	Frequency   int
	LastChecked time.Time
	NextUpdate  int
}

// Video : Represents a YouTube video in SoftTube
type Video struct {
	ID         string
	ChannelID  string
	Title      string
	Added      time.Time
	Published  time.Time
	Duration   string
	Downloaded bool
	Watched    bool
}
