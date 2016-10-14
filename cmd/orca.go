// This command implements the Mora console utility that can run one of two
// background processes: reflectors and generators. These processes should
// be managed with LaunchAgent or Upstart on OS X and Ubunut machines.
package main

import (
	"fmt"
	"os"

	"github.com/bbengfort/orca"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var orcaApp *orca.App

func main() {

	// Load the .env file if it exists
	godotenv.Load()

	// Instantiate the command line application.
	app := cli.NewApp()
	app.Name = "orca"
	app.Usage = "run orca listener or generator in the background"
	app.Version = orca.Version
	app.Author = "Benjamin Bengfort"
	app.Email = "bengfort@cs.umd.edu"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c, config",
			Usage: "specify the path to a yaml configuration",
		},
	}
	app.Before = initOrca
	app.Commands = []cli.Command{
		{
			Name:   "reflect",
			Usage:  "run the reflector daemon",
			Action: startReflector,
		},
		{
			Name:   "generate",
			Usage:  "run the generator daemon",
			Action: startGenerator,
		},
		{
			Name:   "config",
			Usage:  "print the configuration and exit",
			Action: printConfig,
		},
	}

	app.Run(os.Args)
}

func initOrca(c *cli.Context) error {
	var err error

	// Initialize the application
	if orcaApp, err = orca.Init(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	// Modify the config from the command line if necessary
	if c.String("config") != "" {
		path := c.String("config")
		if err = orcaApp.Config.Read(path); err != nil {
			msg := fmt.Sprintf("Unable to read configuration at %s", path)
			return cli.NewExitError(msg, 2)
		}
	}

	return nil
}

func startReflector(c *cli.Context) error {

	if err := orcaApp.Reflect(); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	return nil
}

func startGenerator(c *cli.Context) error {
	addr, err := orca.ResolveAddr("")
	device := &orca.Device{
		Name:   "apollo",
		IPAddr: addr,
	}

	reply, err := orca.Ping(device)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println(reply)
	return nil

	// if err := orcaApp.Generate(); err != nil {
	// 	return cli.NewExitError(err.Error(), 2)
	// }
	//
	// return nil
}

func printConfig(c *cli.Context) error {
	// Print the configuration and exit
	fmt.Println(orcaApp.Config.String())
	return nil
}
