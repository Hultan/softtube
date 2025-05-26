package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// TODO : Try to exchange the MySQL driver with
// 		gorm.io/driver/mysql
// and
// 		gorm.io/gorm
// Check video : https://www.youtube.com/watch?v=zTnkskp-xWs

type ConnectionInfo struct {
	Server   string
	Port     int
	Database string
	Username string
	Password string
}

// Database represents connection to the SoftTube database
type Database struct {
	ConnectionInfo

	Connection    *sql.DB
	Subscriptions *SubscriptionTable
	Videos        *VideosTable
	Download      *DownloadTable
	Log           *LogTable
}

type Table struct {
	*Database
}

// NewDatabase creates a new database object
func NewDatabase(server string, port int, database, username, password string) *Database {
	c := ConnectionInfo{Server: server, Port: port, Database: database, Username: username, Password: password}
	return &Database{ConnectionInfo: c}
}

// Open the database
func (d *Database) Open() error {
	// Open database
	conn, err := sql.Open(constDriverName, d.ConnectionInfo.String())
	if err != nil {
		return err
	}
	d.Connection = conn
	t := &Table{d}
	d.Subscriptions = &SubscriptionTable{t}
	d.Videos = &VideosTable{t}
	d.Download = &DownloadTable{t}
	d.Log = &LogTable{t}

	return nil
}

// Close the database
func (d *Database) Close() {
	if d.Connection != nil {
		_ = d.Connection.Close()
	}
}

// getConnectionString returns the connection string
func (c *ConnectionInfo) String() string {
	// "user:pwd@tcp(server:port)/database"
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", c.Username, c.Password, c.Server, c.Port, c.Database)
}
