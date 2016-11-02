// This command implements the Mora console utility that can run one of two
// background processes: reflectors and generators. These processes should
// be managed with LaunchAgent or Upstart on OS X and Ubunut machines.
package main

import (
	"os"

	"github.com/bbengfort/orca"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

func main() {

	// Load the .env file if it exists
	godotenv.Load()

	// Instantiate the command line application.
	app := cli.NewApp()
	app.Name = "orca"
	app.Usage = "run orca listener or ping"
	app.Version = orca.Version
	app.Author = "Benjamin Bengfort"
	app.Email = "bengfort@cs.umd.edu"
	app.Commands = []cli.Command{
		{
			Name:      "listen",
			Usage:     "listen for pings",
			Action:    listen,
			ArgsUsage: " ",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "s, silent",
					Usage: "do not log messages",
				},
				cli.StringFlag{
					Name:  "a, addr",
					Usage: "specify the address to listen on",
					Value: ":3265",
				},
			},
		},
		{
			Name:      "ping",
			Usage:     "send echo requests to the specified address",
			Action:    ping,
			ArgsUsage: "ipaddr",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "s, silent",
					Usage: "do not log messages",
				},
				cli.IntFlag{
					Name:  "n, num",
					Usage: "limit the number of echo requests",
					Value: 4,
				},
			},
		},
	}

	app.Run(os.Args)
}

func listen(c *cli.Context) error {

	// Get the arguments from the command line
	silent := c.Bool("silent")
	addr := c.String("addr")

	orcaApp, err := orca.Init(silent, addr)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err = orcaApp.Reflect(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func ping(c *cli.Context) error {

	// Validate arguments
	if c.NArg() != 1 {
		return cli.NewExitError("Specify an IP address to ping", 1)
	}

	// Get the arguments from the command line
	silent := c.Bool("silent")
	count := c.Int("num")
	addr := c.Args()[0]

	// Initialize the app
	orcaApp, err := orca.Init(silent, "")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err := orcaApp.Generate(addr, count); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
