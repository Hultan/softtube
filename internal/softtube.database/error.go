package database

import (
	"errors"
	"fmt"
)

type ErrVideo struct {
	Message string
	VideoID string
}

func (e ErrVideo) Error() string {
	return fmt.Sprintf("%s. VideoID: %s\n\n%s\n", e.Message, e.VideoID)
}

func newErrVideo(message string, videoID string) error {
	return ErrVideo{
		Message: message,
		VideoID: videoID,
	}
}

var ErrDatabaseNotOpened = errors.New("database not opened")
