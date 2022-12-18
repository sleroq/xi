package main

import (
	"fmt"
	"strings"
	pm "xi/src/package_manager"
	"xi/src/package_manager/util"
)

func contains[T comparable](slice []T, value T) bool {
	for _, name := range slice {
		if name == value {
			return true
		}
	}

	return false
}

func syncPackages(config Configuration) error {
	if config.Managers.Xbps != nil {
		xbps, err := pm.New("xbps")
		if err != nil {
			return fmt.Errorf("getting xbps: %s", err)
		}
		err = syncGroups(xbps, config.Managers.Xbps)
		if err != nil {
			return fmt.Errorf("syncing xbps groups: %s", err)
		}
	}
	if config.Managers.Pacman != nil {
		pacman, err := pm.New("pacman")
		if err != nil {
			return fmt.Errorf("getting pacman: %s", err)
		}
		err = syncGroups(pacman, config.Managers.Pacman)
		if err != nil {
			return fmt.Errorf("syncing pacman groups: %s", err)
		}
	}
	if config.Managers.Paru != nil {
		paru, err := pm.New("paru")
		if err != nil {
			return fmt.Errorf("getting paru: %s", err)
		}
		err = syncGroups(paru, config.Managers.Paru)
		if err != nil {
			return fmt.Errorf("syncing paru groups: %s", err)
		}
	}
	if config.Managers.Yay != nil {
		yay, err := pm.New("pacman")
		if err != nil {
			return fmt.Errorf("getting yay: %s", err)
		}
		err = syncGroups(yay, config.Managers.Yay)
		if err != nil {
			return fmt.Errorf("syncing yay groups: %s", err)
		}
	}

	return nil
}

func syncGroups(pm pm.PM, pkgsGroups Groups) error {
	statedPackages := []util.Package{}

	// Installing per group:
	for _, group := range pkgsGroups {
		err := InstallMissings(pm, group)
		if err != nil {
			return fmt.Errorf("installing packages: %w", err)
		}

		statedPackages = append(statedPackages, group...)
	}

	err := RemoveRedundant(pm, statedPackages)
	if err != nil {
		return fmt.Errorf("removing packages: %w", err)
	}

	return nil
}

func InstallMissings(pm pm.PM, pkgs []util.Package) error {
	installed, err := pm.GetInstalled()
	if err != nil {
		return fmt.Errorf("getting list of installed packages: %w", err)
	}
	installedList := util.PkgsToStrings(installed)

	toInstall := []util.Package{}
	for _, pkg := range pkgs {
		if !contains(installedList, pkg.Name) {
			toInstall = append(toInstall, pkg)
		}

	}

	if len(toInstall) > 0 {
		err = pm.Install(toInstall...)
		if err != nil {
			return fmt.Errorf("installing package: %w", err)
		}
	}

	return nil
}

func RemoveRedundant(pm pm.PM, pkgs []util.Package) error {
	installed, err := pm.GetInstalled()
	if err != nil {
		return fmt.Errorf("getting list of installed packages: %w", err)
	}
	stated := util.PkgsToStrings(pkgs)

	toRemove := []util.Package{}
	for _, pkg := range installed {
		if !contains(stated, pkg.Name) {
			toRemove = append(toRemove, pkg)

		}

	}

	toRemovString := strings.Join(util.PkgsToStrings(toRemove), ", ")
	fmt.Println("Packages to be removed:", toRemovString)

	if len(toRemove) > 0 {
		err = pm.Remove(toRemove...)
		if err != nil {
			return fmt.Errorf("removing package: %w", err)
		}
	}

	return nil
}
