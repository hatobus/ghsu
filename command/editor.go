package command

import (
	"os"

	"github.com/urfave/cli"
	"golang.org/x/xerrors"
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
