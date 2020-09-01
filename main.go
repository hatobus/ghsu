package main

import (
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

	log.Println(githubClient)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
