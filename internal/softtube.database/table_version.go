package database

import (
	"database/sql"
	"errors"
)

// VersionTable : The version table in the database
type VersionTable struct {
	Connection *sql.DB
}

// sql : Get version
const sqlStatementGetVersion = "select major,minor,revision from Version limit 1"

// GetVersion : Get the version number of a SoftTube database
func (v VersionTable) GetVersion() (Version, error) {
	var version Version

	// Check that database is opened
	if v.Connection == nil {
		return version, errors.New("database not opened")
	}

	// Get rows from database
	rows, err := v.Connection.Query(sqlStatementGetVersion)
	if err != nil {
		return version, err
	}

	// Get version from rows
	for rows.Next() {
		err = rows.Scan(&version.Major, &version.Minor, &version.Revision)
		if err != nil {
			return version, err
		}
	}

	_ = rows.Close()

	// Return current database version
	return version, nil
}
