package main

import (
	"github.com/hultan/softtube/internal/softtube.database"
	"path"

	_ "github.com/go-sql-driver/mysql"
	core "github.com/hultan/softtube/internal/softtube.core"
)

var (
	logger   *core.Logger
	config   *core.Config
	db       *database.Database
	softTube *SoftTube
)

func main() {
	// Load the config file
	loadConfig()

	// Setup the client logging
	startLogging()
	defer stopLogging()

	// Open the SoftTube database
	_ = openDatabase()
	defer closeDatabase()

	startApplication(db)
}

func loadConfig() {
	// Load config file
	config = new(core.Config)
	err := config.Load("main")
	if err!=nil {
		panic("failed to load config")
	}
}

func startLogging() {
	// Start logging
	logger = core.NewLog(path.Join(config.ServerPaths.Log, config.Logs.SoftTube))
	logger.LogStart("softtube client")
}

func stopLogging() {
	// Close log file
	logger.LogFinished("softtube client")
	logger.Close()
}

func openDatabase() *database.Database {
	// Create the database object, and get all subscriptions
	conn := config.Connection
	crypto := core.Crypt{}
	password, err := crypto.Decrypt(conn.Password)
	if err != nil {
		logger.Log("Failed to decrypt MySQL password!")
		logger.LogError(err)
		panic(err)
	}

	db = database.New(conn.Server, conn.Port, conn.Database, conn.Username, password)
	err = db.OpenDatabase()
	if err!=nil {
		return nil
	}
	return db
}

func closeDatabase() {
	db.CloseDatabase()
}

func startApplication(db *database.Database) {
	// Create a new application.
	softTube = new(SoftTube)
	err := softTube.StartApplication(db)
	if err != nil {
		logger.Log("Failed to start application!")
		logger.LogError(err)
		panic(err)
	}
}
