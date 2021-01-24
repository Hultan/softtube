package main

//const constAppID = "com.github.hultan.softtube"
const constAppTitle = "SoftPlan"
const constAppVersion = "2.0"
const constDateLayout = "2006-01-02"

// LISTSTORE COLUMNS

const listStoreColumnImage = 0
const listStoreColumnChannelName = 1
const listStoreColumnDate = 2
const listStoreColumnTitle = 3
const listStoreColumnProgress = 4
const listStoreColumnBackground = 5
const listStoreColumnVideoID = 6
const listStoreColumnDuration = 7
const listStoreColumnProgressText = 8
const listStoreColumnForeground = 9

const constColorNotDownloaded = "#444444"
const constColorDownloading = "Dodger Blue"
const constColorDownloaded = "Dark Slate Blue"
const constColorWatched = "Dark Green"
const constColorDeleted = "Dark Slate Gray"
const constColorWarning = "Coral"
const constColorLive = "Goldenrod"
const constColorSaved = "Dark cyan"

const constStatusNotDownloaded = 0
const constStatusDownloading = 1
const constStatusDownloaded = 2
const constStatusWatched = 3
const constStatusDeleted = 4
const constStatusSaved = 5

const constFilterModeSubscriptions = 0
const constFilterModeToWatch = 1
const constFilterModeToDelete = 2
const constFilterModeSaved = 3

const constLogDownload = 0
const constLogPlay = 1
const constLogDelete = 2
const constLogSetWatched = 3
const constLogSetUnwatched = 4
const constLogError = 5

const constSetAsSaved = "Set as saved"
const constSetAsNotSaved = "Set as not saved"
const constSetAsWatched = "Set as watched"
const constSetAsUnwatched = "Set as unwatched"
const constSetAsNotDownloaded = "Set as not downloaded"

const constThumbnailCommand string = "%s --write-thumbnail --skip-download --no-overwrites -o '%s' -- '%s'"
const constThumbnailLocation string = "%s/%s"