package main

import (
	database "github.com/hultan/softtube/internal/softtube.database"
)

const constStatusDownloading database.VideoStatusType = 1
const constStatusDownloaded database.VideoStatusType = 2

const (
	errorOpenConfig   = 1
	errorOpenDatabase = 2
)
