package main

import (
	"fmt"
	"os"

	"github.com/prometheus/log"

	"github.com/sascha-andres/tardeploy"
)

func config() *tardeploy.Configuration {
	cfg, err := tardeploy.LoadConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("TarballDirectory: [%s]", cfg.Directories.TarballDirectory))
	log.Println(fmt.Sprintf("ApplicationDirectory: [%s]", cfg.Directories.ApplicationDirectory))
	log.Println(fmt.Sprintf("WebRootDirectory: [%s]", cfg.Directories.WebRootDirectory))
	log.Println(fmt.Sprintf("WebOwner: [%s]", cfg.Directories.WebOwner))

	mustExist("WebRootDirectory", cfg.Directories.WebRootDirectory)
	mustExist("ApplicationDirectory", cfg.Directories.ApplicationDirectory)
	mustExist("TarballDirectory", cfg.Directories.TarballDirectory)

	return cfg
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
