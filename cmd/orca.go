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
		{
			Name:   "createdb",
			Usage:  "create the sqlite3 database",
			Action: createDatabase,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "p, path",
					Usage: "specify a path to create the database",
				},
			},
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
			return cli.NewExitError(msg, 1)
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

	if err := orcaApp.Generate(); err != nil {
		return cli.NewExitError(err.Error(), 3)
	}

	return nil
}

func printConfig(c *cli.Context) error {
	// Print the configuration and exit
	fmt.Println(orcaApp.Config.String())
	return nil
}

func createDatabase(c *cli.Context) error {
	var path string

	// Modify the config from the command line if necessary
	if c.String("path") != "" {

		app := &orca.App{}
		app.Config = &orca.Config{}
		path = c.String("path")
		app.Config.DBPath = path

		// Force the reconnection to the new path
		if err := app.ConnectDB(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if err := app.CreateDB(); err != nil {
			return cli.NewExitError(err.Error(), 4)
		}
	} else {

		path = orcaApp.Config.DBPath
		if err := orcaApp.CreateDB(); err != nil {
			return cli.NewExitError(err.Error(), 4)
		}
	}

	fmt.Printf("Created Orca DB at %s\n", path)
	return nil
}
