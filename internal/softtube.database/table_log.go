package database

// LogTable is the log table in the database
type LogTable struct {
	*Table
}

// Insert a new video into the database
func (l *LogTable) Insert(logType LogType, logMessage string) error {
	// Check that the database is opened
	if l.Connection == nil {
		return ErrDatabaseNotOpened
	}

	// Execute insert statement
	_, err := l.Connection.Exec(sqlLogInsert, logType, logMessage)
	if err != nil {
		return err
	}

	return nil
}

// GetLatest returns the version number of a SoftTube database
func (l *LogTable) GetLatest() ([]Log, error) {

	// Check that the database is opened
	if l.Connection == nil {
		return nil, ErrDatabaseNotOpened
	}

	// Get rows from the database
	rows, err := l.Connection.Query(sqlLogGetLatest)
	if err != nil {
		return nil, err
	}

	var logs []Log

	// Get logs from rows
	for rows.Next() {
		log := new(Log)
		err = rows.Scan(&log.ID, &log.Type, &log.Message)
		if err != nil {
			return nil, err
		}
		logs = append(logs, *log)
	}

	_ = rows.Close()

	// Return the logs
	return logs, nil
}
