// Package orca provides the library for a systems experiment that measures
// latency and uptime of mobile nodes against fixed responder nodes.
package orca

// Version specifies the current version of the Orca library.
const Version = "0.1"

// App is the primary application that maintains references to the config
// and device details as well as initializes the environment and runs the
// reflect and generate commands.
type App struct {
	Silent    bool      // Whether to log command line messages or not
	IPAddr    string    // Current external IP address of the machine
	Latencies []float64 // List of latencies to compute metrics on
}

// Init the orca application
// NOTE: SyncLocation should not be called in init!
func Init(silent bool) (*App, error) {
	app := new(App)
	app.Silent = silent

	ipaddr, err := app.GetListenAddr()
	if err != nil {
		return nil, err
	}

	app.IPAddr = ipaddr
	return app, nil
}

// GetListenAddr looks up the IP address on the config or gets an external IP
// This is meant to be used by reflect mode to respond to echo requests.
func (app *App) GetListenAddr() (string, error) {
	return ResolveAddr(app.IPAddr)
}
