package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	core "github.com/hultan/softtube/softtube.core"
)

var (
	logger  core.Logger
	config  *core.Config
	db      core.Database
	session core.SessionIdentifier
)

func main() {

	// Load config file
	config = new(core.Config)
	config.Load("main")

	// Setup logging
	logger = core.NewLog(config.Client.Log)
	defer logger.Close()

	// Log start and finish
	logger.LogStart(config, "softtube client")
	defer logger.LogFinished("softtube client")

	// Create the database object, and get all subscriptions
	conn := config.Connection
	db = core.New(conn.Server, conn.Port, conn.Database, conn.Username, conn.Password)
	db.OpenDatabase()
	defer db.CloseDatabase()

	// Create session
	session, err := core.CreateSession(db)
	if err != nil {
		panic(err)
	}

	logger.Log(fmt.Sprintf("Session started : %s", session.Name))
}
