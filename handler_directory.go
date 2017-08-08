package tardeploy

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (configuration *Configuration) ensureDirectories(application, versionPath string) error {
	if err := configuration.ensureAppDirectory(application); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not setup %s", application))
	}
	log.Infof("Ensuring path for timestamp (%s) exists", versionPath)
	if err := ensureDirectory(versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not create %s", versionPath))
	}
	return nil
}

func (configuration *Configuration) ensureAppDirectory(application string) error {
	log.Infof("Ensuring app directory for %s", application)
	return ensureDirectory(path.Join(configuration.Directories.ApplicationDirectory, application))
}

func ensureDirectory(directory string) error {
	log.Debugf("Ensuring directory %s", directory)
	if ok, err := exists(directory); !ok {
		log.Debugf("Creating %s", directory)
		err := os.MkdirAll(directory, 0750)
		if err != nil {
			return errors.Wrap(err, "Could not create directory")
		}
	} else {
		if ok && nil != err {
			return errors.Wrap(err, "Could not check if directory exists")
		}
	}
	return nil
}
