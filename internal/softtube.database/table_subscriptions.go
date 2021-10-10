package database

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

// SubscriptionTable : SubscriptionTable in the SoftTube database
type SubscriptionTable struct {
	Connection *sql.DB
}

// sql : Get all subscriptions
const sqlStatementGetAllSubscriptions = "SELECT id, name, frequency, last_checked, next_update FROM Subscriptions"
const sqlStatementGetSubscription = "SELECT id, name, frequency, last_checked, next_update FROM Subscriptions WHERE id=?"
const sqlStatementUpdateLastChecked = "UPDATE Subscriptions SET last_checked=?, next_update=? WHERE id=?"

// GetAll : Returns all subscriptions
func (s SubscriptionTable) GetAll() ([]Subscription, error) {
	// Check that database is opened
	if s.Connection == nil {
		return nil, errors.New("database not opened")
	}

	rows, err := s.Connection.Query(sqlStatementGetAllSubscriptions)
	if err != nil {
		return []Subscription{}, err
	}

	var subs []Subscription

	for rows.Next() {
		sub := new(Subscription)
		err = rows.Scan(&sub.ID, &sub.Name, &sub.Frequency, &sub.LastChecked, &sub.NextUpdate)
		if err != nil {
			return []Subscription{}, err
		}
		subs = append(subs, *sub)
	}

	_ = rows.Close()

	return subs, nil
}

// Get : Returns a subscription
func (s SubscriptionTable) Get(id string) (Subscription, error) {
	// Check that database is opened
	if s.Connection == nil {
		return Subscription{}, errors.New("database not opened")
	}

	row := s.Connection.QueryRow(sqlStatementGetSubscription, id)

	sub := Subscription{}
	err := row.Scan(&sub.ID, &sub.Name, &sub.Frequency, &sub.LastChecked, &sub.NextUpdate)

	// Return subscription
	return sub, err
}

// UpdateLastChecked : Update last_checked and next_update for a subscription
func (s SubscriptionTable) UpdateLastChecked(subscription *Subscription, interval int) error {
	// Check that database is opened
	if s.Connection == nil {
		return errors.New("database not opened")
	}

	rand.Seed(time.Now().UnixNano())
	now := time.Now().UTC().Format(constDateLayout) // last_checked
	next := int(float32(interval)*0.5 + (float32(interval) * rand.Float32()))

	// Execute insert statement
	_, err := s.Connection.Exec(sqlStatementUpdateLastChecked, now, next, subscription.ID)
	if err != nil {
		return err
	}

	return nil
}

// NeedsUpdate : Does the subscription need an update?
func (s Subscription) NeedsUpdate() bool {
	// If the fields is null, we must update
	if !s.LastChecked.Valid {
		return true
	}
	if !s.NextUpdate.Valid {
		return true
	}
	next := s.LastChecked.Time.Add(time.Duration(s.NextUpdate.Int32) * time.Second)
	now := time.Now().Local()
	result := next.Before(now)
	return result
}
