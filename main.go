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
	cmdDelete = "delete"
)

const (
	subCmdEnv     = "env"
	subCmdFile    = "file"
	subCmdEditor  = "editor"
	subCmdSecrets = "secrets"
)

const (
	usageShow   = "show up the value you want to upload"
	usageSet    = "set your environment variable"
	usageEditor = "set your ghsu editor"
	usageDelete = "delete variable"
)

const (
	subUsageShowEnv       = "show from .env file format"
	subUsageShowFile      = "show up from user's file"
	subUsageShowFileText  = "\"key\" is the key name of github secrets, \"filename\" is the specific file name want to upload (Base64 encrypted)."
	subUsageShowSecrets   = "show up set github secrets"
	subUsageSetEnv        = "set environment variable from .env format file"
	subUsageSetEditor     = "environment variable editor mode"
	subUsageSetEditorText = "\".env\" file format"
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
					Name:   subCmdEnv,
					Usage:  subUsageShowEnv,
					Action: command.ShowFromEnvFile(),
				},
				{
					Name:      subCmdFile,
					Usage:     subUsageShowFile,
					UsageText: subUsageShowFileText,
					Action:    command.ShowFromUserFile(),
				},
				{
					Name:   subCmdSecrets,
					Usage:  subUsageShowSecrets,
					Action: command.ShowSecret(githubClient),
				},
			},
		},
		{
			Name:  cmdSet,
			Usage: usageSet,
			Subcommands: []cli.Command{
				{
					Name:   subCmdEnv,
					Usage:  subUsageSetEnv,
					Action: command.SetEnvironmentFromEnv(githubClient),
				},
				{
					Name:      subCmdEditor,
					Usage:     subUsageSetEditor,
					UsageText: subUsageSetEditorText,
					Action:    command.SetEditor(),
				},
			},
		},
		{
			Name:   cmdEditor,
			Usage:  usageEditor,
			Action: command.UploadFromEditor(githubClient),
		},
		{
			Name:   cmdDelete,
			Usage:  usageDelete,
			Action: command.DeleteSecrets(githubClient),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
