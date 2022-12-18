package package_manager

import (
	"fmt"
	"os/exec"
	"xi/src/package_manager/pacman"
	"xi/src/package_manager/paru"
	"xi/src/package_manager/util"
	"xi/src/package_manager/yay"
)

type PM interface {
	Install(pkgs ...util.Package) error
	Remove(pkgs ...util.Package) error
	GetInstalled() ([]util.Package, error)
}

func New(pmName string) (PM, error) {
	switch {
	// case pmName == "xbps":
	// 	return xbps.New(), nil
	case pmName == "pacman":
		return pacman.New(), nil
	case pmName == "paru":
		return paru.New(), nil
	case pmName == "yay":
		return yay.New(), nil
	default:
		return nil, fmt.Errorf("don't know such package manager")
	}
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func addIfAvailible(pms map[string]PM, s string) (map[string]PM, error) {
	if commandExists(s) {
		pm, err := New(s)
		if err != nil {
			return nil, fmt.Errorf("creating \"%s\" package manager interface: %w", s, err)
		}
		pms[s] = pm
	}

	return pms, nil
}

func AvailibleManagers() (map[string]PM, error) {
	pms := make(map[string]PM)

	pms, err := addIfAvailible(pms, "xbps")
	if err != nil {
		return nil, err
	}
	pms, err = addIfAvailible(pms, "pacman")
	if err != nil {
		return nil, err
	}
	pms, err = addIfAvailible(pms, "paru")
	if err != nil {
		return nil, err
	}
	pms, err = addIfAvailible(pms, "yay")
	if err != nil {
		return nil, err
	}

	return pms, nil
}
