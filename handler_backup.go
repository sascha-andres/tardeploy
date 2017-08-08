package tardeploy

import (
	"path"

	"io/ioutil"

	"sort"

	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (configuration *Configuration) backup(application string) error {
	if configuration.Application.NumberOfBackups == -1 {
		return nil
	}

	appDirectory := path.Join(configuration.Directories.ApplicationDirectory, application)

	files, err := ioutil.ReadDir(appDirectory)
	if err != nil {
		return errors.Wrap(err, "Could not handle backups")
	}

	var directories []string
	for _, value := range files {
		if value.IsDir() {
			directories = append(directories, value.Name())
		}
	}
	if len(directories) <= configuration.Application.NumberOfBackups+1 {
		return nil
	}
	sort.Strings(directories)
	directoriesToRemove := directories[configuration.Application.NumberOfBackups+1:]

	for _, value := range directoriesToRemove {
		if err := os.RemoveAll(path.Join(appDirectory, value)); err != nil {
			log.Warnf("Could not remove %s: %#v", value, err)
		}
	}

	return nil
}
