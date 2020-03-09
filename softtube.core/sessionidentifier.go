package core

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// SessionIdentifier : Uniquely identifies the computer running this session
type SessionIdentifier struct {
	MachineID string
	Name      string
}

// CreateSession : Creates a new session and stores
// it in the database if needed
func CreateSession(db Database) (SessionIdentifier, error) {
	s := getSession()
	if s.MachineID == "" || s.Name == "" {
		return s, errors.New("invalid session")
	}

	err := storeSession(s, db)
	if err != nil {
		return s, err
	}

	db.SessionIdentifier = s

	return s, nil
}

func storeSession(s SessionIdentifier, db Database) error {
	err := db.Session.Insert(s)
	if err != nil {
		return err
	}
	return nil
}

func getSession() SessionIdentifier {
	var s SessionIdentifier
	machineID := getFileContents("/etc/machine-id")
	if len(machineID) > constMachineIDMaxLength {
		machineID = machineID[:constMachineIDMaxLength]
	}
	s.MachineID = machineID

	computer := getFileContents("/etc/hostname")
	if len(computer) > constHostNameMaxLength {
		computer = computer[:constHostNameMaxLength]
	}
	s.Name = computer

	return s
}

func getFileContents(path string) string {
	command := fmt.Sprintf("cat %s", path)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.Trim(string(output), "\n ")
}
