package main

import (
	"path"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hultan/softteam-tools/pkg/crypt"
	"github.com/hultan/softteam-tools/pkg/log"
	core "github.com/hultan/softtube/internal/softtube.core"
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
	crypto := crypt.Crypt{}
	password, err := crypto.Decrypt(conn.Password)
	if err != nil {
		logger.Log("Failed to decrypt MySQL password!")
		logger.LogError(err)
		panic(err)
	}

	db = core.New(conn.Server, conn.Port, conn.Database, conn.Username, password)
	err = db.OpenDatabase()
	if err!=nil {
		return nil
	}
	return db
}

func closeDatabase() {
	db.CloseDatabase()
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
