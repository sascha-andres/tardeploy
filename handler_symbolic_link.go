package tardeploy

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
)

func (configuration *Configuration) recreateWebSymbolicLink(application, versionPath string) error {
	if err := configuration.removeWebSymbolicLink(application); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not remove old symbolic link for %s", application))
	}
	if err := configuration.createWebSymbolicLink(application, versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not create new symbolic link for %s", application))
	}
	return nil
}

func (configuration *Configuration) removeWebSymbolicLink(application string) error {
	if ok, err := exists(path.Join(configuration.Directories.WebRootDirectory, application)); !ok {
		return nil
	} else {
		if err != nil {
			return errors.Wrap(err, "Could not remove symbolic link to old version")
		}
	}
	return os.Remove(path.Join(configuration.Directories.WebRootDirectory, application))
}

func (configuration *Configuration) createWebSymbolicLink(application, versionPath string) error {
	return os.Symlink(versionPath, path.Join(configuration.Directories.WebRootDirectory, application))
}
