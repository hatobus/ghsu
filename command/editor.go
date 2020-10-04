package command

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/xerrors"

	"github.com/hatobus/ghsu/editor"
	"github.com/hatobus/ghsu/env"
	"github.com/hatobus/ghsu/updator"
)

func SetEditor() func(c *cli.Context) error {
	return func(c *cli.Context) error {

		editorName := c.Args().First()
		if editorName == "" {
			return xerrors.Errorf("please input editor name")
		}

		if err := os.Setenv("GHSU_EDITOR", editorName); err != nil {
			return err
		}

		log.Printf("Successfully set ghsu editor %v \n", editorName)

		return nil
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

		if err := gc.ShowUpSetSecrets(); err != nil {
			log.Println(err)
			return err
		}

		return nil
	}
}
