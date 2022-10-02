package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"xi/src/package_manager"
	"path"
	"xi/src/syncer"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	Packages []syncer.Package
}

func getConfigDir() string {
	configPath := os.Getenv("XDG_CONFIG_HOME")

	if configPath == "" {
		configPath = path.Join(os.Getenv("HOME"), ".config")
	}

	return path.Join(configPath, "xi")
}

func readConfig(pm package_manager.Pm) (config *Configuration, err error) {
	dir := getConfigDir()
	file := path.Join(dir, "config.toml")

	if _, err = os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("creating config directory: %w", err)
		}
	}

	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(file)
		if err != nil {
			return nil, fmt.Errorf("creating config file: %w", err)
		}
		defer f.Close()

		config := new(Configuration)

		packages, err := pm.GetInstalled()
		if err != nil {
			return nil, fmt.Errorf("getting list of installed packages: %w", err)
		}

		for _, name := range packages {
			config.Packages = append(config.Packages, syncer.Package(name))
		}

		encoder := toml.NewEncoder(f)

		err = encoder.Encode(config)
		if err != nil {
			return nil, fmt.Errorf("encoding new config: %w", err)
		}

		return config, nil
	}

	_, err = toml.DecodeFile(file, &config)
	if err != nil {
		return nil, fmt.Errorf("config decoding: %w", err)
	}

	return config, nil
}

func main() {
	pm, err := package_manager.New()
	if err != nil {
		log.Fatalln("Getting package manager wrapper:", err)
	}
	configuration, err := readConfig(pm)
	if err != nil {
		log.Fatalln("Reading configuration", err)
	}

	err = syncer.SyncPackages(pm, configuration.Packages)
	if err != nil {
		log.Fatalln("Syncing packages", err)
	}
}
