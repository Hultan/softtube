package database

import (
	"database/sql"
	"errors"

	entities "github.com/hultan/softtube/softtube.entities"
)

// VersionTable : The version table in the database
type VersionTable struct {
	Database *sql.DB
}

// sql : Get version
const sqlStatementGetVersion = "select version from Version limit 1"

// GetVersion : Get the version number of a SoftTube database
func (v VersionTable) GetVersion() (entities.Version, error) {
	// Check that database is opened
	if v.Database == nil {
		return entities.Version{Major: 0}, errors.New("database not opened")
	}

	rows, err := v.Database.Query(sqlStatementGetVersion)
	if err != nil {
		return entities.Version{Major: 0}, err
	}
	defer rows.Close()

	var version int
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			return entities.Version{Major: 0}, err
		}
	}

	// Return current database version
	return entities.Version{Major: version}, nil
}
