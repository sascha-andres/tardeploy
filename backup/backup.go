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

package backup

import (
	"path"

	"io/ioutil"

	"sort"

	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Execute(appDirectory string, numberOfBackups int) error {
	if numberOfBackups < 0 {
		return nil
	}

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
	if len(directories) <= numberOfBackups+1 {
		return nil
	}
	sort.Sort(sort.Reverse(sort.StringSlice(directories)))

	directoriesToRemove := directories[numberOfBackups+1:]

	for _, value := range directoriesToRemove {
		deploymentDirectory := path.Join(appDirectory, value)
		if err := os.RemoveAll(deploymentDirectory); err != nil {
			log.Warnf("Could not remove %s: %#v", value, err)
		} else {
			log.Infof("Removed old deployment directory %s", deploymentDirectory)
		}
	}

	return nil
}
