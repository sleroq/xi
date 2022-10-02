package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"xi/src/package_manager"
	"xi/src/syncer"

	"github.com/BurntSushi/toml"
)

func getConfigDir() string {
	configPath := os.Getenv("XDG_CONFIG_HOME")

	if configPath == "" {
		configPath = path.Join(os.Getenv("HOME"), ".config")
	}

	return path.Join(configPath, "xi")
}

type Configuration struct {
	Packages syncer.Groups
}

func generateConfig(filePath string, pm package_manager.Pm) (*Configuration, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("creating config file: %w", err)
	}
	defer f.Close()

	config := new(Configuration)
	config.Packages = make(syncer.Groups)

	packages, err := pm.GetInstalled()
	if err != nil {
		return nil, fmt.Errorf("getting list of installed packages: %w", err)
	}

	group := syncer.PackageGroup{}
	for _, name := range packages {
		group = append(group, syncer.Package(name))
	}

	config.Packages["Preinstalled"] = group
	config.Packages["MyPackages"] = syncer.PackageGroup{}

	encoder := toml.NewEncoder(f)
	encoder.Indent = "  "

	err = encoder.Encode(config)
	if err != nil {
		return nil, fmt.Errorf("encoding new config: %w", err)
	}

	err = finalizeConfig(filePath)
	if err != nil {
		return nil, fmt.Errorf("editing new config: %w", err)
	}

	return config, nil
}

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

	splittedConfig := strings.Split(config, "Preinstalled = ")

	splittedConfig[1] = strings.ReplaceAll(splittedConfig[1], "[", "[\n    ")
	splittedConfig[1] = strings.ReplaceAll(splittedConfig[1], ", ", ",\n    ")
	splittedConfig[1] = strings.ReplaceAll(splittedConfig[1], "]", "\n  ]")

	config = strings.Join(splittedConfig, "Preinstalled = ")

	err = os.WriteFile(filePath, []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("saving new config: %w", err)
	}

	return nil
}

func ReadConfig(pm package_manager.Pm) (config *Configuration, err error) {
	dir := getConfigDir()
	file := path.Join(dir, "config.toml")

	if _, err = os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("creating config directory: %w", err)
		}
	}

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		config, err = generateConfig(file, pm)
		if err != nil {
			return nil, fmt.Errorf("generating new config: %w", err)
		}

		return config, nil
	}

	_, err = toml.DecodeFile(file, &config)
	if err != nil {
		return nil, fmt.Errorf("config decoding: %w", err)
	}

	return config, nil
}
