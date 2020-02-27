package database

import "log"

func checkErrors(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
