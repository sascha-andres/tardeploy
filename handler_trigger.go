package tardeploy

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (configuration *Configuration) beforeRunTrigger(application string) error {
	return configuration.trigger(application, "before")
}

func (configuration *Configuration) afterRunTrigger(application string) error {
	return configuration.trigger(application, "after")
}

func (configuration *Configuration) trigger(application, status string) error {
	var cmd string
	switch status {
	case "before":
		if "" == configuration.Trigger.Before {
			return nil
		}
		cmd = configuration.Trigger.Before
		break
	case "after":
		if "" == configuration.Trigger.After {
			return nil
		}
		cmd = configuration.Trigger.After
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
