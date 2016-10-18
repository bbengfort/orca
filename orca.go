// Package orca provides the library for a systems experiment that measures
// latency and uptime of mobile nodes against fixed responder nodes.
package orca

import "database/sql"

// Version specifies the current version of the Orca library.
const Version = "0.1"

// App is the primary application that maintains references to the config
// and device details as well as initializes the environment and runs the
// reflect and generate commands.
type App struct {
	Config     *Config        // The configuration loaded from the YAML file
	GeoIP      *MaxMindClient // GeoIP Lookup API client
	Location   *Location      // Current location of the application
	ExternalIP string         // Current external IP address of the machine
	db         *sql.DB        // Connection to the database stored on the app
}

// Init the orca application
// NOTE: SyncLocation should not be called in init!
func Init() (*App, error) {
	var err error
	app := new(App)

	// Load the configuration from the YAML files
	app.Config = LoadConfig()

	// Connect to the database
	if err = app.ConnectDB(); err != nil {
		return nil, err
	}

	// Initialize the MaxMindClient for GeoIP lookup
	app.GeoIP = NewMaxMindClient(
		app.Config.MaxMind.Username, app.Config.MaxMind.License,
	)

	return app, nil
}

// SyncLocation checks the external IP address against the current IP address,
// if they're different then it performs another location lookup to track
// mobility in the generator application, but does not perform GeoIP lookups
// if they're not necessary (to save bandwidth and cost).
func (app *App) SyncLocation() error {

	// Get the external IP address
	eip, err := ExternalIP()
	if err != nil {
		return err
	}

	// Compare to current IP address and if different, fetch new location.
	if eip != app.ExternalIP {
		// Store the current external IP address on the app
		app.ExternalIP = eip

		// Initialize the current location for geographic tracking
		loc, err := app.GeoIP.GetCurrentLocation()
		if err != nil {
			return err
		}

		// Set the location on the app and save to database
		if err = app.SetLocation(loc, true); err != nil {
			return err
		}
	}

	return nil
}

// SetLocation is a wrapper method that sets the location on the app struct,
// but also does a check about whether or not to save it to the database.
func (app *App) SetLocation(loc *Location, save bool) error {
	if save {
		// Check to make sure the location isn't in the database.
		loc.IPExists(app.db)

		if _, err := loc.Save(app.db); err != nil {
			return err
		}
	}

	// Store the location with the application
	app.Location = loc
	return nil
}

// GetListenAddr looks up the IP address on the config or gets an external IP
// This is meant to be used by reflect mode to respond to echo requests.
func (app *App) GetListenAddr() (string, error) {
	return ResolveAddr(app.Config.Addr)
}
