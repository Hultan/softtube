package main

const constAppID = "com.github.hultan.softtube"
const constAppVersion = "0.1"
const constDateLayout = "2006-01-02"

// LISTSTORE COLUMNS
const liststoreColumnImage = 0
const liststoreColumnChannelName = 1
const liststoreColumnDate = 2
const liststoreColumnTitle = 3
const liststoreColumnProgress = 4
const liststoreColumnBackground = 5
const liststoreColumnVideoID = 6
const liststoreColumnDuration = 7
const liststoreColumnProgressText = 8
const liststoreColumnForeground = 9

const constColorNotDownloaded = "#444444"
const constColorDownloading = "Dodger Blue"
const constColorDownloaded = "Dark Slate Blue"
const constColorWatched = "Dark Green"
const constColorDeleted = "Dark Slate Gray"
const constColorWarning = "Coral"
const constColorSaved = "Dark cyan"

const constStatusNotDownloaded = 0
const constStatusDownloading = 1
const constStatusDownloaded = 2
const constStatusWatched = 3
const constStatusDeleted = 4
const constStatusSaved = 5          // Not used yet
const constStatusDownloadFailed = 6 // Not used yet

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
