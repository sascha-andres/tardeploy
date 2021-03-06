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

package symlink

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Recreate(webRootDirectory, application, versionPath string) error {
	var err error
	if err = Remove(webRootDirectory, application); err != nil {
		return err
	}
	if err = Create(webRootDirectory, application, versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not create new symbolic link for %s", application))
	}
	return err
}

func Remove(webRootDirectory, application string) error {
	symlinkPath := path.Join(webRootDirectory, application)
	if _, err := os.Lstat(symlinkPath); err == nil {
		log.Debugf("Removing symbolic link %s")
		return os.Remove(symlinkPath)
	}
	return nil
}

func Create(webRootDirectory, application, versionPath string) error {
	deploymentDirectory := path.Join(webRootDirectory, application)
	log.Debugf("Link from %s to %s", versionPath, deploymentDirectory)
	return os.Symlink(versionPath, deploymentDirectory)
}
