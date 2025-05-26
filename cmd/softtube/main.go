package main

import (
	"path"

	"github.com/hultan/crypto"
	"github.com/hultan/softtube/internal/logger"
	"github.com/hultan/softtube/internal/softtube.database"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hultan/softtube/internal/softtube"
	core "github.com/hultan/softtube/internal/softtube.core"
)

var (
	log      *logger.Logger
	config   *core.Config
	db       *database.Database
	softTube *softtube.SoftTube
)

func main() {
	// Init the config file
	loadConfig()

	// Set up the client logging
	startLogging()
	defer stopLogging()

	// Open the SoftTube database
	_ = openDatabase()
	defer closeDatabase()

	startApplication()
}

func loadConfig() {
	// Init config file
	config = new(core.Config)
	err := config.Load("main")
	if err != nil {
		panic("failed to load config")
	}
}

func startLogging() {
	var err error

	// Start logging
	log, err = logger.NewStandardLogger(path.Join(config.ServerPaths.Log, config.Logs.SoftTube))
	if err != nil {
		panic(err)
	}
}

func stopLogging() {
	// Close the log file
	log.Close()
}

func openDatabase() *database.Database {
	// Create the database object and get all subscriptions
	conn := config.Connection
	password, err := crypto.Decrypt(conn.Password)
	if err != nil {
		log.Info.Println("Failed to decrypt MySQL password!")
		log.Info.Println(err)
		panic(err)
	}

	db = database.NewDatabase(conn.Server, conn.Port, conn.Database, conn.Username, password)
	err = db.Open()
	if err != nil {
		return nil
	}
	return db
}

func closeDatabase() {
	db.Close()
}

func startApplication() {
	// Create a new application.
	softTube = &softtube.SoftTube{
		Config: config,
		Logger: log,
		DB:     db,
	}
	err := softTube.StartApplication()
	if err != nil {
		log.Info.Println("Failed to start application!")
		log.Info.Println(err)
		panic(err)
	}
}
