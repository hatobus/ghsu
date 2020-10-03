package command

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/xerrors"

	"github.com/hatobus/ghsu/editor"
	"github.com/hatobus/ghsu/env"
	"github.com/hatobus/ghsu/updator"
)

func SetEditor() func(c *cli.Context) error {
	return func(c *cli.Context) error {

		editor := c.Args().First()
		if editor == "" {
			return xerrors.Errorf("please input editor name")
		}

		return os.Setenv("GHSU_EDITOR", editor)
	}
}

func UploadFromEditor(gc *updator.GithubClient) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		fname, err := editor.RunEditor()
		if err != nil {
			return err
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
