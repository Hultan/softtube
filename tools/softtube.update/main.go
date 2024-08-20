package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/hultan/crypto"
	log "github.com/hultan/softtube/internal/logger"
	"github.com/hultan/softtube/internal/softtube.database"

	core "github.com/hultan/softtube/internal/softtube.core"
)

const applicationVersion string = "1.11"
const maxUpdates = 50

var (
	logger *log.Logger
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
	config = &core.Config{}
	err = config.Load("main")
	if err != nil {
		fmt.Println("ERROR (Open config) : ", err.Error())
		os.Exit(1)
	}

	// Setup logging
	logger, err = log.NewStandardLogger(path.Join(config.ServerPaths.Log, config.Logs.Update))
	defer logger.Close()

	logger.Info.Println("")
	logger.Info.Println("---------------")
	logger.Info.Println("softtube.update")
	logger.Info.Println("---------------")
	logger.Info.Println("")

	// Start updating the softtube database
	conn := config.Connection
	c := &crypto.Crypto{}
	password, err := c.Decrypt(conn.Password)
	if err != nil {
		logger.Error.Println("Failed to decrypt MySQL password!")
		logger.Error.Println(err)
		os.Exit(1)
	}

	// Create the database object, and get all subscriptions
	db = database.NewDatabase(conn.Server, conn.Port, conn.Database, conn.Username, password)
	err = db.Open()
	if err != nil {
		logger.Error.Println("Open database failed!")
		logger.Error.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// Get subscriptions
	subs, err := db.Subscriptions.GetAll()
	if err != nil {
		logger.Error.Println("Get all subscriptions failed!")
		logger.Error.Println(err)
		os.Exit(1)
	}

	// Log result
	logger.Info.Printf("Loaded %d subscriptions.\n", len(subs))

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

	logger.Info.Printf("softtube-update needs to update %d (of %d subscriptions).\n", subsThatNeedsUpdate, len(subs))

	// Create a waitgroup to sync the goroutines
	var waitGroup sync.WaitGroup
	waitGroup.Add(subsThatNeedsUpdate)

	count := 0

	for i := 0; i < len(subs); i++ {
		sub := subs[i]

		if sub.NeedsUpdate() {
			if count > maxUpdates {
				logger.Warning.Println("Max number of updates reached!")
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
}

func updateSubscription(subscription *database.Subscription) {
	logger.Info.Printf("Updating channel '%s'.\n", subscription.Name)

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
			logger.Error.Println(errorMessage)
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
			logger.Error.Println(errorMessage)
			continue
		}
	}

	waitGroup.Wait()

	// Mark subscription as updated
	interval, err := getInterval(subscription.Frequency)
	if err != nil {
		logger.Error.Println(err)
	}
	_ = db.Subscriptions.UpdateLastChecked(subscription, interval)
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
		logger.Error.Println(errorMessage)
		return nil
	}

	if contains404(rss) {
		errorMessage := fmt.Sprintf(
			"Channel %s (%s) has been deleted. Please remove channel...",
			subscription.Name, subscription.ID,
		)
		logger.Error.Println(errorMessage)
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
		logger.Error.Println(errorMessage)
		return nil
	}

	// Get videos in the RSS
	videos, err := feed.getVideos()
	if err != nil {
		errorMessage := fmt.Sprintf(
			"Failed to get videos for %s (%s) : %s",
			subscription.Name, subscription.ID, err.Error(),
		)
		logger.Error.Println(errorMessage)
		return nil
	}

	return videos
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
		logger.Error.Println(msg)
		return err
	} else {
		msg := fmt.Sprintf("Inserted video '%s' in database : Success!", video.Title)
		logger.Info.Println(msg)
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
			logger.Error.Println(msg)
			return
		} else {
			msg := fmt.Sprintf("Updated duration for video '%s' : Success!", video.Title)
			logger.Info.Println(msg)
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
			logger.Error.Println(msg)
			return
		} else {
			msg := fmt.Sprintf("Downloaded thumbnail for video '%s': Success!", video.Title)
			logger.Info.Println(msg)
		}
	}()

	return nil
}

func clean(title string) string {
	re := regexp.MustCompile("[^[:ascii:]åäöÅÄÖ]")
	// re := regexp.MustCompile("^\\p{Cc}")

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
