package main

import (
	"log"
	"os"
	"sort"

	"github.com/albertogviana/port-service/internal/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{}
	app.UseShortOptionHandling = true
	app.Description = "The port-service is a service that manages the port information"
	app.Commands = []*cli.Command{
		{
			Name:  "import",
			Usage: "The import command imports the port data from a file to a database",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "file",
					Usage: "The file to import from (required)",
				},
			},
			Action: commands.Import,
		},
		{
			Name:   "start",
			Usage:  "Starts the Rest API.",
			Action: commands.Start,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
