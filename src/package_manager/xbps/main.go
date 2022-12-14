package xbps

import (
	"fmt"
	"strings"
	"xi/src/package_manager/util"
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
	out, err := util.Run("xbps-query", "-m")
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
