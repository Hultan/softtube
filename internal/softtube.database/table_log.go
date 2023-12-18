package database

import (
	"errors"
)

// LogTable : The version table in the database
type LogTable struct {
	*Table
}

// sql : Get version
const sqlLogInsert = `INSERT INTO Log (type, message, time) VALUES (?, ?, NOW());`
const sqlLogGetLatest = `SELECT id, type, message FROM Log                 
ORDER BY id desc
LIMIT 50`

// Insert : Insert a new video into the database
func (l *LogTable) Insert(logType LogType, logMessage string) error {
	// Check that database is opened
	if l.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute insert statement
	_, err := l.Connection.Exec(sqlLogInsert, logType, logMessage)
	if err != nil {
		return err
	}

	return nil
}

// GetLatest : Get the version number of a SoftTube database
func (l *LogTable) GetLatest() ([]Log, error) {

	// Check that database is opened
	if l.Connection == nil {
		return nil, errors.New("database not opened")
	}

	// Get rows from database
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
