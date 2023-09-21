package util

import (
	"fmt"
	"os"
	"os/exec"
)

func GetDockerfilePath(lang string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\\docker_runtimes\\%s", wd, lang), nil
}

func RunDockerCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	// Set the output and error pipes to capture the command's output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	return err
}
