package database

import (
	"database/sql"
	"time"
)

// Download : Represents a download request
type Download struct {
	ID string
}

// Version : Represents the SoftTube database version
type Version struct {
	Major    int
	Minor    int
	Revision int
}

// Log : Represents a SoftTube log entry
type Log struct {
	ID      int
	Message string
	Type    int
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
	ID               string
	SubscriptionID   string
	SubscriptionName string
	Title            string
	Added            time.Time
	Published        time.Time
	Duration         sql.NullString
	Status           int
	Saved            bool
}
