package core

import (
	"database/sql"
	"errors"
)

// DownloadTable : VideosTable in the SoftTube database
type DownloadTable struct {
	Connection *sql.DB
}

// TODO : Make max downloads a setting
const sqlStatementGetDownloads = "SELECT video_id FROM Download LIMIT 5"
const sqlStatementInsertDownload = "INSERT INTO Download (video_id) VALUES (?)"
const sqlStatementSetAsDownloaded = "DELETE FROM Download WHERE video_id=?"

// Insert : Insert a new download request into the database
func (d DownloadTable) Insert(id string) error {
	// Check that database is opened
	if d.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute insert statement
	_, err := d.Connection.Exec(sqlStatementInsertDownload, id)
	if err != nil {
		return err
	}

	return nil
}

// SetAsDownloaded : Deletes the row from the downloaded list
func (d DownloadTable) SetAsDownloaded(id string) error {
	// Check that database is opened
	if d.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute statement
	_, err := d.Connection.Exec(sqlStatementSetAsDownloaded, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAll : Returns all download requests
func (d DownloadTable) GetAll() ([]Download, error) {
	// Check that database is opened
	if d.Connection == nil {
		return nil, errors.New("database not opened")
	}

	rows, err := d.Connection.Query(sqlStatementGetDownloads)
	if err != nil {
		return []Download{}, err
	}

	var downloads []Download

	for rows.Next() {
		download := new(Download)
		err = rows.Scan(&download.ID)
		if err != nil {
			return []Download{}, err
		}
		downloads = append(downloads, *download)
	}

	_ = rows.Close()

	return downloads, nil
}
