package pacman

import (
	"fmt"
	"strings"
	"xi/src/package_manager/util"
)

type Pacman struct {
	name string
}

func New() *Pacman {
	pm := new(Pacman)
	pm.name = "pacman"

	return pm
}

func (*Pacman) Install(pkgs ...util.Package) error {
	cmd := []string{"pacman", "-S", "--noconfirm"}
	res, err := util.SudoRun(append(cmd, util.PkgsToStrings(pkgs)...)...)
	if err != nil {
		return fmt.Errorf("executing pacman -S: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (*Pacman) Remove(pkgs ...util.Package) error {
	cmd := []string{"pacman", "-R"} //, "--noconfirm"}
	res, err := util.SudoRun(append(cmd, util.PkgsToStrings(pkgs)...)...)
	if err != nil {
		return fmt.Errorf("executing pacman -R: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (*Pacman) GetInstalled() ([]util.Package, error) {
	out, err := util.Run("pacman", "-Qe")
	if err != nil {
		return nil, fmt.Errorf("executing pacman -Qe: %w", err)
	}

	splitted := strings.Split(string(out), "\n")

	pkgs := []util.Package{}

	for _, pkgString := range splitted {
		if strings.Trim(pkgString, " ") != "" {
			pkg := strings.Split(pkgString, " ")
			if len(pkg) == 2 {
				pkgs = append(pkgs, util.Package{Name: pkg[0], Version: pkg[1]})
			}
		}
	}

	return pkgs, nil
}
