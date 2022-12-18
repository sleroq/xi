package main

import (
	"fmt"
	"log"
)

func main() {
	configuration, err := readConfig()
	if err != nil {
		log.Fatalln("Reading configuration", err)
	}

	fmt.Println(configuration.Managers)
	// err = syncer.SyncGroups(pm, configuration.Packages)
	// if err != nil {
	// 	log.Fatalln("Syncing packages", err)
	// }
}
