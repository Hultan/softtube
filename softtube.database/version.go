package database

import (
	"database/sql"

	entities "github.com/hultan/softtube/softtube.entities"
)

// VersionTable : The version table in the database
type VersionTable struct {
	Path string
}

// sql : Get version
const sqlStatementGetVersion = "select version from Version limit 1"

// GetVersion : Get the version number of a SoftTube database
func (v VersionTable) GetVersion() (entities.Version, error) {
	// Open database
	connectionString := getConnectionString(v.Path)
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return entities.Version{Major: 0}, err
	}

	rows, err := db.Query(sqlStatementGetVersion)
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
