package main

import (
	"fmt"
	"sync"

	database "github.com/hultan/softtube/softtube.database"
	entities "github.com/hultan/softtube/softtube.entities"
)

func main() {
	config := new(Config)
	config.Load()

	// Process command line args
	err := processCommandLineArgs(config)
	if err != nil {
		invalidCommandLineArg(err)
	}

	// Setup logging
	log := createAndOpenLog(config.Paths.Log)
	defer log.close()

	// Start updating the softtube database
	log.logStart(config)

	// Create the database object, and get all subscriptions
	db := database.New(config.Paths.Database)
	subs, err := db.Subscriptions.GetAll()

	// Handle errors
	if err != nil {
		log.log(err.Error())
	}

	// Log result
	log.logFormat("Loaded ", len(subs), " subscriptions.")

	// Check how many subscriptions that needs update
	subsThatNeedsUpdate := 0
	for i := 0; i < len(subs); i++ {
		sub := subs[i]

		if sub.NeedsUpdate() {
			subsThatNeedsUpdate++
		}
	}

	log.logFormat("softtube-update needs to update ", subsThatNeedsUpdate, " (of ", len(subs), " subscriptions).")

	// Create a waitgroup to sync the goroutines
	var waitGroup sync.WaitGroup
	waitGroup.Add(subsThatNeedsUpdate)

	for i := 0; i < len(subs); i++ {
		sub := subs[i]

		if sub.NeedsUpdate() {
			// Start goroutine to update subscription
			go func() {
				defer waitGroup.Done()
				update(&db, &log, config, &sub.Entity)
			}()
		}
	}

	waitGroup.Wait()

	// update(&log, config, &subs[21].Entity)

	log.logFinished()
}

func update(db *database.Database, log *Log, config *Config, subscription *entities.Subscription) {
	log.logFormat("Updating channel '", subscription.Name, "'.")

	youtube := new(youtube)
	rss, err := youtube.getSubscriptionRSS(subscription.ID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error updating %s : %s", subscription.ID, err.Error())
		log.log(errorMessage)
	}
	feed := new(Feed)
	feed.parse(rss)
	videos := feed.getVideos()
	var waitGroup sync.WaitGroup

	for i := 0; i < len(videos); i++ {
		video := videos[i]
		// Check if the video already exists in the database
		exists, err := db.Videos.Exists(video.ID)
		if err != nil {
			fmt.Println(err.Error())
		}

		if !exists {
			waitGroup.Add(3)

			go func() {
				defer waitGroup.Done()
				// Insert the video in the database
				db.Videos.Insert(video.ID, video.ChannelID, video.Title, "", video.Published)
			}()
			go func() {
				// Get duration
				defer waitGroup.Done()
				youtube.getDuration(db, config, video.ID)
			}()
			go func() {
				// Get thumbnail
				defer waitGroup.Done()
				youtube.getThumbnail(config, video.ID)
			}()
		}
	}

	waitGroup.Wait()
}
