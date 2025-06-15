package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// VideosTable in the SoftTube database
type VideosTable struct {
	*Table
}

// Get a subscription
func (v VideosTable) Get(id string) (Video, error) {
	// Check that the database is opened
	if v.Connection == nil {
		return Video{}, ErrDatabaseNotOpened
	}

	row := v.Connection.QueryRow(sqlVideosGet, id)

	video := Video{}
	var saved uint8
	err := row.Scan(
		&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published,
		&video.Added, &video.Status, &saved, &video.Seconds,
	)
	if saved == 1 {
		video.Saved = true
	}
	// Return video
	return video, err
}

// Exists returns true if a video already exists in the database
func (v VideosTable) Exists(videoID string) (bool, error) {
	// Check that the database is opened
	if v.Connection == nil {
		return false, ErrDatabaseNotOpened
	}

	// Execute select query
	rows, err := v.Connection.Query(sqlVideosExists, videoID)
	if err != nil {
		return false, err
	}

	if rows.Next() {
		var result bool
		err = rows.Scan(&result)
		if err != nil {
			return false, err
		}
		return result, nil
	}

	_ = rows.Close()

	return false, newErrVideo("failed to check if video exists", videoID)
}

// GetStatus gets the video status
func (v VideosTable) GetStatus(videoID string) (int, error) {
	// Check that the database is opened
	if v.Connection == nil {
		return -1, ErrDatabaseNotOpened
	}

	// Execute select query
	rows, err := v.Connection.Query(sqlVideosGetStatus, videoID)
	if err != nil {
		return -1, err
	}

	if rows.Next() {
		var result int
		err = rows.Scan(&result)
		if err != nil {
			return -1, err
		}
		return result, nil
	}

	_ = rows.Close()

	return -1, newErrVideo("failed to check if video exists", videoID)
}

// Insert a new video into the database
func (v VideosTable) Insert(
	id string, subscriptionID string, title string, duration string, published time.Time,
) error {
	// Check that the database is opened
	if v.Connection == nil {
		return ErrDatabaseNotOpened
	}

	if !strings.HasPrefix(subscriptionID, "UC") {
		subscriptionID = "UC" + subscriptionID
	}

	now := time.Now().UTC().Format(constDateLayout) // Added

	// Execute insert statement
	_, err := v.Connection.Exec(sqlVideosInsert, id, subscriptionID, title, duration, published, now,
		v.getSeconds(duration))
	if err != nil {
		return err
	}

	return nil
}

// UpdateStatus updates the status for a video
func (v VideosTable) UpdateStatus(id string, status VideoStatusType) error {
	// Check that the database is opened
	if v.Connection == nil {
		return ErrDatabaseNotOpened
	}

	// Check if the connection is still valid
	if err := v.Connection.Ping(); err != nil {
		return ErrDatabaseNotOpened
	}

	// Execute the update statement
	_, err := v.Connection.Exec(sqlVideosUpdateStatus, status, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSave updates the saved flag for a video
func (v VideosTable) UpdateSave(id string, saved bool) error {
	// Check that the database is opened
	if v.Connection == nil {
		return ErrDatabaseNotOpened
	}

	// Execute the update statement
	_, err := v.Connection.Exec(sqlVideosUpdateSave, saved, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDuration updates duration for a video
func (v VideosTable) UpdateDuration(videoID string, duration string) error {
	// Check that the database is opened
	if v.Connection == nil {
		return ErrDatabaseNotOpened
	}

	// Execute insert statement
	_, err := v.Connection.Exec(sqlVideosUpdateDuration, duration, v.getSeconds(duration), videoID)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a video from the database
func (v VideosTable) Delete(id string) error {
	// Check that the database is opened
	if v.Connection == nil {
		return ErrDatabaseNotOpened
	}

	// Execute delete query
	_, err := v.Connection.Exec(sqlVideosDelete, id)
	if err != nil {
		return err
	}

	return nil
}

// Search searches for videos
func (v VideosTable) Search(text string) ([]Video, error) {
	// Check that the database is opened
	if v.Connection == nil {
		return nil, ErrDatabaseNotOpened
	}

	search := fmt.Sprintf("%%%s%%", text)
	rows, err := v.Connection.Query(sqlVideosSearch, search, search)
	if err != nil {
		return nil, err
	}

	var videos []Video
	var saved uint8

	for rows.Next() {
		video := new(Video)
		err = rows.Scan(
			&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added,
			&video.Status, &video.SubscriptionName, &saved, &video.Seconds,
		)
		if err != nil {
			return nil, err
		}
		if saved == 1 {
			video.Saved = true
		}
		videos = append(videos, *video)
	}

	_ = rows.Close()

	return videos, nil
}

// GetVideos gets a list of the latest videos
func (v VideosTable) GetVideos(failed, savedView bool) ([]Video, error) {
	// Check that the database is opened
	if v.Connection == nil {
		return nil, ErrDatabaseNotOpened
	}

	var sqlString string
	if failed {
		sqlString = sqlVideosGetFailed
	} else {
		sqlString = sqlVideosGetLatest
		if savedView {
			sqlString = strings.Replace(sqlString, "$ORDER$", "subscription_id, added asc", -1)
		} else {
			sqlString = strings.Replace(sqlString, "$ORDER$", "added desc", -1)
		}
	}

	rows, err := v.Connection.Query(sqlString)
	if err != nil {
		return []Video{}, err
	}

	var videos []Video
	var saved uint8

	for rows.Next() {
		video := new(Video)
		err = rows.Scan(
			&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added,
			&video.Status, &video.SubscriptionName, &saved, &video.Seconds,
		)
		if err != nil {
			return []Video{}, err
		}
		if saved == 1 {
			video.Saved = true
		}
		videos = append(videos, *video)
	}

	_ = rows.Close()

	return videos, nil
}

func (v VideosTable) getSeconds(duration string) int {
	if duration == "" || duration == "LIVE" || duration == "MEMBER" || duration == "ERROR" {
		return 0
	}

	// Split the duration into parts
	parts := strings.Split(duration, ":")

	// Handle different formats based on the number of parts
	switch len(parts) {
	case 1:
		// Format: SS or S
		seconds, _ := strconv.Atoi(parts[0])
		return seconds

	case 2:
		// Format: MM:SS or M:SS
		minutes, _ := strconv.Atoi(parts[0])
		seconds, _ := strconv.Atoi(parts[1])
		seconds = minutes*60 + seconds
		return seconds

	case 3:
		// Format: HH:MM:SS or H:MM:SS
		hours, _ := strconv.Atoi(parts[0])
		minutes, _ := strconv.Atoi(parts[1])
		seconds, _ := strconv.Atoi(parts[2])
		seconds = hours*3600 + minutes*60 + seconds
		return seconds

	default:
		return 0
	}
}

// GetStats gets a list of the videos in DB
func (v VideosTable) GetStats() ([]string, error) {
	// Check that the database is opened
	if v.Connection == nil {
		return nil, ErrDatabaseNotOpened
	}

	rows, err := v.Connection.Query(sqlVideosGetStats)
	if err != nil {
		return []string{}, err
	}

	var videos []string

	for rows.Next() {
		video := ""
		err = rows.Scan(&video)
		if err != nil {
			return []string{}, err
		}
		videos = append(videos, video)
	}

	_ = rows.Close()

	return videos, nil
}

// HasVideosToDelete returns true if there are videos to delete
func (v VideosTable) HasVideosToDelete() bool {
	// Check that the database is opened
	if v.Connection == nil {
		return false
	}

	var sqlString = "SELECT 1 FROM softtube.Videos WHERE Videos.status IN (3) AND Videos.save=0 order by Videos.Added;"
	var exists int
	err := v.Connection.QueryRow(sqlString).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	} else if err != nil {
		return false
	}
	return true
}
