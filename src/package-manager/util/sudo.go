package util

import (
	"fmt"
	"os"
	"os/exec"
)

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
