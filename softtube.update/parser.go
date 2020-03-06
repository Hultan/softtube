package main

import (
	"encoding/xml"
	"fmt"
	"time"

	core "github.com/hultan/softtube/softtube.core"
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
	Title     string `xml:"title"`
	Published string `xml:"published"`
}

func (f *Feed) parse(rss string) {
	bytes := []byte(rss)
	xml.Unmarshal(bytes, &f)
}

func (f Feed) getVideos() []core.Video {
	var videoList []core.Video
	for i := 0; i < len(f.Entries); i++ {
		var video core.Video
		video.ID = f.Entries[i].ID
		video.SubscriptionID = f.ChannelID
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
