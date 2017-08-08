package tardeploy

import "github.com/pkg/errors"

func (configuration *Configuration) beforeRunTrigger(application string) error {
	return configuration.trigger(application, "before")
}

func (configuration *Configuration) afterRunTrigger(application string) error {
	return configuration.trigger(application, "after")
}

func (configuration *Configuration) trigger(application, status string) error {
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
