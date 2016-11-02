// Package orca provides the library for a systems experiment that measures
// latency and uptime of mobile nodes against fixed responder nodes.
package orca

import (
	"fmt"
	"math"
)

// Version specifies the current version of the Orca library.
const Version = "ping"

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
func Init(silent bool, addr string) (*App, error) {
	app := new(App)
	app.Silent = silent
	app.IPAddr = addr

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

// ComputeStats aggregates the latencies and prints out a description of the
// statistical measurements of the pings recorded.
func (app *App) ComputeStats() {
	// Compute the statistics of the latencies
	num := 0.0
	sum := 0.0
	ssq := 0.0
	min := -1.0
	max := -1.0

	for _, latency := range app.Latencies {
		num += 1.0
		sum += latency

		if min < 0.0 || latency < min {
			min = latency
		}

		if max < 0.0 || latency > max {
			max = latency
		}
	}

	mean := sum / num

	for _, latency := range app.Latencies {
		ssq += math.Pow(latency-mean, 2)
	}

	sd := math.Sqrt(ssq / num)

	fmt.Println("--- echo statistics ---")
	fmt.Printf("round-trip min/avg/max/stddev = %0.3f/%0.3f/%0.3f/%0.3f ms\n", min, mean, max, sd)
}
