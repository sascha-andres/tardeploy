package tardeploy

import (
	"fmt"

	"github.com/prometheus/log"
)

// to handle gzip: compress/gzip
// to handle  tar: archive/tar

func (configuration *Configuration) SetupApplication(application string) error {
	log.Println(fmt.Sprintf("Setup for %s", application))
	return nil
}

func (configuration *Configuration) RemoveApplication(application string) error {
	log.Println(fmt.Sprintf("Remove for %s", application))
	return nil
}
