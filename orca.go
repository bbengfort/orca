// Package orca provides the library for a systems experiment that measures
// latency and uptime of mobile nodes against fixed responder nodes.
package orca

// Version specifies the current version of the Orca library.
const Version = "0.1"

// App is the primary application that maintains references to the config
// and device details as well as initializes the environment and runs the
// reflect and generate commands.
type App struct {
	Config *Config // The configuration loaded from the YAML file

}

// Init the orca application
func Init() (*App, error) {
	app := new(App)

	// Load the configuration from the YAML files
	app.Config = LoadConfig()

	return app, nil
}

// GetAddr looks up the IP address on the config or gets an external IP
func (app App) GetAddr() (string, error) {
	return ResolveAddr(app.Config.Addr)
}
