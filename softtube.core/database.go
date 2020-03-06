package core

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Database : A connection to the SoftTube database
type Database struct {
	Connection       *sql.DB
	ConnectionString string
	Server           string
	Port             int
	Database         string
	Username         string
	Password         string
	Subscriptions    SubscriptionTable
	Videos           VideosTable
	Version          VersionTable
}

// New : Creates a new database object
func New(server string, port int, database, username, password string) Database {
	return Database{Server: server, Port: port, Database: database, Username: username, Password: password}
}

// OpenDatabase : Open the database
func (d *Database) OpenDatabase() error {
	// Open database
	d.ConnectionString = d.getConnectionString()
	conn, err := sql.Open(constDriverName, d.ConnectionString)
	if err != nil {
		d.ConnectionString = ""
		return err
	}
	d.Connection = conn
	d.Subscriptions = SubscriptionTable{Connection: conn}
	d.Videos = VideosTable{Connection: conn}
	d.Version = VersionTable{Connection: conn}
	return nil
}

// CloseDatabase : Close the database
func (d *Database) CloseDatabase() {
	if d.Connection != nil {
		d.Connection.Close()
	}
}

// getConnectionString : Returns the connections string
func (d *Database) getConnectionString() string {
	//"user:pwd@tcp(server:port)/database"
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", d.Username, d.Password, d.Server, d.Port, d.Database)
}
