package softtube

import (
	"fmt"
	"os"
	"slices"

	"github.com/hultan/dialog"
)

type Stats struct {
	Disk     []string
	Database []string
}

func (s *SoftTube) showStats() {
	onDisk, err := s.collectDiskStats()
	if err != nil {
		_, _ = dialog.Title("Failed to collect disk stats").
			Text("An error occurred while trying to collect disk stats:").
			ExtraExpand(err.Error()).ExtraHeight(100).
			Width(500).ErrorIcon().OkButton().Show()
	}
	inDb, err := s.collectDBStats()
	if err != nil {
		_, _ = dialog.Title("Failed to collect DB stats").
			Text("An error occurred while trying to collect DB stats:").
			ExtraExpand(err.Error()).ExtraHeight(100).
			Width(500).ErrorIcon().OkButton().Show()
	}

	if len(onDisk) == len(inDb) {
		return
	}

	diff := getDiff(onDisk, inDb)

	stats := fmt.Sprintf("Files that differ:\n")
	stats += fmt.Sprintf("============\n")
	for i := 0; i < len(diff); i++ {
		stats += fmt.Sprintf("%s\n", diff[i])
	}
	stats += fmt.Sprintf("============\n")
	stats += fmt.Sprintf("Disk: %d   DB: %d", len(onDisk), len(inDb))

	text := "SoftTube missing files and missing database entries:"
	_, _ = dialog.Title("SoftTube statistics").
		Text(text).
		ExtraExpand(stats).
		WarningIcon().
		Size(400, 250).
		OkButton().
		Show()
}

func (s *SoftTube) collectDiskStats() ([]string, error) {
	dir, err := os.Open(s.Config.ClientPaths.Videos)
	if err != nil {
		return nil, err
	}
	files, err := dir.Readdirnames(0)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *SoftTube) collectDBStats() ([]string, error) {
	files, err := s.DB.Videos.GetStats()
	if err != nil {
		return nil, err
	}
	return files, nil
}

func getDiff(disk []string, db []string) []string {
	var missing []string
	for _, file := range disk {
		if !slices.Contains(db, file) {
			missing = append(missing, fmt.Sprintf("%s (not in DB)", file))
		}
	}
	for _, file := range db {
		if !slices.Contains(disk, file) {
			missing = append(missing, fmt.Sprintf("%s (not on disk)", file))
		}
	}
	return missing
}
