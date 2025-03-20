package softtube

import (
	database "github.com/hultan/softtube/internal/softtube.database"
)

const constAppTitle = "SoftTube"
const constAppVersion = "3.0.5"
const constDateLayout = "2006-01-02"

type listStoreColumnType int

const (
	listStoreColumnImage listStoreColumnType = iota
	listStoreColumnChannelName
	listStoreColumnDate
	listStoreColumnTitle
	listStoreColumnProgress
	listStoreColumnBackground
	listStoreColumnVideoID
	listStoreColumnDuration
	listStoreColumnProgressText
	listStoreColumnForeground
)

type colorType string

const (
	constColorNotDownloaded colorType = "#444444"
	constColorDownloading             = "Dodger Blue"
	constColorDownloaded              = "Dark Slate Blue"
	constColorWatched                 = "Dark Green"
	constColorDeleted                 = "Dark Slate Gray"
	constColorWarning                 = "Coral"
	constColorLive                    = "Goldenrod"
	constColorSaved                   = "Dark cyan"
)

const (
	constStatusNotDownloaded database.VideoStatusType = iota
	constStatusDownloading
	constStatusDownloaded
	constStatusWatched
	constStatusDeleted
	constStatusSaved
)

type viewType int

const (
	viewSubscriptions viewType = iota + 1
	viewDownloads
	viewToWatch
	viewSaved
	viewToDelete
)

const (
	constLogDownload database.LogType = iota
	constLogPlay
	constLogDelete
	constLogSetWatched
	constLogSetUnwatched
	constLogError
)

const (
	constSetAsSaved         = "Set as saved"
	constSetAsNotSaved      = "Set as not saved"
	constSetAsWatched       = "Set as watched"
	constSetAsUnwatched     = "Set as unwatched"
	constSetAsNotDownloaded = "Set as not downloaded"
)

const constThumbnailCommand string = "%s --write-thumbnail --skip-download --no-overwrites -o '%s' -- '%s'"
const constThumbnailLocation string = "%s/%s"
