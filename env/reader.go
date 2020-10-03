package env

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

func fileExist(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}

func tmpFileDelete(fname string) error {
	if strings.HasPrefix(fname, "/tmp") {
		return os.Remove(fname)
	}
	return nil
}

func ReadFromDotEnvFile(fname string) (map[string]string, error) {
	if fname == "" {
		fname = "./.env"
	}

	absPath, err := filepath.Abs(fname)
	if err != nil {
		return nil, err
	}

	if !fileExist(absPath) {
		return nil, xerrors.New("no such file, check your input file name")
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

		elems := strings.SplitN(scanned, "=", 2)
		if len(elems) < 2 {
			return nil, xerrors.New("read data failed invalid data format")
		}

		values[elems[0]] = elems[1]
	}

	if len(values) == 0 {
		return nil, xerrors.New("read data failed invalid data format")
	}

	err = tmpFileDelete(fname)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func ReadFromFileEncryptBase64(fname string) (string, error) {
	if !fileExist(fname) {
		return "", xerrors.New("no such file, check your input file name")
	}

	absPath, err := filepath.Abs(fname)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}
