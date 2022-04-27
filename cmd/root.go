package cmd

import (
	"os"

	"github.com/enmand/cf-tf-diff/internal/diff"
	"github.com/enmand/cf-tf-diff/internal/terraform"
	"github.com/jbowes/cling"
	"github.com/urfave/cli/v2"
)

func Execute() error {
	app := &cli.App{
		Name:  "cf-tf-diff",
		Usage: "A CLI tool to compare CloudFlare and Terraform state",

		Commands: []*cli.Command{
			{
				Name:    "compare",
				Aliases: []string{"c"},
				Usage:   "Compare CloudFlare and Terraform state",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Aliases:  []string{"p"},
						Usage:    "The path to the Terraform project",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "cloudflare-email",
						Aliases:  []string{"c"},
						Usage:    "The email address of the CloudFlare account",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "cloudflare-key",
						Aliases:  []string{"k"},
						Usage:    "The API key of the CloudFlare account",
						Required: true,
					},
				},
				Action: compare,
			},
		},
	}

	return app.Run(os.Args)
}

func compare(c *cli.Context) error {
	rs, err := terraform.GetBackend(c.String("path"))
	if err != nil {
		return cling.Wrap(err, "unable to get backend")
	}

	sf, err := rs.GetStateFile()
	if err != nil {
		return cling.Wrap(err, "unable to get state file")
	}

	cfs, err := diff.ParseState(sf)
	if err != nil {
		return cling.Wrap(err, "unable to parse state file")
	}
	_ = cfs

	return nil
}
