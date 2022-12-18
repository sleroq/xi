package util

type Package struct {
	Name    string
	Version string
}

func PkgsToStrings(pkgs []Package) []string {
	res := []string{}

	for _, pkg := range pkgs {
		res = append(res, pkg.Name)
	}

	return res
}
