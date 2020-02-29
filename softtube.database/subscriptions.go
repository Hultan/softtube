package database

import (
	"database/sql"
	"errors"
	"time"

	entities "github.com/hultan/softtube/softtube.entities"
)

// SubscriptionTable : SubscriptionTable in the SoftTube database
type SubscriptionTable struct {
	Database *sql.DB
}

// SubscriptionRecord : A single subscription in the SubscriptionTable
type SubscriptionRecord struct {
	Entity entities.Subscription
}

// sql : Get all subscriptions
const sqlStatementGetAllSubscriptions = "select channel_id, channel_name, frequency, cast(last_checked as varchar), next_update from Subscriptions"

// GetAll : Returns all subscriptions
func (s SubscriptionTable) GetAll() ([]SubscriptionRecord, error) {
	// Check that database is opened
	if s.Database == nil {
		return nil, errors.New("database not opened")
	}

	rows, err := s.Database.Query(sqlStatementGetAllSubscriptions)
	if err != nil {
		return []SubscriptionRecord{}, err
	}
	defer rows.Close()

	var subs []SubscriptionRecord

	for rows.Next() {
		sub := new(SubscriptionRecord)
		var dateString string = ""
		err = rows.Scan(&sub.Entity.ID, &sub.Entity.Name, &sub.Entity.Frequency, &dateString, &sub.Entity.NextUpdate)
		if err != nil {
			return []SubscriptionRecord{}, err
		}
		if dateString == "" {
			sub.Entity.LastChecked = time.Time{}
		} else {
			// Parse the date
			date, err := time.Parse(dateLayout, dateString)
			if err != nil {
				sub.Entity.LastChecked = time.Time{}
			} else {
				sub.Entity.LastChecked = localTime(date)
			}
		}
		subs = append(subs, *sub)
	}

	// Return subscriptions
	return subs, nil
}

// NeedsUpdate : Does the subscription need an update?
func (s SubscriptionRecord) NeedsUpdate() bool {

	next := s.Entity.LastChecked.Add(time.Duration(s.Entity.NextUpdate) * time.Second)
	now := time.Now().Local()
	result := next.Before(now)
	return result
}

func localTime(datetime time.Time) time.Time {
	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		panic(err)
	}
	return datetime.In(loc)
}
