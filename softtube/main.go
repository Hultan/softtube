package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	core "github.com/hultan/softtube/softtube.core"
)

func main() {
	testCreateDatabase()
}

func testCreateDatabase() {
	db := core.New("192.168.1.3", 3306, "softtube", "per", "KnaskimGjwQ6M!")
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

	versionNumber := fmt.Sprintf("Version : %v.%v.%v", version.Major, version.Minor, version.Revision)
	fmt.Println(versionNumber)
}
