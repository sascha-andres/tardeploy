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

package trigger

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type (
	// TriggerConfiguration provides a place to configure triggers ( external programs ) called before or after a deployment
	Configuration struct {
		Before string
		After  string
	}
)

func (configuration *Configuration) BeforeRunTrigger(application string) error {
	return configuration.execute(application, "before")
}

func (configuration *Configuration) AfterRunTrigger(application string) error {
	return configuration.execute(application, "after")
}

func (configuration *Configuration) execute(application, status string) error {
	var cmd string
	switch status {
	case "before":
		if "" == configuration.Before {
			return nil
		}
		cmd = configuration.Before
		break
	case "after":
		if "" == configuration.After {
			return nil
		}
		cmd = configuration.After
		break
	}

	log.Infof("Trigger for %s: %s", application, status)

	command := exec.Command(cmd, application, status)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
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
