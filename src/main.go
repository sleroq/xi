package main

import (
	"log"
	"xi/src/config"
	"xi/src/package_manager"
	"xi/src/syncer"
)

func main() {
	pm, err := package_manager.New()
	if err != nil {
		log.Fatalln("Getting package manager wrapper:", err)
	}
	configuration, err := config.ReadConfig(pm)
	if err != nil {
		log.Fatalln("Reading configuration", err)
	}

	err = syncer.SyncGroups(pm, configuration.Packages)
	if err != nil {
		log.Fatalln("Syncing packages", err)
	}
}
