package editor

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func startUpEditor(editorCmd string, fname string) error {
	// stdin, stdout, stdout を cmd と紐づけるとエディタを開くことができる
	cmd := exec.Command(editorCmd, fname)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RunEditor() (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer func() {
		f.Close()
	}()

	editorCmd := os.Getenv("GHSU_EDITOR")
	if editorCmd == "" {
		editorCmd = "vi"
	}

	if err = startUpEditor(editorCmd, f.Name()); err != nil {
		return "", err
	}

	return f.Name(), nil
}
