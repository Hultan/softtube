package core

import (
	"database/sql"
	"time"
)

// Version : Represents the SoftTube database version
type Version struct {
	Major    int
	Minor    int
	Revision int
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
	LastChecked sql.NullTime
	NextUpdate  sql.NullInt32
}

// Video : Represents a YouTube video in SoftTube
type Video struct {
	ID             string
	SubscriptionID string
	Title          string
	Added          time.Time
	Published      time.Time
	Duration       sql.NullString
	Status         int
}
