package syncer

import (
	"fmt"
	"xi/src/package_manager"
)

type Package string
type PackageGroup []Package
type Groups map[string]PackageGroup

func contains[T comparable](slice []T, value T) bool {
	for _, name := range slice {
		if name == value {
			return true
		}
	}

	return false
}

func SyncGroups(pm package_manager.Pm, groups Groups) error {
	configPkgs := []Package{}
	for _, pkgs := range groups {
		configPkgs = append(configPkgs, pkgs...)
	}

	err := InstallMissings(pm, configPkgs)
	if err != nil {
		return fmt.Errorf("installing packages: %w", err)
	}

	err = RemoveRedundant(pm, configPkgs)
	if err != nil {
		return fmt.Errorf("removing packages: %w", err)
	}

	return nil
}

func InstallMissings(pm package_manager.Pm, pkgs []Package) error {
	installed, err := pm.GetInstalled()
	if err != nil {
		return fmt.Errorf("getting list of installed packages: %w", err)
	}

	for _, pkg := range pkgs {
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

func RemoveRedundant(pm package_manager.Pm, configPkgs []Package) error {
	installed, err := pm.GetInstalled()
	if err != nil {
		return fmt.Errorf("getting list of installed packages: %w", err)
	}

	for _, pkg := range installed {
		if contains(configPkgs, Package(pkg)) {
			continue
		}

		err := pm.Remove(pkg)
		if err != nil {
			return fmt.Errorf("installing package: %w", err)
		}
	}

	return nil
}
