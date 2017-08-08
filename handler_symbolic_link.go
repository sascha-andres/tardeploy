package tardeploy

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
)

func (configuration *Configuration) recreateWebSymbolicLink(application, versionPath string) error {
	var err error
	if err = configuration.removeWebSymbolicLink(application); err != nil {
		return err
	}
	if err = configuration.createWebSymbolicLink(application, versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not create new symbolic link for %s", application))
	}
	return err
}

func (configuration *Configuration) removeWebSymbolicLink(application string) error {
	symlinkPath := path.Join(configuration.Directories.WebRootDirectory, application)
	if _, err := os.Lstat(symlinkPath); err == nil {
		return os.Remove(symlinkPath)
	}
	return nil
}

func (configuration *Configuration) createWebSymbolicLink(application, versionPath string) error {
	return os.Symlink(versionPath, path.Join(configuration.Directories.WebRootDirectory, application))
}
