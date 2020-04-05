package main

import (
	"path"

	_ "github.com/go-sql-driver/mysql"
	core "github.com/hultan/softtube/softtube.core"
)

var (
	logger   core.Logger
	config   *core.Config
	db       core.Database
	softtube *SoftTube
)

func main() {
	// Load the config file
	loadConfig()

	// Setup the client logging
	startLogging()
	defer stopLogging()

	// Open the SoftTube database
	openDatabase()
	defer db.CloseDatabase()

	startApplication(&db)
}

func loadConfig() {
	// Load config file
	config = new(core.Config)
	config.Load("main")
}

func startLogging() {
	// Start logging
	logger = core.NewLog(path.Join(config.ServerPaths.Log, config.Logs.SoftTube))
	logger.LogStart(config, "softtube client")
}

func stopLogging() {
	// Close log file
	logger.LogFinished("softtube client")
	logger.Close()
}

func openDatabase() core.Database {
	// Create the database object, and get all subscriptions
	conn := config.Connection
	crypt := core.Crypt{}
	password, err := crypt.Decrypt(conn.Password)
	if err != nil {
		logger.Log("Failed to decrypt MySQL password!")
		logger.LogError(err)
		panic(err)
	}

	db = core.New(conn.Server, conn.Port, conn.Database, conn.Username, password)
	db.OpenDatabase()
	return db
}

func startApplication(db *core.Database) {
	// Create a new application.
	softtube = new(SoftTube)
	err := softtube.StartApplication(db)
	if err != nil {
		logger.Log("Failed to start application!")
		logger.LogError(err)
		panic(err)
	}
}
