package test

import (
	"testing"
	"time"

	"github.com/hultan/softtube/internal/softtube.database"

	core "github.com/hultan/softtube/internal/softtube.core"
)

// TestSubscriptionGetAll : test Subscription.GetAll()
func TestSubscriptionGetAll(t *testing.T) {
	db, err := openDatabase(t)
	if err != nil {
		t.Errorf("TestSubscriptionGetAll: Failed to open database : %s", err.Error())
		return
	}
	subs, err := db.Subscriptions.GetAll()
	if err != nil {
		t.Errorf("TestSubscriptionGetAll: GetAll() returned an error : %s", err.Error())
		return
	}
	if len(subs) == 0 {
		t.Errorf("TestSubscriptionGetAll: GetAll() returned o rows")
		return
	}

	if subs[0].Name == "" {
		t.Errorf("TestSubscriptionGetAll: First subscription does not have a name")
		return
	}
}

// TestSubscriptionGet : test Subscription.Get()
func TestSubscriptionGet(t *testing.T) {
	db, err := openDatabase(t)
	if err != nil {
		t.Errorf("TestSubscriptionGet: Failed to open database : %s", err.Error())
		return
	}

	// Get a subscription
	sub, err := db.Subscriptions.Get("UCrjliuC6PeXJ4P8XYl_iETQ")
	if err != nil {
		t.Errorf("TestSubscriptionGet: Get() returned an error : %s", err.Error())
		return
	}

	// Check if the name is correct
	if sub.Name != "ThioJoeTech" {
		t.Errorf("TestSubscriptionGet: Get() returned an invalid subscription!")
		return
	}
}

// TestSubscriptionUpdateLastChecked : test Subscription.UpdateLastChecked()
func TestSubscriptionUpdateLastChecked(t *testing.T) {
	db, err := openDatabase(t)
	if err != nil {
		t.Errorf("TestSubscriptionUpdateLastChecked: Failed to open database : %s", err.Error())
		return
	}

	// Get a subscription
	sub, err := db.Subscriptions.Get("UCrjliuC6PeXJ4P8XYl_iETQ")
	if err != nil {
		t.Errorf("TestSubscriptionUpdateLastChecked: Get() returned an error : %s", err.Error())
		return
	}

	next := sub.NextUpdate

	// Update last checked
	err = db.Subscriptions.UpdateLastChecked(&sub, 3600)
	if err != nil {
		t.Errorf("TestSubscriptionUpdateLastChecked: UpdateLastChecked() returned an error : %s", err.Error())
		return
	}

	// Get the subscription again
	sub, err = db.Subscriptions.Get("UCrjliuC6PeXJ4P8XYl_iETQ")
	if err != nil {
		t.Errorf("TestSubscriptionUpdateLastChecked: Second Get() returned an error : %s", err.Error())
		return
	}

	// Make sure it is updated
	if next == sub.NextUpdate {
		t.Errorf("TestSubscriptionUpdateLastChecked: Interval not updated!")
	}

}

// TestVideosExists : test if a video exists
func TestVideosExists(t *testing.T) {
	db, err := openDatabase(t)
	if err != nil {
		t.Errorf("TestVideosExists: Failed to open database : %s", err.Error())
		return
	}

	// Check if video exists
	exists, err := db.Videos.Exists("AYmTPeEj7_Q")
	if err != nil {
		t.Errorf("TestVideosExists: Get() returned an error : %s", err.Error())
		return
	}

	if !exists {
		t.Errorf("TestVideosExists: Second Get() returned an error : %s", err.Error())
		return
	}
}

// TestVideosInsert : test inserting a video
func TestVideosInsert(t *testing.T) {
	db, err := openDatabase(t)
	if err != nil {
		t.Errorf("TestVideosInsert: Failed to open database : %s", err.Error())
		return
	}

	now := time.Now().UTC()

	// Insert a new video
	err = db.Videos.Insert("TO_DELETE!", "UCrjliuC6PeXJ4P8XYl_iETQ", "TO DELETE!", "0:00", now)
	if err != nil {
		t.Errorf("TestVideosInsert: Insert() returned an error : %s", err.Error())
		return
	}

	// Get the new video
	_, err = db.Videos.Get("TO_DELETE!")
	if err != nil {
		t.Errorf("TestVideosInsert: Get() returned an error : %s", err.Error())
		return
	}

	// Delete the new video
	err = db.Videos.Delete("TO_DELETE!")
	if err != nil {
		t.Errorf("TestVideosInsert: Delete() returned an error : %s", err.Error())
		return
	}
}

func openDatabase(t *testing.T) (*database.Database, error) {
	config, err := getTestConfig()
	if err != nil {
		t.Errorf("openDatabase: Failed to get config : %s", err.Error())
		return nil, err
	}
	db := database.NewDatabase(config.Connection.Server, config.Connection.Port, config.Connection.Database,
		config.Connection.Username, config.Connection.Password)
	err = db.Open()
	if err != nil {
		t.Errorf("openDatabase: Failed to open database : %s", err.Error())
		return nil, err
	}

	return db, nil
}

func getTestConfig() (*core.Config, error) {
	// Init config file
	config := new(core.Config)
	err := config.Load("test")

	return config, err
}
