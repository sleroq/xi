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

	err = syncPackages(*configuration)
	if err != nil {
		log.Fatalln("Syncing packages", err)
	}

	fmt.Println("All synced!")
}
