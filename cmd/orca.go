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
			Name:   "listen",
			Usage:  "listen for pings",
			Action: listen,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a, addr",
					Usage: "specify the address to listen on",
				},
			},
		},
		{
			Name:   "ping",
			Usage:  "send echo requests to the specified address",
			Action: ping,
		},
	}

	app.Run(os.Args)
}

func listen(c *cli.Context) error {

	orcaApp, err := orca.Init(false)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err = orcaApp.Reflect(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func ping(c *cli.Context) error {

	orcaApp, err := orca.Init(false)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if c.NArg() != 1 {
		return cli.NewExitError("Specify an IP address to ping", 1)
	}

	addr := c.Args()[0]

	if err := orcaApp.Generate(addr); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
