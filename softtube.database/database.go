package database

import (
	"database/sql"
)

// Database : A connection to the SoftTube database
type Database struct {
	Database      *sql.DB
	Path          string
	Subscriptions SubscriptionTable
	Videos        VideosTable
	Version       VersionTable
}

// New : Creates a new database object
func New(path string) Database {
	return Database{Path: path}
}

// OpenDatabase : Open the database
func (d *Database) OpenDatabase() error {
	// Open database
	connectionString := d.ConnectionString(d.Path)
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return err
	}
	d.Database = db
	d.Subscriptions = SubscriptionTable{Database: db}
	d.Videos = VideosTable{Database: db}
	d.Version = VersionTable{Database: db}
	return nil
}

// CloseDatabase : Close the database
func (d *Database) CloseDatabase() {
	if d.Database != nil {
		d.Database.Close()
	}
}

// ConnectionString : Returns the connections string
func (d *Database) ConnectionString(databasePath string) string {
	return "file:" + databasePath + "?parseTime=true&_timeout=5000"
}
