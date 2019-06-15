package main

import (
	"fmt"
	"os"

	"github.com/GodsBoss/pullbuddy"

	"gopkg.in/urfave/cli.v2"
)

func main() {
	if err := createApp().Run(os.Args); err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
}

func createApp() *cli.App {
	app := &cli.App{
		Name:    "pullbuddy",
		Version: "1.0.0",
		Commands: []*cli.Command{
			{
				Name: "server",
				Subcommands: []*cli.Command{
					{
						Name: "start",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "addr",
							},
						},
						Action: func(ctx *cli.Context) error {
							server := &pullbuddy.Server{
								Addr: ctx.String("addr"),
							}
							return server.Start()
						},
					},
				},
			},
			{
				Name: "schedule",
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
			{
				Name: "status",
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
		},
	}
	return app
}