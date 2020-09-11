package tools

import (
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

// RunInDir ...
func RunInDir(cmd, dir string) ([]byte, error) {
	if os.Getenv("DEBUG") != "" {
		log.Println(cmd)
	}
	command := exec.Command("sh", "-c", "set -o pipefail && "+cmd)
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil {
		return output, cli.NewExitError(err, 1)
	}
	return output, nil
}

// Run ...
func Run(cmd string) ([]byte, error) {
	return RunInDir(cmd, ".")
}

// RunInteractiveInDir ...
func RunInteractiveInDir(cmd, dir string) error {
	if os.Getenv("DEBUG") != "" {
		log.Println(cmd)
	}
	command := exec.Command("sh", "-c", "set -o pipefail && "+cmd)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Dir = dir
	err := command.Run()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

// RunWithInput ...
func RunWithInput(cmd string, input []byte) ([]byte, error) {
	command := exec.Command("sh", "-c", cmd)
	command.Stdin = bytes.NewReader(input)
	data, err := command.Output()
	if err != nil {
		return data, cli.NewExitError(err.Error(), 1)
	}
	return data, nil
}

// RunInteractive ...
func RunInteractive(cmd string) error {
	return RunInteractiveInDir(cmd, "")
}
