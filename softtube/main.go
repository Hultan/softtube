package main

import (
	"fmt"

	database "github.com/hultan/softtube/softtube.database"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	testCreateDatabase()
}

func testCreateDatabase() {
	db := database.New("/home/per/temp/test.db")
	//database := database.Database{Path: "/home/per/temp/test.db"}
	err := db.OpenDatabase()
	if err != nil {
		panic(err)
	}
	defer db.CloseDatabase()

	version, err := db.Version.GetVersion()
	if err != nil {
		panic(err)
	}

	fmt.Println("Major : ", version.Major)
}
