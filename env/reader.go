package env

import (
	"bufio"
	"encoding/base64"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func fileExist(fname string) bool {
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
