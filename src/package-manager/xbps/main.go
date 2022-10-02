package xbps

import (
	"fmt"
	"os/exec"
	"strings"
	"xi/src/package-manager/util"
)

type Xbps struct {
	name string
}

func New() *Xbps {
	pm := new(Xbps)
	pm.name = "xbps"

	return pm
}

func (pm *Xbps) Install(pkgs ...string) error {
	cmd := []string{"xbps-install", "-y"}
	res, err := util.SudoRun(append(cmd, pkgs...)...)
	if err != nil {
		return fmt.Errorf("executing xpbs-install: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (pm *Xbps) Remove(pkgs ...string) error {
	cmd := []string{"xbps-remove", "-y"}
	res, err := util.SudoRun(append(cmd, pkgs...)...)
	if err != nil {
		return fmt.Errorf("executing xpbs-install: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (*Xbps) GetInstalled() ([]string, error) {
	cmd := exec.Command("xbps-query", "-m")

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("executing xbps-query -m: %w", err)
	}

	splitted := strings.Split(string(out), "\n")


	pkgs := []string{}

	for _, name := range splitted {
		if strings.Trim(name, " ") != "" {
			pkgs = append(pkgs, name)
		}
	}

	return pkgs, nil
}
