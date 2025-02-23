package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// VideosTable : VideosTable in the SoftTube database
type VideosTable struct {
	*Table
}

// Get : Returns a subscription
func (v VideosTable) Get(id string) (Video, error) {
	// Check that database is opened
	if v.Connection == nil {
		return Video{}, errors.New("database not opened")
	}

	row := v.Connection.QueryRow(sqlVideosGet, id)

	video := Video{}
	var saved uint8
	err := row.Scan(
		&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added, &video.Status,
		&saved,
	)
	if saved == 1 {
		video.Saved = true
	}
	// Return video
	return video, err
}

// Exists : Does a video already exist in the database?
func (v VideosTable) Exists(videoID string) (bool, error) {
	// Check that database is opened
	if v.Connection == nil {
		return false, errors.New("database not opened")
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

	return false, fmt.Errorf("failed to check if video '%s' exists", videoID)
}

// GetStatus : Get the video status
func (v VideosTable) GetStatus(videoID string) (int, error) {
	// Check that database is opened
	if v.Connection == nil {
		return -1, errors.New("database not opened")
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

	return -1, fmt.Errorf("failed to check if video '%s' exists", videoID)
}

// Insert : Insert a new video into the database
func (v VideosTable) Insert(
	id string, subscriptionID string, title string, duration string, published time.Time,
) error {
	// Check that database is opened
	if v.Connection == nil {
		return errors.New("database not opened")
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

// UpdateStatus : Update the status for a video
func (v VideosTable) UpdateStatus(id string, status VideoStatusType) error {
	// Check that database is opened
	if v.Connection == nil {
		return errors.New("database not opened")
	}

	// Check if the connection is still valid
	if err := v.Connection.Ping(); err != nil {
		return fmt.Errorf("database connection lost: %v", err)
	}

	// Execute the update statement
	_, err := v.Connection.Exec(sqlVideosUpdateStatus, status, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSave : Update saved flag for a video
func (v VideosTable) UpdateSave(id string, saved bool) error {
	// Check that database is opened
	if v.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute the update statement
	_, err := v.Connection.Exec(sqlVideosUpdateSave, saved, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDuration : Update duration for a video
func (v VideosTable) UpdateDuration(videoID string, duration string) error {
	// Check that database is opened
	if v.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute insert statement
	_, err := v.Connection.Exec(sqlVideosUpdateDuration, duration, v.getSeconds(duration), videoID)
	if err != nil {
		return err
	}

	return nil
}

// Delete : Delete a video from the database
func (v VideosTable) Delete(id string) error {
	// Check that database is opened
	if v.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute delete query
	_, err := v.Connection.Exec(sqlVideosDelete, id)
	if err != nil {
		return err
	}

	return nil
}

// Search : Searches for videos
func (v VideosTable) Search(text string) ([]Video, error) {
	// Check that database is opened
	if v.Connection == nil {
		return nil, errors.New("database not opened")
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
			&video.Status, &video.SubscriptionName, &saved,
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

// GetVideos : Gets a list of the latest videos
func (v VideosTable) GetVideos(failed, savedView bool) ([]Video, error) {
	// Check that database is opened
	if v.Connection == nil {
		return nil, errors.New("database not opened")
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
			&video.Status, &video.SubscriptionName, &saved,
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
	fmt.Println("Duration: " + duration)
	if duration == "" {
		fmt.Println("No duration : cause empty")
		return 0
	}
	if duration == "LIVE" || duration == "MEMBER" || duration == "ERROR" {
		fmt.Println("No duration : cause '" + duration + "'")
		return 0
	}

	// Split the duration into parts
	parts := strings.Split(duration, ":")

	// Handle different formats based on the number of parts
	switch len(parts) {
	case 1:
		// Format: SS or S
		seconds, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("No duration : cause '" + err.Error() + "'")
			return 0
		}
		fmt.Printf("Duration : %d seconds\n", seconds)
		return seconds

	case 2:
		// Format: MM:SS or M:SS
		minutes, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("No duration : cause '" + err.Error() + "'")
			return 0
		}
		seconds, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("No duration : cause '" + err.Error() + "'")
			return 0
		}
		seconds = minutes*60 + seconds
		fmt.Printf("Duration : %d seconds\n", seconds)
		return seconds

	case 3:
		// Format: HH:MM:SS or H:MM:SS
		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("No duration : cause '" + err.Error() + "'")
			return 0
		}
		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("No duration : cause '" + err.Error() + "'")
			return 0
		}
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("No duration : cause '" + err.Error() + "'")
			return 0
		}
		seconds = hours*3600 + minutes*60 + seconds
		fmt.Printf("Duration : %d seconds\n", seconds)
		return seconds

	default:
		fmt.Println("No duration : cause invalid length")
		return 0
	}
}
