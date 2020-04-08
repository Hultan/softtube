package core

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

const sqlStatementVideoExists = "SELECT EXISTS(SELECT 1 FROM Videos WHERE id=?);"
const sqlStatementGetVideo = "SELECT id, subscription_id, title, duration, published, added, status FROM Videos WHERE id=?"
const sqlStatementInsertVideo = `INSERT IGNORE INTO Videos (id, subscription_id, title, duration, published, added) 
								VALUES (?, ?, ?, ?, ?, ?);`
const sqlStatementUpdateDuration = "UPDATE Videos SET duration=? WHERE id=?"
const sqlStatementDeleteVideo = "DELETE FROM Videos WHERE id=?"
const sqlStatementUpdateStatus = "UPDATE Videos SET status=? WHERE id=?"
const sqlStatementSearchVideos = `SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, 
									Videos.status, Subscriptions.name 
									FROM Videos 
									INNER JOIN Subscriptions ON Subscriptions.id = Videos.subscription_id 
									WHERE Videos.title LIKE ? OR Subscriptions.name LIKE ? 
									ORDER BY Videos.Added DESC`

// TODO : Make a setting of max number of videos
//const sqlStatementGetLatestVideos = "SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, Videos.status, Subscriptions.name FROM Videos INNER JOIN Subscriptions ON Subscriptions.id = Videos.subscription_id ORDER BY Videos.Added DESC LIMIT 200;"
const sqlStatementGetLatest = `SELECT * FROM 
(SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, Videos.status, Subscriptions.name 
FROM Videos 
INNER JOIN Subscriptions ON Videos.subscription_id = Subscriptions.id 
ORDER BY added desc
LIMIT 200) as Newest

UNION

SELECT * FROM
	(SELECT Videos.id, Videos.subscription_id, Videos.title, Videos.duration, Videos.published, Videos.added, Videos.status, Subscriptions.name 
	FROM Videos 
	INNER JOIN Subscriptions ON Videos.subscription_id = Subscriptions.id 
	WHERE Videos.status NOT IN (0,4)) as Downloaded

ORDER BY added desc`

// Get : Returns a subscription
func (v VideosTable) Get(id string) (Video, error) {
	// Check that database is opened
	if v.Connection == nil {
		return Video{}, errors.New("database not opened")
	}

	row := v.Connection.QueryRow(sqlStatementGetVideo, id)

	video := Video{}
	err := row.Scan(&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added, &video.Status)
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
	defer rows.Close()

	if rows.Next() {
		var result bool
		err = rows.Scan(&result)
		if err != nil {
			return false, err
		}
		return result, nil
	}

	return false, fmt.Errorf("failed to check if video '%s' exists", videoID)
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
	defer rows.Close()

	var videos []Video

	for rows.Next() {
		video := new(Video)
		err = rows.Scan(&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added, &video.Status, &video.SubscriptionName)
		if err != nil {
			return nil, err
		}
		videos = append(videos, *video)
	}

	return videos, nil
}

// GetVideos : Gets a list of the latest videos
func (v VideosTable) GetVideos() ([]Video, error) {
	// Check that database is opened
	if v.Connection == nil {
		return nil, errors.New("database not opened")
	}

	rows, err := v.Connection.Query(sqlStatementGetLatest)
	if err != nil {
		return []Video{}, err
	}
	defer rows.Close()

	var videos []Video
	//const sqlStatementGetVideo = "SELECT id, subscription_id, title, duration, published, added, status FROM Videos WHERE id=?"

	for rows.Next() {
		video := new(Video)
		err = rows.Scan(&video.ID, &video.SubscriptionID, &video.Title, &video.Duration, &video.Published, &video.Added, &video.Status, &video.SubscriptionName)
		if err != nil {
			return []Video{}, err
		}
		videos = append(videos, *video)
	}

	return videos, nil
}
