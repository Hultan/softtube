package main

import (
	"database/sql"
	"testing"
	"time"

	database "github.com/hultan/softtube/internal/softtube.database"
)

func Test_updateSubscription(t *testing.T) {
	subscription := database.Subscription{
		ID:          testKarenPuzzlesChannelId,
		Name:        "Karen puzzles",
		Frequency:   1,
		LastChecked: sql.NullTime{Time: time.Now(), Valid: true},
		NextUpdate:  sql.NullInt32{Int32: 0},
	}

	updateSubscription(&subscription)
}

func Test_updateSubscription2(t *testing.T) {
	subscription := database.Subscription{
		ID:          "UCsr4lPNPq2GZ76ocdNMAArA",
		Name:        "Mathew Santoro 2",
		Frequency:   1,
		LastChecked: sql.NullTime{Time: time.Now(), Valid: true},
		NextUpdate:  sql.NullInt32{Int32: 0},
	}

	updateSubscription(&subscription)
}
