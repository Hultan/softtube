package main

import (
	"fmt"
	"os"
	"testing"

	core "github.com/hultan/softtube/internal/softtube.core"
)

const (
	ytdlpPath          = "/usr/local/bin/yt-dlp"
	testVideoId        = "K7TiNRdHV84"
	testFailedDuration = "ERROR: [youtube] xxxxxxxxxxx: Video unavailable\n"
)

func initConfig(t *testing.T) {
	// Init config file
	config = new(core.Config)
	err := config.Load("main")
	if err != nil {
		t.Errorf("failed to init config : %v", err)
	}
}

func Test_youtube_getYoutubePath(t *testing.T) {
	initConfig(t)
	yt := &youtube{}
	got := yt.getYoutubePath()
	if got != ytdlpPath {
		t.Errorf("getYoutubePath() = %v, want %v", got, ytdlpPath)
	}
}

func Test_youtube_getDurationInternal(t *testing.T) {
	type args struct {
		videoID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{""}, "", true},
		{"invalid", args{"xxxxxxxxxxx"}, testFailedDuration, true},
		{"karen puzzles", args{testVideoId}, "10:02", false},
		{"python warning", args{"F7AEii-r4R8"}, "26:24", false},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				initConfig(t)

				y := youtube{}
				got, err := y.getDurationInternal(tt.args.videoID)
				if (err != nil) != tt.wantErr {
					t.Errorf("getDurationInternal() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("getDurationInternal() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func Test_youtube_getThumbnailInternal(t *testing.T) {
	type args struct {
		videoID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{""}, "", true},
		{"invalid", args{"xxxxxxxxxxx"}, "", true},
		{"karen puzzles", args{testVideoId}, "x", false},
	}

	initConfig(t)
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// First clean up, in case last test failed
				path := fmt.Sprintf("/softtube/thumbnails/%s.webp", testVideoId)
				_ = os.Remove(path)

				y := youtube{}
				_, err := y.getThumbnailInternal(tt.args.videoID)

				if (err != nil) != tt.wantErr {
					t.Errorf("getThumbnailInternal() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.want != "" {
					// Check that file exists
					if _, err = os.Stat(path); os.IsNotExist(err) {
						t.Errorf("getThumbnailInternal() failed to download thumbnail")
					}
					// Cleanup after us
					err = os.Remove(path)
					if err != nil {
						t.Errorf("getThumbnailInternal() failed to delete thumbnail")
					}
				}
			},
		)
	}
}

func Test_youtube_getSubscriptionRSS(t *testing.T) {
	type args struct {
		channelID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want404 bool
		wantErr bool
	}{
		{"empty", args{""}, "", false, true},
		{"invalid", args{"UC2umy62ojMfxzzHkVcgxxxx"}, "", true, false},
		{"karen puzzles", args{"UC2umy62ojMfxzzHkVcgEUUA"}, "x", false, false},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				y := youtube{}
				got, err := y.getSubscriptionRSS(tt.args.channelID)
				if (err != nil) != tt.wantErr {
					t.Errorf("getSubscriptionRSS() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.want404 {
					if !contains404(got) {
						t.Errorf("getSubscriptionRSS() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				} else if tt.want == "" {
					if got != tt.want {
						t.Errorf("getSubscriptionRSS() got = %v, want %v", got, tt.want)
					}
				}
			},
		)
	}
}
