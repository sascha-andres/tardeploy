package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/sascha-andres/tardeploy"
)

func config() *tardeploy.Configuration {
	cfg, err := tardeploy.LoadConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("TarballDirectory: [%s]", cfg.Directories.TarballDirectory)
	log.Infof("ApplicationDirectory: [%s]", cfg.Directories.ApplicationDirectory)
	log.Infof("WebRootDirectory: [%s]", cfg.Directories.WebRootDirectory)
	log.Infof("WebOwner: [%s:%s]", cfg.Directories.Security.User, cfg.Directories.Security.Group)

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
