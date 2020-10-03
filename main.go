package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/hatobus/ghsu/command"
	"github.com/hatobus/ghsu/updator"
)

const (
	cmdShow   = "show"
	cmdSet    = "set"
	cmdEditor = "editor"
)

const (
	usageShow   = "show up the value you want to upload"
	usageSet    = "set your environment variable"
	usageEditor = "set your ghsu editor"
)

func main() {
	app := cli.NewApp()
	app.Name = "ghsu"
	app.Usage = "update github secret"

	githubClient, err := updator.NewGithubClient()
	if err != nil {
		switch err {
		case updator.ErrClientGithubKeyNotFound:
			log.Fatal("Github Token nof found. please set your github token")
		case updator.ErrClientGithubKeyNotFound:
			log.Fatal("failed to get repository data, please check your current directory")
		case updator.ErrFailedToGetOwnerOrRepo:
			log.Fatal("failed to get your public key")
		default:
			log.Fatal(err)
		}
	}

	app.Commands = []cli.Command{
		{
			Name:  cmdShow,
			Usage: usageShow,
			Subcommands: []cli.Command{
				{
					Name:   "env",
					Usage:  "show from .env file format",
					Action: command.ShowFromEnvFile(),
				},
				{
					Name:      "file",
					Usage:     "show up from user's file",
					UsageText: "\"key\" is the key name of github secrets, \"filename\" is the specific file name want to upload (Base64 encrypted).",
					Action:    command.ShowFromUserFile(),
				},
			},
		},
		{
			Name:  cmdSet,
			Usage: usageSet,
			Subcommands: []cli.Command{
				{
					Name:   "env",
					Usage:  "set environment variable from .env format file",
					Action: command.SetEnvironmentFromEnv(githubClient),
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
