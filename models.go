package orca

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bbengfort/orca/echo"
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
	Name     string       // Hostname of the device
	IPAddr   string       // IP Address of the device
	Domain   string       // Domain name of the device
	Sequence int64        // The response/reply counter for a device
	echo     *echo.Device // The protocol buffer representation
	ModelMeta
}

// Devices is a collection of Device objects loaded from the database.
type Devices []*Device

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
	Note         string  // Any additional annotations by the user
	ModelMeta
}

// Ping is a timeseries record of latency requests reflected from echo servers.
type Ping struct {
	ID       int64           //  Unique ID of the record
	Source   *Device         // Source device of the ping (always the local node)
	Target   *Device         // Target device that the ping was sent to
	Location *Location       // The location of the source at the time of the ping
	Request  int64           // Request sequence number for the source/target pair
	Response int64           // Response sequence number for the target/source pair
	Sent     time.Time       // The time that the ping was sent
	Recv     time.Time       // The time that the ping was received
	Latency  sql.NullFloat64 // The latency in milliseconds of the ping
}

/////////////////////////////////////////////////////////////////////////////
// Device Methods
/////////////////////////////////////////////////////////////////////////////

// Get a device from the database by ID and populate the struct fields.
func (d *Device) Get(id int64, db *sql.DB) error {
	row := db.QueryRow("SELECT * FROM devices WHERE id = $1", id)
	err := row.Scan(&d.ID, &d.Name, &d.IPAddr, &d.Domain, &d.Sequence, &d.Created, &d.Updated)

	return err
}

// GetByName a device from the database and populate the struct fields.
func (d *Device) GetByName(name string, db *sql.DB) error {

	row := db.QueryRow("SELECT * FROM devices WHERE name = $1", name)
	err := row.Scan(&d.ID, &d.Name, &d.IPAddr, &d.Domain, &d.Sequence, &d.Created, &d.Updated)
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
		query := "UPDATE devices SET name=$1, ipaddr=$2, domain=$3, sequence=$4, updated=$5 WHERE id = $6"
		_, err := db.Exec(query, d.Name, d.IPAddr, d.Domain, d.Sequence, d.Updated, d.ID)

		return false, err
	}

	// This is the INSERT method, so return true
	// Set the created and updated timestamps on the device
	d.Created = time.Now()
	d.Updated = time.Now()

	// Create the query to insert the device into the database
	query := "INSERT INTO devices (name, ipaddr, domain, sequence, created, updated) VALUES ($1, $2, $3, $4, $5, $6)"

	// Execute the INSERT query against the dtabase
	res, err := db.Exec(query, d.Name, d.IPAddr, d.Domain, d.Sequence, d.Created, d.Updated)
	if err != nil {
		return false, err
	}

	// Look up the last inserted ID from sqlite3
	did, err := res.LastInsertId()
	if err != nil {
		return false, err
	}

	// Store the ID and return
	d.ID = did
	return true, err
}

// Delete a device from the database. Returns true if the number of rows
// affected is 1 or false otherwise.
func (d *Device) Delete(db *sql.DB) (bool, error) {
	return deleteFromDatabase(db, "devices", d.ID)
}

// Exists checks if the specified device is in the database.
func (d *Device) Exists(id int64, db *sql.DB) (bool, error) {
	if id == 0 {
		id = d.ID
	}
	return existsInDatabase(db, "devices", id)
}

// Echo converts a device to an echo.Device protocol buffer message.
func (d *Device) Echo() *echo.Device {
	if d.echo == nil {
		d.echo = new(echo.Device)
		d.echo.Name = d.Name
		d.echo.IPAddr = d.IPAddr
		d.echo.Domain = d.Domain
	}

	return d.echo
}

// String returns the textual string representation of the device
func (d *Device) String() string {
	output := "%s (%s)"

	var addr string
	switch {
	case d.Domain != "":
		addr = d.Domain
	case d.IPAddr != "":
		addr = d.IPAddr
	default:
		addr = "Unknown Address"
	}

	return fmt.Sprintf(output, d.Name, addr)
}

/////////////////////////////////////////////////////////////////////////////
// Location Methods
/////////////////////////////////////////////////////////////////////////////

// Get a location from the database by ID and populate the struct fields.
func (loc *Location) Get(id int64, db *sql.DB) error {
	row := db.QueryRow("SELECT * FROM locations WHERE id = $1", id)
	err := row.Scan(
		&loc.ID, &loc.IPAddr, &loc.Latitude, &loc.Longitude,
		&loc.City, &loc.PostCode, &loc.Country, &loc.Organization,
		&loc.Domain, &loc.Note, &loc.Created, &loc.Updated,
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
		query := "UPDATE locations SET ipaddr=$1, latitude=$2, longitude=$3, city=$4, postcode=$5, country=$6, organization=$7, domain=$8, note=$9, updated=$10 WHERE id = $11"
		_, err := db.Exec(query, loc.IPAddr, loc.Latitude, loc.Longitude, loc.City, loc.PostCode, loc.Country, loc.Organization, loc.Domain, loc.Note, loc.Updated, loc.ID)

		return false, err
	}

	// This is the INSERT method, so return true
	// Set the created and updated timestamps on the device
	loc.Created = time.Now()
	loc.Updated = time.Now()

	// Construct the query
	query := "INSERT INTO locations (ipaddr, latitude, longitude, city, postcode, country, organization, domain, note, created, updated)"
	query += " VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	// Execute the INSERT query against the dtabase
	res, err := db.Exec(query, loc.IPAddr, loc.Latitude, loc.Longitude, loc.City, loc.PostCode, loc.Country, loc.Organization, loc.Domain, loc.Note, loc.Created, loc.Updated)
	if err != nil {
		return false, err
	}

	// Get the last inserted ID from SQLite3
	lid, err := res.LastInsertId()
	if err != nil {
		return false, err
	}

	// Store the ID on the location and return
	loc.ID = lid
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

// IPExists sets the ID on the location if the location's IP address is
// already in the database, otherwise it sets it to zero.
func (loc *Location) IPExists(db *sql.DB) error {
	query := "SELECT id FROM locations WHERE ipaddr = $1 LIMIT 1"
	row := db.QueryRow(query, loc.IPAddr)
	err := row.Scan(&loc.ID)

	return err
}

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

/////////////////////////////////////////////////////////////////////////////
// Ping Methods
/////////////////////////////////////////////////////////////////////////////

// Get a ping from the database by ID and populate the struct fields.
func (p *Ping) Get(id int64, db *sql.DB) error {

	// Construct the Ping query
	query := "SELECT * from pings p "
	query += "   JOIN devices s on p.source_id = s.id "
	query += "   JOIN devices t on p.target_id = t.id "
	query += "   JOIN locations l on p.location_id = l.id "
	query += "WHERE p.id=$1"

	// Create the empty struct targets
	p.Source = new(Device)
	p.Target = new(Device)
	p.Location = new(Location)

	// Execute the query and scann the ping
	row := db.QueryRow(query, id)
	err := row.Scan(
		&p.ID, &p.Source.ID, &p.Target.ID, &p.Location.ID, &p.Request, &p.Response, &p.Sent, &p.Recv, &p.Latency,
		&p.Source.ID, &p.Source.Name, &p.Source.IPAddr, &p.Source.Domain, &p.Source.Sequence, &p.Source.Created, &p.Source.Updated,
		&p.Target.ID, &p.Target.Name, &p.Target.IPAddr, &p.Target.Domain, &p.Target.Sequence, &p.Target.Created, &p.Target.Updated,
		&p.Location.ID, &p.Location.IPAddr, &p.Location.Latitude, &p.Location.Longitude, &p.Location.City, &p.Location.PostCode,
		&p.Location.Country, &p.Location.Organization, &p.Location.Domain, &p.Location.Note, &p.Location.Created, &p.Location.Updated,
	)

	return err
}

// Save a ping struct to the database. This function checks if the ping
// has an ID or not. If it does, it will execute a SQL UPDATE, otherwise it
// will execute a SQL INSERT. Returns a boolean if the device was inserted.
// This method handles setting the created and updated timestamps as well.
func (p *Ping) Save(db *sql.DB) (bool, error) {
	if p.ID > 0 {
		// This is the UPDATE method so return false.
		// Execute the query against the database
		query := "UPDATE pings SET "
		query += "source_id=$1, target_id=$2, location_id=$3, request=$4, "
		query += "response=$5, sent=$6, recv=$7, latency=$8 "
		query += "WHERE id = $9"
		_, err := db.Exec(query, p.Source.ID, p.Target.ID, p.Location.ID, p.Request, p.Response, p.Sent, p.Recv, p.Latency, p.ID)

		return false, err
	}

	// This is the INSERT method, so return true
	// Create the query to insert the device into the database
	query := "INSERT INTO pings "
	query += "(source_id, target_id, location_id, request, response, sent, recv, latency) "
	query += "VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	// Execute the INSERT query against the dtabase
	res, err := db.Exec(query, p.Source.ID, p.Target.ID, p.Location.ID, p.Request, p.Response, p.Sent, p.Recv, p.Latency)
	if err != nil {
		return false, err
	}

	// Look up the last inserted ID from sqlite3
	pid, err := res.LastInsertId()
	if err != nil {
		return false, err
	}

	// Store the ID and return
	p.ID = pid
	return true, err
}

// Delete a ping from the database. Returns true if the number of rows
// affected is 1 or false otherwise.
func (p *Ping) Delete(db *sql.DB) (bool, error) {
	return deleteFromDatabase(db, "pings", p.ID)
}

// String returns a pretty representation of the ping
func (p *Ping) String() string {
	output := "%s -> %s order=%d seq=%d %0.3fms"
	return fmt.Sprintf(output, p.Source.Name, p.Target.Name, p.Request, p.Response, p.Latency)
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
