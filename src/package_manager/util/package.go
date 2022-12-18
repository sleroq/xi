package util

import (
	"fmt"
	"strings"
)

type Package struct {
	Name    string
	Version string
}

func PkgsToString(pkgs []Package) string {
	s := ""

	for _, pkg := range pkgs {
		s += fmt.Sprint(pkg.Name, " ")
	}

	return strings.Trim(s, " ")
}
