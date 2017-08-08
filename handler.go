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

	"time"

	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// SetupApplication deploys an application
func (configuration *Configuration) SetupApplication(tarball string) error {
	log.Infof("Setup for %s", tarball)

	var (
		application string
		err         error
	)

	if application, err = makeApplication(tarball); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not determine application name for %s", tarball))
	}

	if err := configuration.beforeRunTrigger(application); err != nil {
		return errors.Wrap(err, "Could not run trigger")
	}
	defer func() {
		if err := configuration.afterRunTrigger(application); err != nil {
			log.Warnf("Could not run after trigger: %s", err.Error())
		}
	}()

	timestamp := time.Now().Format("20060102150405")
	log.Debugf("Using timestamp %s", timestamp)
	versionPath := path.Join(configuration.Directories.ApplicationDirectory, application, timestamp)
	log.Infof("Deployment path: %s", versionPath)

	if err := configuration.ensureDirectories(application, versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not handle directories %s", application))
	}

	if err := configuration.ensureFiles(path.Join(configuration.Directories.TarballDirectory, tarball), versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not handle files %s", tarball))
	}

	if err := configuration.recreateWebSymbolicLink(application, versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not handle symbolic link for %s", application))
	}

	configuration.backup(application)

	return nil
}

// RemoveApplication removes an application from the server
func (configuration *Configuration) RemoveApplication(tarball string) error {
	log.Infof("Remove for %s", tarball)

	var (
		application string
		err         error
	)

	if application, err = makeApplication(tarball); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not determine application name for %s", tarball))
	}

	if err := configuration.removeWebSymbolicLink(application); err != nil {
		return err
	}

	return os.RemoveAll(path.Join(configuration.Directories.ApplicationDirectory, application))
}

func makeApplication(tarball string) (string, error) {
	if !strings.HasSuffix(tarball, ".tgz") && !strings.HasSuffix(tarball, ".tar.gz") {
		return "", errors.New("Only tar.gz or tgz allowed")
	}

	if strings.HasSuffix(tarball, ".tar.gz") {
		if len(tarball) == 7 {
			return "", errors.New("Expected at least one character as application name")
		}
		return tarball[0 : len(tarball)-7], nil
	}

	if len(tarball) == 4 {
		return "", errors.New("Expected at least one character as application name")
	}

	return tarball[0 : len(tarball)-4], nil
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
