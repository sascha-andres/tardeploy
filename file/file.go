// Copyright © 2017 Sascha Andres <sascha.andres@outlook.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package file

import (
	"fmt"
	"os"
	"strconv"

	"path/filepath"

	"os/exec"
	"os/user"

	"github.com/pkg/errors"
	"github.com/sascha-andres/tardeploy/deflate"
	log "github.com/sirupsen/logrus"
)

func callTarCommand(tarcommand, tarball, directory string) error {
	log.Infof("Calling '%s xzf %s' in %s", tarcommand, tarball, directory)
	command := exec.Command(tarcommand, "xzf", tarball)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Dir = directory
	var err error
	if err = command.Start(); err != nil {
		return errors.Wrap(err, "could not start command")
	}
	err = command.Wait()
	if err != nil {
		return errors.Wrap(err, "Could not wait for command")
	}

	return nil
}

func Ensure(tarCommand, user, group, tarball, versionPath string) error {
	var (
		userID  int
		groupID int
		err     error
	)
	if tarCommand == "" {
		if err := deflate.Tarball(tarball, versionPath); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Could not deflate %s", tarball))
		}
	} else {
		if err := callTarCommand(tarCommand, tarball, versionPath); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Could not exec tar command for %s", tarball))
		}
	}

	if userID, err = getUIDForUser(user); err != nil {
		return err
	}
	if groupID, err = getGIDForGroup(group); err != nil {
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

func getUIDForUser(userName string) (int, error) {
	var (
		value int
		err   error
		usr   *user.User
	)
	if value, err = strconv.Atoi(userName); err == nil {
		return value, nil
	}
	usr, err = user.Lookup(userName)
	if err != nil {
		return -1, errors.Wrap(err, fmt.Sprintf("Could not get userid for %s", userName))
	}
	if value, err = strconv.Atoi(usr.Uid); err == nil {
		return value, nil
	}
	return -1, errors.Wrap(err, fmt.Sprintf("Could not get uid for %s [returned value is %s]", userName, usr.Uid))
}

func getGIDForGroup(groupName string) (int, error) {
	var (
		value int
		err   error
		grp   *user.Group
	)
	if value, err = strconv.Atoi(groupName); err == nil {
		return value, nil
	}
	grp, err = user.LookupGroup(groupName)
	if err != nil {
		return -1, errors.Wrap(err, fmt.Sprintf("Could not get userid for %s", groupName))
	}
	if value, err = strconv.Atoi(grp.Gid); err == nil {
		return value, nil
	}
	return -1, errors.Wrap(err, fmt.Sprintf("Could not get gid for %s [returned value is %s]", groupName, grp.Gid))
}
