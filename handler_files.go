package tardeploy

import (
	"fmt"
	"os"
	"strconv"

	"path/filepath"

	"os/user"

	"github.com/pkg/errors"
	"github.com/sascha-andres/tardeploy/deflate"
	log "github.com/sirupsen/logrus"
)

func (configuration *Configuration) ensureFiles(tarball, versionPath string) error {
	var (
		userID  int
		groupID int
		err     error
	)
	if err := deflate.Tarball(tarball, versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not deflate %s", tarball))
	}

	if userID, err = configuration.getUIDForUser(); err != nil {
		return err
	}
	if groupID, err = configuration.getGIDForGroup(); err != nil {
		return err
	}

	err = filepath.Walk(versionPath, func(path string, info os.FileInfo, err error) error {
		log.Debugf("Ownership for %s", path)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Could not set user rights for %s", path))
		}
		return os.Chown(path, userID, groupID)
	})
	if err != nil {
		return err
	}

	return nil
}

func (configuration *Configuration) getUIDForUser() (int, error) {
	var (
		value int
		err   error
		usr   *user.User
	)
	if value, err = strconv.Atoi(configuration.Directories.Security.User); err == nil {
		return value, nil
	}
	usr, err = user.Lookup(configuration.Directories.Security.User)
	if err != nil {
		return -1, errors.Wrap(err, fmt.Sprintf("Could not get userid for %s", configuration.Directories.Security.User))
	}
	if value, err = strconv.Atoi(usr.Uid); err == nil {
		return value, nil
	}
	return -1, errors.Wrap(err, fmt.Sprintf("Could not get uid for %s [returned value is %s]", configuration.Directories.Security.User, usr.Uid))
}

func (configuration *Configuration) getGIDForGroup() (int, error) {
	var (
		value int
		err   error
		grp   *user.Group
	)
	if value, err = strconv.Atoi(configuration.Directories.Security.Group); err == nil {
		return value, nil
	}
	grp, err = user.LookupGroup(configuration.Directories.Security.Group)
	if err != nil {
		return -1, errors.Wrap(err, fmt.Sprintf("Could not get userid for %s", configuration.Directories.Security.User))
	}
	if value, err = strconv.Atoi(grp.Gid); err == nil {
		return value, nil
	}
	return -1, errors.Wrap(err, fmt.Sprintf("Could not get gid for %s [returned value is %s]", configuration.Directories.Security.User, grp.Gid))
}
