package paru

import (
	"fmt"
	"strings"
	"xi/src/package_manager/util"
)

type Paru struct {
	name string
}

func New() *Paru {
	pm := new(Paru)
	pm.name = "paru"

	return pm
}

func (*Paru) Install(pkgs ...util.Package) error {
	cmd := []string{"paru", "-S", "--noconfirm"}
	res, err := util.RunIn(append(cmd, util.PkgsToStrings(pkgs)...)...)
	if err != nil {
		return fmt.Errorf("executing paru -S: %w", err)
	}

	fmt.Println("Result:", string(res))
	return nil
}

func (*Paru) Remove(pkgs ...util.Package) error {
	cmd := []string{"paru", "-R", "--noconfirm"}
	res, err := util.RunIn(append(cmd, util.PkgsToStrings(pkgs)...)...)
	if err != nil {
		return fmt.Errorf("executing paru -R: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (*Paru) GetInstalled() ([]util.Package, error) {
	out, err := util.Run("paru", "-Qe")
	if err != nil {
		return nil, fmt.Errorf("executing paru -Qe: %w", err)
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
