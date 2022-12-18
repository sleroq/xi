package paru

import (
	"fmt"
	"os/exec"
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
	res, err := util.SudoRun(append(cmd, util.PkgsToString(pkgs))...)
	if err != nil {
		return fmt.Errorf("executing paru -S: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (*Paru) Remove(pkgs ...util.Package) error {
	cmd := []string{"paru", "-R", "--noconfirm"}
	res, err := util.SudoRun(append(cmd, util.PkgsToString(pkgs))...)
	if err != nil {
		return fmt.Errorf("executing paru -R: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (*Paru) GetInstalled() ([]util.Package, error) {
	cmd := exec.Command("paru", "-Qe")

	out, err := cmd.Output()
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
