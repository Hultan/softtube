package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// TODO : Testa att byta ut databas drivern mot
// gorm.io/driver/mysql
// samt
// gorm.io/gorm
// Se video : https://www.youtube.com/watch?v=zTnkskp-xWs

type ConnectionInfo struct {
	Server   string
	Port     int
	Database string
	Username string
	Password string
}

// Database : A connection to the SoftTube database
type Database struct {
	ConnectionInfo

	Connection    *sql.DB
	Subscriptions *SubscriptionTable
	Videos        *VideosTable
	Version       *VersionTable
	Download      *DownloadTable
	Log           *LogTable
}

// NewDatabase : Creates a new database object
func NewDatabase(server string, port int, database, username, password string) *Database {
	c := ConnectionInfo{Server: server, Port: port, Database: database, Username: username, Password: password}
	return &Database{ConnectionInfo: c}
}

// Open : Open the database
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
	d.Version = &VersionTable{t}
	d.Download = &DownloadTable{t}
	d.Log = &LogTable{t}

	return nil
}

// Close : Close the database
func (d *Database) Close() {
	if d.Connection != nil {
		_ = d.Connection.Close()
	}
}

// getConnectionString : Returns the connections string
func (c *ConnectionInfo) String() string {
	// "user:pwd@tcp(server:port)/database"
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", c.Username, c.Password, c.Server, c.Port, c.Database)
}
