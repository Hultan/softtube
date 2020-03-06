package main

import (
	"errors"
	"fmt"
	"sync"

	core "github.com/hultan/softtube/softtube.core"
)

const applicationVersion string = "1.00"
const maxUpdates = 50

var (
	logger Log
	config *core.Config
	db     core.Database
)

func main() {
	// Process command line args
	err := processCommandLineArgs()
	if err != nil {
		invalidCommandLineArg(err)
	}

	// Load config file
	config = new(core.Config)
	config.Load("main")

	// Setup logging
	logger = createAndOpenLog(config.Update.Log)
	defer logger.close()

	// Start updating the softtube database
	logger.logStart(config)
	defer logger.logFinished()

	conn := config.Connection

	// Create the database object, and get all subscriptions
	db = core.New(conn.Server, conn.Port, conn.Database, conn.Username, conn.Password)
	db.OpenDatabase()
	defer db.CloseDatabase()
	subs, err := db.Subscriptions.GetAll()

	// Handle errors
	if err != nil {
		logger.log(err.Error())
	}

	// Log result
	logger.logFormat("Loaded ", len(subs), " subscriptions.")

	// Check how many subscriptions that needs update
	subsThatNeedsUpdate := 0
	for i := 0; i < len(subs); i++ {
		sub := subs[i]

		if sub.NeedsUpdate() {
			subsThatNeedsUpdate++
		}
	}

	if subsThatNeedsUpdate > maxUpdates {
		subsThatNeedsUpdate = maxUpdates
	}

	logger.logFormat("softtube-update needs to update ", subsThatNeedsUpdate, " (of ", len(subs), " subscriptions).")

	// Create a waitgroup to sync the goroutines
	var waitGroup sync.WaitGroup
	waitGroup.Add(subsThatNeedsUpdate)

	count := 0

	for i := 0; i < len(subs); i++ {
		sub := subs[i]

		if sub.NeedsUpdate() {
			if count > maxUpdates {
				logger.log("Max number of updates reached!")
				break
			}
			count++
			// Start goroutine to update subscription
			go func() {
				defer waitGroup.Done()
				updateSubscription(&sub)
			}()
		}
	}

	waitGroup.Wait()

	//updateSubscription(&subs[21])
}

func updateSubscription(subscription *core.Subscription) {
	logger.logFormat("Updating channel '", subscription.Name, "'.")

	// Download the subscription RSS
	youtube := new(youtube)
	rss, err := youtube.getSubscriptionRSS(subscription.ID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error updating %s : %s", subscription.ID, err.Error())
		logger.log(errorMessage)
	}

	// Parse the RSS
	feed := new(Feed)
	feed.parse(rss)

	// Get videos in the RSS
	videos := feed.getVideos()

	// Create the waitgroup to syncronize the goroutines below
	var waitGroup sync.WaitGroup

	for i := 0; i < len(videos); i++ {
		video := videos[i]

		// Check if the video already exists in the database
		exists, err := db.Videos.Exists(video.ID)
		if err != nil {
			fmt.Println(err.Error())
		}

		if !exists {
			message := fmt.Sprintf("New video for channel '%s' : '%s'", subscription.Name, video.Title)
			logger.log(message)

			// Insert the video in the database
			// This must be executed before getDuration()
			err := db.Videos.Insert(video.ID, video.SubscriptionID, video.Title, "", video.Published)
			if err != nil {
				message = fmt.Sprintf("Inserted video '%s' in database : Failed! (Reason : %s)", video.Title, err.Error())
				logger.log(message)
				return
			}
			message = fmt.Sprintf("Inserted video '%s' in database : Success!", video.Title)
			logger.log(message)

			waitGroup.Add(2)

			go func() {
				// Get duration
				defer waitGroup.Done()
				err := youtube.getDuration(video.ID, logger)
				if err != nil {
					message = fmt.Sprintf("Updated duration for video '%s' : Failed! (Reason : %s)", video.Title, err.Error())
					logger.log(message)
					return
				}
				message = fmt.Sprintf("Updated duration for video '%s' : Success!", video.Title)
				logger.log(message)
			}()
			go func() {
				// Get thumbnail
				defer waitGroup.Done()
				err := youtube.getThumbnail(video.ID, config.Update.Thumbnails, logger)
				if err != nil {
					message = fmt.Sprintf("Downloaded thumbnail for video '%s': Failed! (Reason : %s)", video.Title, err.Error())
					logger.log(message)
					return
				}
				message = fmt.Sprintf("Downloaded thumbnail for video '%s': Success!", video.Title)
				logger.log(message)
			}()
		}
	}

	// Mark subscription as updated
	interval, err := getInterval(subscription.Frequency)
	if err != nil {
		logger.logError(err)
	}
	db.Subscriptions.UpdateLastChecked(subscription, interval)

	waitGroup.Wait()
}

func getInterval(frequency int) (int, error) {
	switch frequency {
	case 1:
		return config.Intervals.High, nil
	case 2:
		return config.Intervals.Medium, nil
	case 3:
		return config.Intervals.Low, nil
	}
	return 0, errors.New("invalid frequency")
}
