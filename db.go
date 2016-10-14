package orca

import (
	"database/sql"
	"errors"

	// Imports the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

// ConnectDB establishes a connection to the Sqlite3 database
func (app *App) ConnectDB() error {
	if app.db != nil {
		return errors.New("A database connection already exists!")
	}

	var err error
	app.db, err = sql.Open("sqlite3", app.Config.DBPath)
	return err
}

// CreateDB  executes the create table statements in the schema.sql stored as
// binary data in the application (as well as any alter table statements).
func (app *App) CreateDB() error {

	// Load the schema from the binary data
	schema, err := Asset("fixtures/schema.sql")
	if err != nil {
		return err
	}

	// Execute the schema SQL
	_, err = app.db.Exec(string(schema))
	if err != nil {
		return err
	}

	return nil
}
