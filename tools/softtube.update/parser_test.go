package main

import (
	"testing"
)

const testKarenPuzzlesChannelId = "UC2umy62ojMfxzzHkVcgEUUA"

func TestFeed_parse(t *testing.T) {
	type args struct {
		channelId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want404 bool
		wantErr bool
	}{
		{"empty", args{""}, "", false, true},
		{"invalid", args{"xxxxxxxxxxx"}, "", true, false},
		{"karen puzzles", args{testKarenPuzzlesChannelId}, "", false, false},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				y := youtube{}
				rss, err := y.getSubscriptionRSS(tt.args.channelId)
				if (err != nil) != tt.wantErr {
					t.Errorf("getSubscriptionRSS() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if rss == "" {
					return
				}

				if tt.want404 {
					if !contains404(rss) {
						t.Errorf("getSubscriptionRSS() should return 404 for invalid channelId")
						return
					}
				} else {
					f := &Feed{}
					if err = f.parse(rss); (err != nil) != tt.wantErr {
						t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
					}

					if f.Title != "Karen Puzzles" {
						t.Errorf("feed title = %v, want %v", f.Title, "Karen Puzzles")
					}
				}
			},
		)
	}
}

func TestFeed_getVideos(t *testing.T) {
	type args struct {
		channelId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"karen puzzles", args{testKarenPuzzlesChannelId}, "", false},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				y := youtube{}
				rss, err := y.getSubscriptionRSS(tt.args.channelId)
				if (err != nil) != tt.wantErr {
					t.Errorf("getSubscriptionRSS() returned error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				f := &Feed{}
				if err = f.parse(rss); (err != nil) != tt.wantErr {
					t.Errorf("parse() returned error = %v, wantErr %v", err, tt.wantErr)
				}

				videos, err := f.getVideos()
				if err != nil {
					t.Errorf("getVideos() return error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if len(videos) != 15 {
					t.Errorf("getSubscriptionRSS() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				for i, video := range videos {
					if video.SubscriptionID != testKarenPuzzlesChannelId {
						t.Errorf(
							"video(%d) has wrong channel id (%v), want (%v)",
							i, video.SubscriptionID, testKarenPuzzlesChannelId,
						)
						return
					}
				}
			},
		)
	}
}
