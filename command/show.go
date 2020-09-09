package command

import (
	"fmt"
	"golang.org/x/xerrors"

	"github.com/hatobus/ghsu/env"
	"github.com/urfave/cli"
)

func ShowFromEnvFile() func(c *cli.Context) error {
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

		fmt.Printf("new variables from \"%v\"\n", fname)
		fmt.Printf("\x1b[33mkey: \x1b[0m\t \x1b[34mvalue\x1b[0m\n")

		for key, val := range values {
			fmt.Printf("\x1b[33m%v\x1b[0m: \t \x1b[34m%v \x1b[0m\n", key, val)
		}

		return nil
	}
}

func ShowFromUserFile() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		key := c.Args().Get(0)
		if key == "" {
			return xerrors.New("please input key name")
		}

		filename := c.Args().Get(1)

		encrypted, err := env.ReadFromFileEncryptBase64(filename)
		if err != nil {
			return err
		}

		fmt.Printf("new variables from \"%v\"\n", filename)
		fmt.Printf("\x1b[33mkey:\x1b[0m\t%v\n", key)
		fmt.Printf("\x1b[34mvalue (Base64 encrypted):\x1b[0m\n%v\n", encrypted)

		return nil
	}
}
