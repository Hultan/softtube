package main

import (
	"encoding/xml"
	"errors"
	"strings"
	"time"

	database "github.com/hultan/softtube/internal/softtube.database"
)

// Feed : A subscription RSS feed
type Feed struct {
	Title     string  `xml:"title"`
	ChannelID string  `xml:"channelId"`
	Entries   []Entry `xml:"entry"`
}

// Entry : a youtube video
type Entry struct {
	ID        string `xml:"videoId"`
	ChannelID string `xml:"channelId"`
	Title     string `xml:"title"`
	Published string `xml:"published"`
}

func (f *Feed) parse(rss string) error {
	if rss == "" {
		return errors.New("rss cannot be empty")
	}
	bytes := []byte(rss)
	return xml.Unmarshal(bytes, &f)
}

func (f *Feed) getVideos() ([]database.Video, error) {
	var videoList []database.Video

	channelId := f.getChannelId() // Set the channel ID
	for _, entry := range f.Entries {
		var video database.Video
		video.ID = entry.ID
		video.SubscriptionID = channelId
		video.Title = entry.Title
		date, err := f.getDate(entry.Published)
		if err != nil {
			return nil, err
		}
		video.Published = date
		videoList = append(videoList, video)
	}

	return videoList, nil
}

func (f *Feed) getDate(dateString string) (time.Time, error) {
	// Parse the date string
	publishedDate, err := time.Parse(constDateLayout, dateString)
	if err != nil {
		return time.Now(), err
	}
	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		return time.Now(), err
	}
	return publishedDate.In(loc), nil
}

func (f *Feed) getChannelId() string {
	channelId := f.ChannelID // Set the channel ID
	for _, entry := range f.Entries {
		// 221120 PH : YouTube have made a change, sometimes the
		// channel ID does not exist in the feed, but it does
		// exist in one of the entries.
		if strings.Trim(channelId, " \n\t") == "" {
			channelId = entry.ChannelID
		}
	}
	return channelId
}
