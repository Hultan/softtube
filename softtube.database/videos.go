package database

import (
	"database/sql"
	"errors"
	"time"

	entities "github.com/hultan/softtube/softtube.entities"
)

// VideosTable : VideosTable in the SoftTube database
type VideosTable struct {
	Database *sql.DB
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
	// Check that database is opened
	if v.Database == nil {
		return false, errors.New("database not opened")
	}

	// Execute select query
	rows, err := v.Database.Query(sqlStatementVideoExists, videoID)
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
	// Check that database is opened
	if v.Database == nil {
		return errors.New("database not opened")
	}

	now := time.Now().UTC().Format(dateLayout) // Added

	// Execute insert statement
	_, err := v.Database.Exec(sqlStatementInsertVideo, videoID, channelID, title, duration, published, now)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDuration : Update duration for a video
func (v VideosTable) UpdateDuration(videoID string, duration string) error {
	// Check that database is opened
	if v.Database == nil {
		return errors.New("database not opened")
	}

	// Execute insert statement
	_, err := v.Database.Exec(sqlStatementUpdateDuration, videoID, duration)
	if err != nil {
		return err
	}

	return nil
}
