package database

// DownloadTable in the SoftTube database
type DownloadTable struct {
	*Table
}

// Insert a new download request into the database
func (d DownloadTable) Insert(id string) error {
	// Check that the database is opened
	if d.Connection == nil {
		return ErrDatabaseNotOpened
	}

	// Execute insert statement
	_, err := d.Connection.Exec(sqlDownloadsInsert, id)
	if err != nil {
		return err
	}

	return nil
}

// Delete the row from the downloaded list
func (d DownloadTable) Delete(id string) error {
	// Check that the database is opened
	if d.Connection == nil {
		return ErrDatabaseNotOpened
	}

	// Execute statement
	_, err := d.Connection.Exec(sqlDownloadsDelete, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAll returns all download requests
func (d DownloadTable) GetAll() ([]Download, error) {
	// Check that the database is opened
	if d.Connection == nil {
		return nil, ErrDatabaseNotOpened
	}

	rows, err := d.Connection.Query(sqlDownloadsGetAll)
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
