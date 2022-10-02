package syncer

import (
	"fmt"
	package_manager "xi/src/package-manager"
)

type Package string

func contains[T comparable](slice []T, value T) bool {
	for _, name := range slice {
		if name == value {
			return true
		}
	}

	return false
}

func SyncPackages(pm package_manager.Pm, packages []Package) error {
	installed, err := pm.GetInstalled()
	if err != nil {
		return fmt.Errorf("getting list of installed packages: %w", err)
	}

	for _, pkg := range packages {
		if contains(installed, string(pkg)) {
			continue
		}

		err := pm.Install(string(pkg))
		if err != nil {
			return fmt.Errorf("installing package: %w", err)
		}
	}

	return nil
}
