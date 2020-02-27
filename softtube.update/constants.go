package main

const configPath string = "/.config/softtube/softtube.config"

const videoDurationCommand = "youtube-dl --get-duration -- '%s'"

const thumbnailCommand string = "youtube-dl --write-thumbnail --skip-download --no-overwrites -o '%s' -- '%s'"
const thumbnailPath string = "%s/%s.jpg"

const subscriptionRSSURL = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"

const dateLayout = "2006-01-02T15:04:05-07:00"
