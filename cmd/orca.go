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
	}

	app.Run(os.Args)
}

func startReflector(c *cli.Context) error {
	fmt.Println("beginning reflector: ", c.Args().First())
	return nil
}

func startGenerator(c *cli.Context) error {
	fmt.Println("beginning generator: ", c.Args().First())
	return nil
}
