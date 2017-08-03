package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sascha-andres/tardeploy"
)

func config() {
	configuration, err := tardeploy.LoadConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("TarballDirectory:     [%s]", configuration.TarballDirectory))
	log.Println(fmt.Sprintf("ApplicationDirectory: [%s]", configuration.ApplicationDirectory))
	log.Println(fmt.Sprintf("WebRootDirectory:     [%s]", configuration.WebRootDirectory))
	log.Println(fmt.Sprintf("WebOwner:             [%s]", configuration.WebOwner))

	mustExist("WebRootDirectory", configuration.WebRootDirectory)
	mustExist("ApplicationDirectory", configuration.ApplicationDirectory)
	mustExist("TarballDirectory", configuration.TarballDirectory)
}

func mustExist(name, path string) {
	ok, err := exists(path)
	if !ok || err != nil {
		log.Fatalf(fmt.Sprintf("%s [%s] does not exist", name, path))
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
