package main

import (
	"fmt"
	"github.com/hatobus/ghsu/env"
	"github.com/hatobus/ghsu/updator"
	"github.com/urfave/cli"
	"log"
	"os"
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

	fmt.Println(githubClient)

	app.Commands = []cli.Command{
		{
			Name: "show",
			Usage: "show up the value you want to upload",
			Subcommands: []cli.Command{
				{
					Name: "env",
					Usage: "show from .env file format",
					Action: func(c *cli.Context) error {
						fname := c.Args().First()

						values, err := env.ReadFromDotEnvFile(fname)
						if err != nil {
							return err
						}

						fmt.Printf("new variables from \"%v\"\n", fname)

						for key, val := range values {
							fmt.Printf("key: %v \t value: %v \n", key, val)
						}

						return nil
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
