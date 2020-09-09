package env

import (
	"bufio"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"strings"
)

func FileExist(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}

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
		scanned := scanner.Text()
		if len(scanned) == 0 {
			continue
		}

		elems := strings.Split(scanned, "=")
		if len(elems) < 2 {
			return nil, xerrors.New("read data failed invalid data format")
		}

		key := elems[0]
		value := strings.Join(elems[1:len(elems)], "=")
		values[key] = value
	}

	if len(values) == 0 {
		return nil, xerrors.New("read data failed invalid data format")
	}

	return values, nil
}
