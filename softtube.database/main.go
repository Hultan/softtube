package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Connect : Connect to a sqlite3 database
func Connect() {
	db, err := sql.Open("sqlite3", "/home/per/temp/test.db")
	checkErrors(err)

	statement := "select * from Videos"
	rows, err := db.Query(statement)
	checkErrors(err)
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name)
	}

	fmt.Println("Connected!")
}
