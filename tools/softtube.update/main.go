package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/hultan/softteam/framework"
	"github.com/hultan/softtube/internal/softtube.database"

	core "github.com/hultan/softtube/internal/softtube.core"
)

const applicationVersion string = "1.00"
const maxUpdates = 50

var (
	logger *core.Logger
	config *core.Config
	db     *database.Database
)

func main() {
	// Process command line args
	err := processCommandLineArgs()
	if err != nil {
		invalidCommandLineArg(err)
	}

	// Init config file
	config = new(core.Config)
	err = config.Load("main")
	if err != nil {
		fmt.Println("ERROR (Open config) : ", err.Error())
		os.Exit(1)
	}

	// Setup logging
	logger = core.NewLog(path.Join(config.ServerPaths.Log, config.Logs.Update))
	defer logger.Close()

	// Start updating the softtube database
	logger.LogStart("softtube update")

	conn := config.Connection
	fw := framework.NewFramework()
	password, err := fw.Crypto.Decrypt(conn.Password)
	if err != nil {
		logger.Log("Failed to decrypt MySQL password!")
		logger.LogError(err)
		panic(err)
	}

	// Create the database object, and get all subscriptions
	db = database.NewDatabase(conn.Server, conn.Port, conn.Database, conn.Username, password)
	err = db.Open()
	if err != nil {
		fmt.Println("ERROR (Open config) : ", err.Error())
		os.Exit(1)
	}

	defer db.Close()
	subs, err := db.Subscriptions.GetAll()

	// Handle errors
	if err != nil {
		logger.Log(err.Error())
	}

	// Log result
	logger.LogFormat("Loaded ", len(subs), " subscriptions.")

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

	logger.LogFormat("softtube-update needs to update ", subsThatNeedsUpdate, " (of ", len(subs), " subscriptions).")

	// Create a waitgroup to sync the goroutines
	var waitGroup sync.WaitGroup
	waitGroup.Add(subsThatNeedsUpdate)

	count := 0

	for i := 0; i < len(subs); i++ {
		sub := subs[i]

		if sub.NeedsUpdate() {
			if count > maxUpdates {
				logger.Log("Max number of updates reached!")
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

	// updateSubscription(&subs[21])
}

func getVideos(subscription *database.Subscription) []database.Video {
	// Download the subscription RSS
	yt := &youtube{}
	rss, err := yt.getSubscriptionRSS(subscription.ID)
	if err != nil {
		errorMessage := fmt.Sprintf(
			"Failed to get RSS feed for %s (%s) : %s",
			subscription.Name, subscription.ID, err.Error(),
		)
		logger.Log(errorMessage)
		return nil
	}

	if contains404(rss) {
		errorMessage := fmt.Sprintf(
			"Channel %s (%s) has been deleted. Please remove channel...",
			subscription.Name, subscription.ID,
		)
		logger.Log(errorMessage)
		return nil
	}

	// Parse the RSS
	feed := &Feed{}
	err = feed.parse(rss)
	if err != nil {
		errorMessage := fmt.Sprintf(
			"Failed to parse RSS feed for %s (%s) : %s",
			subscription.Name, subscription.ID, err.Error(),
		)
		logger.Log(errorMessage)
		return nil
	}

	// Get videos in the RSS
	videos, err := feed.getVideos()
	if err != nil {
		errorMessage := fmt.Sprintf(
			"Failed to get videos for %s (%s) : %s",
			subscription.Name, subscription.ID, err.Error(),
		)
		logger.Log(errorMessage)
		return nil
	}

	return videos
}

func updateSubscription(subscription *database.Subscription) {
	logger.LogFormat("Updating channel '", subscription.Name, "'.")

	// Get videos in the RSS
	videos := getVideos(subscription)
	if videos == nil {
		// We failed to get the videos, it is already logged
		// in getVideos(), so just return
		return
	}

	// Create the waitgroup to synchronize the goroutines below
	var waitGroup sync.WaitGroup

	for _, video := range videos {
		// Check if the video already exists in the database
		exists, err := db.Videos.Exists(video.ID)
		if err != nil {
			errorMessage := fmt.Sprintf(
				"Failed to check if video exists for %s (%s) : %s",
				subscription.Name, subscription.ID, err.Error(),
			)
			logger.Log(errorMessage)
			continue
		}

		if exists {
			continue
		}

		err = handleNewVideo(video, &waitGroup)
		if err != nil {
			errorMessage := fmt.Sprintf(
				"Failed to handle new video for %s (%s) : %s",
				subscription.Name, subscription.ID, err.Error(),
			)
			logger.Log(errorMessage)
			continue
		}
	}

	waitGroup.Wait()

	// Mark subscription as updated
	interval, err := getInterval(subscription.Frequency)
	if err != nil {
		logger.LogError(err)
	}
	_ = db.Subscriptions.UpdateLastChecked(subscription, interval)
}

func handleNewVideo(video database.Video, waitGroup *sync.WaitGroup) error {
	// Clean video title from invalid character
	// TODO : Remove this?
	video.Title = clean(video.Title)

	// Insert the video in the database
	// This must be executed before getDuration()
	err := db.Videos.Insert(video.ID, video.SubscriptionID, video.Title, "", video.Published)
	if err != nil {
		msg := fmt.Sprintf(
			"Inserted video '%s' in database : Failed! (Reason : %s)", video.Title, err.Error(),
		)
		logger.Log(msg)
		return err
	} else {
		msg := fmt.Sprintf("Inserted video '%s' in database : Success!", video.Title)
		logger.Log(msg)
	}

	waitGroup.Add(2)

	go func() {
		// Get duration
		defer waitGroup.Done()
		yt := &youtube{}
		err = yt.getDuration(video.ID)
		if err != nil {
			msg := fmt.Sprintf(
				"Updated duration for video '%s' : Failed! (Reason : %s)", video.Title, err.Error(),
			)
			logger.Log(msg)
			return
		} else {
			msg := fmt.Sprintf("Updated duration for video '%s' : Success!", video.Title)
			logger.Log(msg)
		}
	}()
	go func() {
		// Get thumbnail
		defer waitGroup.Done()
		yt := &youtube{}
		err = yt.getThumbnail(video.ID)
		if err != nil {
			msg := fmt.Sprintf(
				"Downloaded thumbnail for video '%s': Failed! (Reason : %s)", video.Title, err.Error(),
			)
			logger.Log(msg)
			return
		} else {
			msg := fmt.Sprintf("Downloaded thumbnail for video '%s': Success!", video.Title)
			logger.Log(msg)
		}
	}()

	return nil
}

func clean(title string) string {
	re := regexp.MustCompile("[^[:ascii:]åäö]")

	return re.ReplaceAllLiteralString(title, "")
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

func contains404(rss string) bool {
	return strings.Contains(rss, "Error 404 (Not Found)")
}
