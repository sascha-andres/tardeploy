package tardeploy

import (
	"fmt"

	"github.com/prometheus/log"
)

// to handle gzip: compress/gzip
// to handle  tar: archive/tar

func (configuration *Configuration) SetupApplication(application string) error {
	log.Println(fmt.Sprintf("Setup for %s", application))

	// check and create app dir
	// extract tar within a subdirectory YYYMMDDHHmmss
	// delete symbolic link @ web/app
	// create symolic link @ web/app pointing to app/YYYMMDDHHmmss
	// remove directories that are not part of configured backup delay

	return nil
}

func (configuration *Configuration) RemoveApplication(application string) error {
	log.Println(fmt.Sprintf("Remove for %s", application))

	// delete symbolic link @ web/app
	// delete app dir

	return nil
}
