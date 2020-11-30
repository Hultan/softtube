package main

// See https://github.com/ytdl-org/youtube-dl/issues/22641
// It might help adding --user-agent "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
const constVideoDurationCommand = "%s --get-duration -- '%s'"

const constThumbnailCommand string = "%s --write-thumbnail --skip-download --no-overwrites -o '%s' -- '%s'"
const constThumbnailLocation string = "%s/%s.jpg"

const constSubscriptionRSSURL = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"

const constDateLayout = "2006-01-02T15:04:05-07:00"
