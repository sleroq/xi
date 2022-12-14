package yay

import (
	"fmt"
	"strings"
	"xi/src/package_manager/util"
)

type Yay struct {
	name string
}

func New() *Yay {
	pm := new(Yay)
	pm.name = "yay"

	return pm
}

func (*Yay) Install(pkgs ...util.Package) error {
	cmd := []string{"yay", "-S", "--noconfirm"}
	res, err := util.RunIn(append(cmd, util.PkgsToStrings(pkgs)...)...)
	if err != nil {
		return fmt.Errorf("executing yay -S: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (*Yay) Remove(pkgs ...util.Package) error {
	cmd := []string{"yay", "-R"} //, "--noconfirm"}
	res, err := util.RunIn(append(cmd, util.PkgsToStrings(pkgs)...)...)
	if err != nil {
		return fmt.Errorf("executing yay -R: %w", err)
	}

	fmt.Println("Result:", string(res))

	return nil
}

func (*Yay) GetInstalled() ([]util.Package, error) {
	out, err := util.Run("yay", "-Qe")
	if err != nil {
		return nil, fmt.Errorf("executing yay -Qe: %w", err)
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
