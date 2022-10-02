package main

import (
	"log"
	"xi/configuration"
	"xi/src/syncer"
)

func main() {
	packages := configuration.Packages()
	err := syncer.Sync(packages)
	if err != nil {
		log.Fatalln("Syncing packages", err)
	}
}
