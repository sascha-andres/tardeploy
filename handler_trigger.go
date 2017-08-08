package tardeploy

import (
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
	log.Infof("Trigger for %s: %s", application, status)
	switch status {
	case "before":
		if "" == configuration.Trigger.Before {
			return nil
		}
		break
	case "after":
		if "" == configuration.Trigger.After {
			return nil
		}
		break
	}
	return errors.New("Trigger not yet implemented")
}
