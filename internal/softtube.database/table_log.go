package database

import (
	"database/sql"
	"errors"
)

// LogTable : The version table in the database
type LogTable struct {
	Connection *sql.DB
}

// sql : Get version
const sqlStatementInsertLog = `INSERT INTO Log (type, message) VALUES (?, ?);`
const sqlStatementGetLogs = `SELECT id, type, message FROM Log                 
ORDER BY id desc
LIMIT 50`

// Insert : Insert a new video into the database
func (l *LogTable) Insert(logType int, logMessage string) error {
	// Check that database is opened
	if l.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute insert statement
	_, err := l.Connection.Exec(sqlStatementInsertLog, logType, logMessage)
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
	rows, err := l.Connection.Query(sqlStatementGetLogs)
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
