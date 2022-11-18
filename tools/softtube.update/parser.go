package main

import (
	"encoding/xml"
	"fmt"
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

func (f *Feed) parse(rss string) {
	bytes := []byte(rss)
	_ = xml.Unmarshal(bytes, &f)
}

func (f Feed) getVideos() []database.Video {
	var videoList []database.Video
	for i := 0; i < len(f.Entries); i++ {
		var video database.Video
		video.ID = f.Entries[i].ID
		channelId := f.ChannelID
		if strings.Trim(channelId, " \n\t") == "" {
			channelId = f.Entries[i].ChannelID
		}
		video.SubscriptionID = channelId
		video.Title = f.Entries[i].Title
		publishedDate, err := time.Parse(constDateLayout, f.Entries[i].Published)
		if err != nil {
			// TODO : Handle errors
			fmt.Println(err.Error())
		}
		video.Published = localTime(publishedDate)
		videoList = append(videoList, video)
	}
	return videoList
}

func localTime(datetime time.Time) time.Time {
	loc, err := time.LoadLocation("Europe/Stockholm")
	if err != nil {
		panic(err)
	}
	return datetime.In(loc)
}
