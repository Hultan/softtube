package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// VideosTable : VideosTable in the SoftTube database
type VideosTable struct {
	Connection *sql.DB
}

const sqlStatementGetFailedDownloads = `SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, 
										Videos.published, Videos.added, Videos.status, Subscriptions.name, Videos.save 
									FROM Videos 
									INNER JOIN Subscriptions ON Videos.subscription_id = Subscriptions.id
									WHERE Videos.status = 1
									ORDER BY added desc`
const sqlStatementVideoExists = "SELECT EXISTS(SELECT 1 FROM Videos WHERE id=?);"
const sqlStatementGetStatus = "SELECT status FROM Videos WHERE id=?"
const sqlStatementGetVideo = "SELECT id, subscription_id, title, duration, published, added, status, save FROM Videos WHERE id=?"
const sqlStatementInsertVideo = `INSERT IGNORE INTO Videos (id, subscription_id, title, duration, published, added, status, save) 
								VALUES (?, ?, ?, ?, ?, ?, 0, 0);`
const sqlStatementUpdateDuration = "UPDATE Videos SET duration=? WHERE id=?"
const sqlStatementDeleteVideo = "DELETE FROM Videos WHERE id=? AND save=0"
const sqlStatementUpdateStatus = "UPDATE Videos SET status=? WHERE id=?"
const sqlStatementUpdateSave = "UPDATE Videos SET save=? WHERE id=?"
const sqlStatementSearchVideos = `SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, 
									Videos.status, Subscriptions.name , Videos.save
									FROM Videos 
									INNER JOIN Subscriptions ON Subscriptions.id = Videos.subscription_id 
									WHERE Videos.title LIKE ? OR Subscriptions.name LIKE ? 
									ORDER BY Videos.Added DESC`

const sqlStatementGetLatest = `SELECT * FROM 
									(SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, Videos.status, Subscriptions.name, Videos.save 
									FROM Videos 
									INNER JOIN Subscriptions ON Videos.subscription_id = Subscriptions.id 
									ORDER BY added desc
									LIMIT 200) as Newest

									UNION

									SELECT * FROM
										(SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, Videos.status, Subscriptions.name, Videos.save 
										FROM Videos 
										INNER JOIN Subscriptions ON Videos.subscription_id = Subscriptions.id 
										WHERE Videos.status NOT IN (0,4) OR Videos.save=1) as Downloaded

									ORDER BY added desc`

// Get : Returns a subscription
func (v VideosTable) Get(id string) (Video, error) {
	// Check that database is opened
	if v.Connection == nil {
		return Video{}, errors.New("database not opened")
	}

	row := v.Connection.QueryRow(sqlStatementGetVideo, id)

	video := Video{}
	var saved uint8
	err := row.Scan(&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added, &video.Status, &saved)
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
	rows, err := v.Connection.Query(sqlStatementVideoExists, videoID)
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
	rows, err := v.Connection.Query(sqlStatementGetStatus, videoID)
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
func (v VideosTable) Insert(id string, subscriptionID string, title string, duration string, published time.Time) error {
	// Check that database is opened
	if v.Connection == nil {
		return errors.New("database not opened")
	}

	now := time.Now().UTC().Format(constDateLayout) // Added

	// Execute insert statement
	_, err := v.Connection.Exec(sqlStatementInsertVideo, id, subscriptionID, title, duration, published, now)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStatus : Update the status for a video
func (v VideosTable) UpdateStatus(id string, status int) error {
	// Check that database is opened
	if v.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute the update statement
	_, err := v.Connection.Exec(sqlStatementUpdateStatus, status, id)
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
	_, err := v.Connection.Exec(sqlStatementUpdateSave, saved, id)
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
	_, err := v.Connection.Exec(sqlStatementUpdateDuration, duration, videoID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteFromDatabase : Delete a video from the database
func (v VideosTable) DeleteFromDatabase(id string) error {
	// Check that database is opened
	if v.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute delete query
	_, err := v.Connection.Exec(sqlStatementDeleteVideo, id)
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
	rows, err := v.Connection.Query(sqlStatementSearchVideos, search, search)
	if err != nil {
		return nil, err
	}

	var videos []Video
	var saved uint8

	for rows.Next() {
		video := new(Video)
		err = rows.Scan(&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added, &video.Status, &video.SubscriptionName, &saved)
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
func (v VideosTable) GetVideos(failed bool) ([]Video, error) {
	// Check that database is opened
	if v.Connection == nil {
		return nil, errors.New("database not opened")
	}

	var sqlString string
	if failed {
		sqlString = sqlStatementGetFailedDownloads
	} else {
		sqlString = sqlStatementGetLatest
	}

	rows, err := v.Connection.Query(sqlString)
	if err != nil {
		return []Video{}, err
	}

	var videos []Video
	var saved uint8

	for rows.Next() {
		video := new(Video)
		err = rows.Scan(&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added, &video.Status, &video.SubscriptionName, &saved)
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
