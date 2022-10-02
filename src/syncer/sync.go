package syncer

import "fmt"

type Package struct {
	Name string
	PreInstall string
	PostIntall string
	Test string
}

func Sync (packages []Package) error {
	for _, pkg := range packages {
		fmt.Println("package:", pkg.Name)



	}

	return nil
}
