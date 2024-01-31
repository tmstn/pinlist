package main

import (
	"github.com/urfave/cli/v2"
)

func createApp() *cli.App {
	app := &cli.App{
		Name:  "pinlist",
		Usage: "pinboard links extractor",
		Commands: []*cli.Command{
			{
				Name:    "authenticate",
				Aliases: []string{"a", "auth"},
				Usage:   "authenticate with pinboard",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "username",
						Usage:    "username to login to pinboard with",
						Aliases:  []string{"u", "user"},
						Required: true,
					},
					&cli.StringFlag{
						Name:    "passphrase",
						Usage:   "passphrase to encrypt and decrypt credentials with",
						Aliases: []string{"pp", "pass"},
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:   "login",
						Usage:  "login to pinboard",
						Action: login,
					},
					{
						Name:   "logout",
						Usage:  "logout from pinboard",
						Action: logout,
					},
				},
			},
			{
				Name:    "query",
				Aliases: []string{"q", "qry"},
				Usage:   "extract links using a query",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "username",
						Usage:    "username to login to pinboard with",
						Aliases:  []string{"u", "user"},
						Required: true,
					},
					&cli.StringFlag{
						Name:    "passphrase",
						Usage:   "passphrase to encrypt and decrypt credentials with",
						Aliases: []string{"pp", "pass"},
					},
					&cli.StringFlag{
						Name:    "url",
						Usage:   "url to query",
						Aliases: []string{"l", "ul"},
					},
					&cli.StringSliceFlag{
						Name:    "tag",
						Usage:   "tag to query",
						Aliases: []string{"t", "tg"},
					},
					&cli.IntFlag{
						Name:    "page",
						Usage:   "page size per request",
						Aliases: []string{"p", "pg"},
					},
					&cli.IntFlag{
						Name:    "max",
						Usage:   "maximum number of links to extract",
						Aliases: []string{"m", "mx"},
					},
					&cli.BoolFlag{
						Name:    "remove",
						Usage:   "remove links after extraction",
						Aliases: []string{"r", "rm"},
					},
				},
				Action: query,
			},
		},
	}

	return app
}
