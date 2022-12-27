package util

import (
	"fmt"
	"os"
	"os/exec"
)

func Run(args ...string) ([]byte, error) {
	cmd := exec.Command(args[0], args[1:]...)

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("executing %s: %w", cmd, err)
	}

	return out, nil
}

func RunIn(args ...string) ([]byte, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("executing %s: %w", cmd, err)
	}

	return out, nil
}

func RunOut(args ...string) (error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("executing %s: %w", cmd, err)
	}

	return nil
}

func SudoRun(args ...string) ([]byte, error) {
	cmd := exec.Command("sudo", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("executing %s: %w", cmd, err)
	}

	return out, nil
}
