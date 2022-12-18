package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"xi/src/package_manager"
	"xi/src/package_manager/util"

	"github.com/BurntSushi/toml"
)

const defaultIndent = "  "
const defaultGroup = "InstalledPackages"

func configDir() string {
	dir := os.Getenv("XDG_CONFIG_HOME")

	if dir == "" {
		dir = path.Join(os.Getenv("HOME"), ".config")
	}

	return path.Join(dir, "xi")
}

type PackageGroup []util.Package
type Groups map[string]PackageGroup
type Managers struct {
	Xbps   Groups
	Pacman Groups
	Paru   Groups
	Yay    Groups
}

type Configuration struct {
	Managers Managers
}

type rawPackageGroup []interface{}
type rawGroups map[string]rawPackageGroup
type rawManagers map[string]rawGroups

type rawConfiguration struct {
	Packages rawManagers
}

// generateConfig will create config.toml with installed packages
func generateConfig(filePath string) (*Configuration, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("creating config file: %w", err)
	}
	defer f.Close()

	config := new(rawConfiguration)
	config.Packages = make(rawManagers)

	pms, err := package_manager.AvailibleManagers()
	if err != nil {
		return nil, fmt.Errorf("getting availible package managers: %s", err)
	}

	for name, manager := range pms {
		// prefering paru over yay and pacman
		if name == "yay" || name == "pacman" {
			if _, ok := pms["paru"]; ok {
				continue
			}
		}

		pkgs, err := manager.GetInstalled()
		if err != nil {
			return nil, fmt.Errorf("getting installed packages from %s: %s", name, err)
		}

		if len(pkgs) == 0 {
			continue
		}

		group := rawPackageGroup{}

		for _, pkg := range pkgs {
			group = append(group, pkg.Name)
		}

		config.Packages[name] = rawGroups{defaultGroup: group}
	}

	encoder := toml.NewEncoder(f)
	encoder.Indent = defaultIndent

	err = encoder.Encode(config)
	if err != nil {
		return nil, fmt.Errorf("encoding new config: %w", err)
	}

	err = finalizeConfig(filePath)
	if err != nil {
		return nil, fmt.Errorf("editing new config: %w", err)
	}

	return resolveConfig(*config)
}

// finalizeConfig will reformat default configuration and add examples
func finalizeConfig(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("opening new config: %w", err)
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading new config: %w", err)
	}

	config := string(contents)
	indentPkgs := strings.Repeat(defaultIndent, 3)
	indentDouble := strings.Repeat(defaultIndent, 2)

	config = strings.ReplaceAll(config, "[\"", fmt.Sprint("[\n", indentPkgs, "\""))
	config = strings.ReplaceAll(config, ", ", fmt.Sprint(",\n", indentPkgs))
	config = strings.ReplaceAll(config, "\"]", fmt.Sprint("\",\n", indentDouble, "]"))

	toBeReplaced := fmt.Sprintf("%s = [", defaultGroup)
	resWithGroup := fmt.Sprintf("MyGroup = [ \"neofetch\" ]\n%s%s", indentDouble, toBeReplaced)
	config = strings.ReplaceAll(config, toBeReplaced, resWithGroup)

	err = os.WriteFile(filePath, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("saving new config: %w", err)
	}

	return nil
}

func resolveGroups(inputGroups rawGroups) (groups Groups, err error) {
	groups = Groups{}

	for groupName, pkgs := range inputGroups {
		group := PackageGroup{}

		for _, rawPkg := range pkgs {
			if val, ok := rawPkg.(string); ok {
				group = append(group, util.Package{val, ""})
			} else if val, ok := rawPkg.([]interface{}); ok {
				pkg := util.Package{}
				for i, value := range val {
					if s, ok := value.(string); ok {
						switch i {
						case 0:
							pkg.Name = s
						case 1:
							pkg.Version = s
						}
					} else {
						return Groups{}, fmt.Errorf("Can't parse package \"%s\": \"%s\" is not a string", rawPkg, s)
					}
				}

				if len(val) > 0 {
					group = append(group, pkg)
				}
			} else {
				return Groups{}, fmt.Errorf("Can't parse package \"%s\"", rawPkg)
			}
		}

		if len(group) > 0 {
			groups[groupName] = group
		}
	}

	return groups, nil
}

func resolveManagers(inputManagers rawManagers) (managers Managers, err error) {
	for manager, groups := range inputManagers {
		switch manager {
		case "xbps":
			managers.Xbps, err = resolveGroups(groups)
		case "pacman":
			managers.Pacman, err = resolveGroups(groups)
		case "paru":
			managers.Paru, err = resolveGroups(groups)
		case "yay":
			managers.Yay, err = resolveGroups(groups)
		default:
			fmt.Printf("WARN: \"%s\" package manager is not supported\n", manager)
			continue
		}

		if err != nil {
			return Managers{}, fmt.Errorf("resolving groups for \"%s\": %s", manager, err)
		}
	}

	return managers, nil
}

// resolveConfig resolves known types in parsed configuration file
func resolveConfig(parsed rawConfiguration) (config *Configuration, err error) {
	config = &Configuration{}
	config.Managers, err = resolveManagers(parsed.Packages)
	if err != nil {
		return nil, fmt.Errorf("resolving package managers: %s", err)
	}
	return config, nil
}

// readConfig will find and parse config file
// from $XDG_CONFIG_HOME/xi/config.toml or $HOME/.config/xi/config.toml
// If config does not exist it will be generated
func readConfig() (config *Configuration, err error) {
	dir := configDir()
	file := path.Join(dir, "config.toml")

	if _, err = os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("creating config directory: %w", err)
		}
	}

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		config, err = generateConfig(file)
		if err != nil {
			return nil, fmt.Errorf("generating new config: %w", err)
		}

		return config, nil
	}

	var rawConfig rawConfiguration

	_, err = toml.DecodeFile(file, &rawConfig)
	if err != nil {
		return nil, fmt.Errorf("config decoding: %w", err)
	}

	config, err = resolveConfig(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("config resolving: %w", err)
	}

	return config, nil
}
