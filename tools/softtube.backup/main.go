package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	core "github.com/hultan/softtube/internal/softtube.core"
)

const (
	constDateLayoutBackup = "20060102_0304"
)

func main() {
	// Init config file
	config := new(core.Config)
	err := config.Load("main")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = backup("softtube", config.ServerPaths.Backup)
	// backup("softtubeTEST", config.ServerPaths.Backup)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

// Backs up a mysql database
func backup(database, rootBackupPath string) error {
	backupFile := fmt.Sprintf("%s_%s.sql", database, time.Now().Local().Format(constDateLayoutBackup))
	backupPath := path.Join(rootBackupPath, backupFile)
	command := fmt.Sprintf("mysqldump -u per %s > %s", database, backupPath)

	_, err := exec.Command("/bin/bash", "-c", command).Output()
	if err != nil {
		return err
	}

	return nil
}
