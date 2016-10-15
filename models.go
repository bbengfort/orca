package orca

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

/////////////////////////////////////////////////////////////////////////////
// Model Definitions
/////////////////////////////////////////////////////////////////////////////

// Model specifies types that can interact with the database.
type Model interface {
	Get(id int64, db *sql.DB) error            // Populate model fields from the database
	Save(db *sql.DB) (bool, error)             // Insert or update the model
	Delete(db *sql.DB) (bool, error)           // Delete the model from the database
	Exists(id int64, db *sql.DB) (bool, error) // Determine if the model exists by ID
}

// ModelMeta specifies the fields that all models should have via embedding.
type ModelMeta struct {
	ID      int64     // Unique ID of the model
	Created time.Time // Datetime the model was added to the database
	Updated time.Time // Datetime the model was updated in the database
}

// Device is an entity that represents nodes in the network that can be pinged.
// Device objects are stored in the devices table.
type Device struct {
	Name   string // Hostname of the device
	IPAddr string // IP Address of the device
	Domain string // Domain name of the device
	ModelMeta
}

// Location is a geographic record that is usually associated with an IP
// address via the geoip lookup service but could also come from GPS.
type Location struct {
	IPAddr       string  // IP Address associated with the location
	Latitude     float64 // Decimal based latitude
	Longitude    float64 // Decimal based logitude
	City         string  // City returned by MaxMind for the IP address
	PostCode     string  // Postal code returned by MaxMind for the IP address
	Country      string  // Country returned by MaxMind for the IP address
	Organization string  // Organization associated with the given domain (ISP)
	Domain       string  // Domain associated with the IP address (ISP)
	ModelMeta
}

/////////////////////////////////////////////////////////////////////////////
// Device Methods
/////////////////////////////////////////////////////////////////////////////

// Get a device from the database by ID and populate the struct fields.
func (d *Device) Get(id int64, db *sql.DB) error {
	row := db.QueryRow("SELECT * FROM devices WHERE id = $1", id)
	err := row.Scan(&d.ID, &d.Name, &d.IPAddr, &d.Domain, &d.Created, &d.Updated)

	return err
}

// Save a device struct to the database. This function checks if the device
// has an ID or not. If it does, it will execute a SQL UPDATE, otherwise it
// will execute a SQL INSERT. Returns a boolean if the device was inserted.
// This method handles setting the created and updated timestamps as well.
func (d *Device) Save(db *sql.DB) (bool, error) {
	if d.ID > 0 {
		// This is the UPDATE method so return false.
		// Update the updated timestamp on the device.
		d.Updated = time.Now()

		// Execute the query against the database
		query := "UPDATE devices SET name=$1, ipaddr=$2, domain=$3, updated=$4 WHERE id = $6"
		_, err := db.Exec(query, d.Name, d.IPAddr, d.Domain, d.Updated, d.ID)

		return false, err
	}

	// This is the INSERT method, so return true
	// Set the created and updated timestamps on the device
	d.Created = time.Now()
	d.Updated = time.Now()

	// Execute the INSERT query against the dtabase
	query := "INSERT INTO devices (name, ipaddr, domain, created, updated) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	row := db.QueryRow(query, d.Name, d.IPAddr, d.Domain, d.Created, d.Updated)
	err := row.Scan(&d.ID)

	if err != nil {
		return false, err
	}

	return true, err
}

// Delete a device from the database. Returns true if the number of rows
// affected is 1 or false otherwise.
func (d *Device) Delete(db *sql.DB) (bool, error) {
	return deleteFromDatabase(db, "devices", d.ID)
}

// Exists checks if the specified location is in the database.
func (d *Device) Exists(id int64, db *sql.DB) (bool, error) {
	if id == 0 {
		id = d.ID
	}
	return existsInDatabase(db, "locations", id)
}

/////////////////////////////////////////////////////////////////////////////
// Location Methods
/////////////////////////////////////////////////////////////////////////////

// String returns a pretty representation of the location
func (loc *Location) String() string {
	output := fmt.Sprintf("%s is located at %s, %s (%f, %f)", loc.IPAddr, loc.City, loc.Country, loc.Latitude, loc.Longitude)
	if loc.Organization != "" {
		output += fmt.Sprintf("\nOrganization: %s", loc.Organization)
		if loc.Domain != "" {
			output += fmt.Sprintf(" (%s)", loc.Domain)
		}
	}
	return output
}

// Get a location from the database by ID and populate the struct fields.
func (loc *Location) Get(id int64, db *sql.DB) error {
	row := db.QueryRow("SELECT * FROM locations WHERE id = $1", id)
	err := row.Scan(
		&loc.ID, &loc.IPAddr, &loc.Latitude, &loc.Longitude,
		&loc.City, &loc.PostCode, &loc.Country, &loc.Organization,
		&loc.Domain, &loc.Created, &loc.Updated,
	)

	return err
}

// Save a location struct to the database. This function checks if the
// location has an ID or not. If it does, it will execute a SQL UPDATE,
// otherwise it will execute a SQL INSERT. Returns a boolean if the location
// was inserted. This method handles setting the meta timestamps as well.
func (loc *Location) Save(db *sql.DB) (bool, error) {
	if loc.ID > 0 {
		// This is the UPDATE method so return false.
		// Update the updated timestamp on the device.
		loc.Updated = time.Now()

		// Execute the query against the database
		query := "UPDATE locations SET ipaddr=$1, latitude=$2, longitude=$3, city=$4, postcode=$5, country=$6, organization=$7, domain=$8, updated=$9 WHERE id = $10"
		_, err := db.Exec(query, loc.IPAddr, loc.Latitude, loc.Longitude, loc.City, loc.PostCode, loc.Country, loc.Organization, loc.Domain, loc.Updated, loc.ID)

		return false, err
	}

	// This is the INSERT method, so return true
	// Set the created and updated timestamps on the device
	loc.Created = time.Now()
	loc.Updated = time.Now()

	// Execute the INSERT query against the dtabase
	query := "INSERT INTO locations (ipaddr, latitude, longitude, city, postcode, country, organization, domain, created, updated) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	row := db.QueryRow(query)
	err := row.Scan(&loc.ID)

	if err != nil {
		return false, err
	}

	return true, err
}

// Delete a location from the database. Returns true if the number of rows
// affected is 1 or false otherwise.
func (loc *Location) Delete(db *sql.DB) (bool, error) {
	return deleteFromDatabase(db, "locations", loc.ID)
}

// Exists checks if the specified location is in the database.
func (loc *Location) Exists(id int64, db *sql.DB) (bool, error) {
	if id == 0 {
		id = loc.ID
	}
	return existsInDatabase(db, "locations", id)
}

/////////////////////////////////////////////////////////////////////////////
// Model Helper Functions
/////////////////////////////////////////////////////////////////////////////

// Helper function that deletes an item from a table by ID.
func deleteFromDatabase(db *sql.DB, table string, id int64) (bool, error) {
	if id == 0 {
		msg := fmt.Sprintf("Cannot delete a row with id=0 from the %s table", table)
		return false, errors.New(msg)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", table)
	res, err := db.Exec(query, id)

	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()

	switch {
	case err != nil:
		return false, err
	case rows > 1:
		return false, errors.New("Multiple deletions from the database?!")
	case rows == 1:
		return true, nil
	case rows < 1:
		return false, nil
	default:
		return false, errors.New("Unknown case in device deletion")
	}
}

// Helper function that checks if a row exists in a table by ID
func existsInDatabase(db *sql.DB, table string, id int64) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1 LIMIT 1)", table)
	row := db.QueryRow(query, id)
	err := row.Scan(&exists)

	return exists, err
}
