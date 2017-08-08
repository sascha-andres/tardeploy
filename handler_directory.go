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
