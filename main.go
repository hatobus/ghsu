package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/hatobus/ghsu/command"
	"github.com/hatobus/ghsu/updator"
)

func main() {
	app := cli.NewApp()
	app.Name = "ghsu"
	app.Usage = "update github secret"

	githubClient, err := updator.NewGithubClient()
	if err != nil {
		if err == updator.ErrClientGithubKeyNotFound {
			log.Fatal("Github Token nof found. please set your github token")
			return
		} else if err == updator.ErrFailedToGetOwnerOrRepo {
			log.Fatal("failed to get repository data, please check your current directory")
			return
		} else if err == updator.ErrFailedToGetPublicKey {
			log.Fatal("failed to get your public key")
			return
		} else {
			log.Fatal(err)
			return
		}
	}

	app.Commands = []cli.Command{
		{
			Name:  "show",
			Usage: "show up the value you want to upload",
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
			Name:  "set",
			Usage: "set your environment variable",
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
