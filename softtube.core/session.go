package core

import (
	"database/sql"
	"errors"
)

// SessionTable : Handles Session database table
type SessionTable struct {
	Connection *sql.DB
}

const sqlStatementInsertSession = `INSERT IGNORE INTO Session (id, name) 
								VALUES (?, ?);`

// Insert : Creates a session in the Session table of the database
func (s SessionTable) Insert(id SessionIdentifier) error {
	// Check that database is opened
	if s.Connection == nil {
		return errors.New("database not opened")
	}

	// Execute insert statement
	_, err := s.Connection.Exec(sqlStatementInsertSession, id.MachineID, id.Name)
	if err != nil {
		return err
	}

	return nil
}
