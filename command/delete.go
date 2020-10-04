package command

import (
	"context"
	"log"

	"github.com/hatobus/ghsu/updator"
	"github.com/urfave/cli"
)

func DeleteSecrets(gc *updator.GithubClient) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		args := []string{}
		args = append(args, c.Args().First())
		args = append(args, c.Args().Tail()...)

		ctx := context.Background()

		correctlyDeleted := []string{}
		for _, key := range args {
			if gc.ExistRepoSecret(gc.Owner, gc.Repo, key) {
				correctlyDeleted = append(correctlyDeleted, key)
			} else {
				log.Printf("%v is not found\n", key)
			}
		}

		for _, key := range correctlyDeleted {
			_, err := gc.Client.Actions.DeleteSecret(ctx, gc.Owner, gc.Repo, key)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		log.Printf("delete secrets %v\n", correctlyDeleted)

		return nil
	}
}
