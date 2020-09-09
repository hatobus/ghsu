package env

import (
	"bufio"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"strings"
)

func ReadFromDotEnvFile(fname string) (map[string]string, error) {
	if fname == "" {
		fname = "./.env"
	}

	absPath, err := filepath.Abs(fname)
	if err!= nil {
		return nil, err
	}

	fp, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	values := map[string]string{}

	for scanner.Scan() {
		elems := strings.Split(scanner.Text(), "=")
		if len(elems) < 2 {
			return nil, xerrors.Errorf("read data failed invalid data format")
		}

		key := elems[0]
		var value string
		if len(elems) == 2 {
			value = elems[1]
		} else {
			value = strings.Join(elems[1:len(elems)], "=")
		}
		values[key] = value
	}


	return values, nil
}
