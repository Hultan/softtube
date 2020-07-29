package main

import (
	"path"

	_ "github.com/go-sql-driver/mysql"
	crypt "github.com/hultan/softteam/crypt"
	log "github.com/hultan/softteam/log"
	core "github.com/hultan/softtube/softtube.core"
)

var (
	logger   *log.Logger
	config   *core.Config
	db       *core.Database
	softTube *SoftTube
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

	startApplication(db)
}

func loadConfig() {
	// Load config file
	config = new(core.Config)
	config.Load("main")
}

func startLogging() {
	// Start logging
	logger = log.NewLog(path.Join(config.ServerPaths.Log, config.Logs.SoftTube))
	logger.LogStart("softtube client")
}

func stopLogging() {
	// Close log file
	logger.LogFinished("softtube client")
	logger.Close()
}

func openDatabase() *core.Database {
	// Create the database object, and get all subscriptions
	conn := config.Connection
	crypt := crypt.Crypt{}
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
	softTube = new(SoftTube)
	err := softTube.StartApplication(db)
	if err != nil {
		logger.Log("Failed to start application!")
		logger.LogError(err)
		panic(err)
	}
}

// func getExecutablePath() string {
// 	ex, err := os.Executable()
// 	if err != nil {
// 		return ""
// 	}
// 	return filepath.Dir(ex)
// }

// func getResourcePath() string {
// 	return path.Join(getExecutablePath(), "resources")
// }
