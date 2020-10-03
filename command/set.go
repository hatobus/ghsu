package command

import (
	"context"
	"fmt"

	"github.com/urfave/cli"

	"github.com/hatobus/ghsu/env"
	"github.com/hatobus/ghsu/updator"
)

func SetEnvironmentFromEnv(gc *updator.GithubClient) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		var fname string
		if c.Args().First() == "" {
			fname = "./.env"
		} else {
			fname = c.Args().First()
		}

		values, err := env.ReadFromDotEnvFile(fname)
		if err != nil {
			return err
		}

		encrypted, err := gc.GenerateEncryptedSecret(values)
		if err != nil {
			return err
		}

		for _, secret := range encrypted {
			ctx := context.Background()
			_, err := gc.Client.Actions.CreateOrUpdateSecret(ctx, gc.Owner, gc.Repo, secret)
			if err != nil {
				fmt.Printf("Key: %v, registration error", secret.Name)
			}
		}

		return nil
	}
}
