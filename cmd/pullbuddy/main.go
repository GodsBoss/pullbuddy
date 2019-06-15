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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "server-addr",
					},
					&cli.StringFlag{
						Name: "image",
					},
				},
				Action: func(ctx *cli.Context) error {
					client := &pullbuddy.Client{
						Addr: ctx.String("server-addr"),
						Out:  os.Stdout,
					}
					return client.Schedule(ctx.String("image"))
				},
			},
			{
				Name: "status",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "server-addr",
					},
				},
				Action: func(ctx *cli.Context) error {
					client := &pullbuddy.Client{
						Addr: ctx.String("server-addr"),
						Out:  os.Stdout,
					}
					return client.Status()
				},
			},
		},
	}
	return app
}
