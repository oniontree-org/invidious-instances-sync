package main

import (
	"github.com/urfave/cli/v2"
	"time"
)

func (a *Application) commands() {
	a.app = &cli.App{
		Name:    "invidious-instances",
		Version: Version,
		Usage:   "Interact with Invidious's instances API",
		Commands: cli.Commands{
			&cli.Command{
				Name:      "sync",
				Usage:     "Download Invidious instances and update an OnionTree repository",
				ArgsUsage: "<id>",
				Before:    a.handleOnionTreeOpen(),
				Action:    a.handleSyncCommand(),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "url",
						Usage: "Invidious API URL",
						Value: "https://instances.invidio.us/instances.json",
					},
					&cli.BoolFlag{
						Name:  "replace",
						Usage: "replace old URLs",
					},
					&cli.DurationFlag{
						Name:  "timeout",
						Usage: "request timeout",
						Value: 15 * time.Second,
					},
				},
			},
		},
		HideHelpCommand: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "C",
				Value: ".",
				Usage: "change directory to",
			},
		},
	}
}
