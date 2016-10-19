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

// GetDB returns the database connectio object from the app.
// Use with care, I didn't want to expose this outside the app ...
func (app *App) GetDB() *sql.DB {
	if app.db == nil {
		app.CreateDB()
	}
	return app.db
}

// FetchDevices returns a collection of devices, ordered by the created
// timestamp. This function expects you to limit the size of the collection
// by specifying the maximum number of nodes to return in the Devices list.
func (app *App) FetchDevices() (Devices, error) {
	var devices Devices

	rows, err := app.db.Query("SELECT * FROM devices ORDER BY created DESC")
	if err != nil {
		return devices, err
	}

	for rows.Next() {
		var d Device
		if err := rows.Scan(&d.ID, &d.Name, &d.IPAddr, &d.Domain, &d.Sequence, &d.Created, &d.Updated); err != nil {
			return devices, err
		}

		devices = append(devices, d)
	}

	rows.Close()
	return devices, nil
}
