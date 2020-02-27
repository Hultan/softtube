package database

import (
	"database/sql"
	"errors"
	"time"

	entities "github.com/hultan/softtube/softtube.entities"
)

// VideosTable : VideosTable in the SoftTube database
type VideosTable struct {
	Path string
}

// VideosRecord : A single video in the SubscriptionTable
type VideosRecord struct {
	Entity entities.Video
}

const sqlStatementVideoExists = "SELECT EXISTS(SELECT 1 FROM Videos WHERE video_id=?1);"
const sqlStatementInsertVideo = `INSERT OR IGNORE INTO Videos (video_id, channel_id, title, duration, published, added, downloaded, watched) 
								VALUES (?1, ?2, ?3, ?4, ?5, ?6, 0, 0);`
const sqlStatementUpdateDuration = "UPDATE Videos SET duration=?1 WHERE video_id=?2"

// Exists : Does a video already exist in the database?
func (v VideosTable) Exists(videoID string) (bool, error) {
	// Open database
	connectionString := getConnectionString(v.Path)
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return false, err
	}
	defer db.Close()

	// Execute select query
	rows, err := db.Query(sqlStatementVideoExists, videoID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		var result bool
		err = rows.Scan(&result)
		if err != nil {
			return false, err
		}
		return result, nil
	}

	return false, errors.New("unknown failure")
}

// Insert : Insert a new video into the database
func (v VideosTable) Insert(videoID string, channelID string, title string, duration string, published time.Time) error {
	// Open database
	connectionString := getConnectionString(v.Path)
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	now := time.Now().UTC().Format(dateLayout) // Added

	// Execute insert statement
	_, err = db.Exec(sqlStatementInsertVideo, videoID, channelID, title, duration, published, now)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDuration : Update duration for a video
func (v VideosTable) UpdateDuration(videoID string, duration string) error {
	// Open database
	connectionString := getConnectionString(v.Path)
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute insert statement
	_, err = db.Exec(sqlStatementUpdateDuration, videoID, duration)
	if err != nil {
		return err
	}
	return nil
}
