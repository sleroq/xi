package package_manager

import (
	"fmt"
	"xi/src/package-manager/xbps"
)

type Pm interface {
	Install(pkgs ...string) error
	Remove(pkgs ...string) error
	GetInstalled() ([]string, error)
}

func New() (Pm, error) {
	pmName := "xbps"

	switch {
	case pmName == "xbps":
		return xbps.New(), nil
	default:
		return nil, fmt.Errorf("don't know such package manager")
	}
}
