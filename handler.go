package tardeploy

import (
	"fmt"
	"os"

	"path"

	"time"

	"strings"

	"github.com/pkg/errors"
	"github.com/prometheus/log"
	"github.com/sascha-andres/tardeploy/deflate"
)

// to handle gzip: compress/gzip
// to handle  tar: archive/tar

func (configuration *Configuration) SetupApplication(tarball string) error {
	log.Println(fmt.Sprintf("Setup for %s", tarball))

	var (
		application string
		err         error
	)

	if application, err = makeApplication(tarball); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not determine application name for %s", tarball))
	}

	timestamp := time.Now().Format("20060102150405")
	log.Debugf("Using timestamp %s", timestamp)
	versionPath := path.Join(configuration.Directories.ApplicationDirectory, application, timestamp)

	if err := configuration.ensureDirectories(application, versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not handle directories %s", application))
	}

	// TODO ensure files
	if err := deflate.DeflateTarball(path.Join(configuration.Directories.TarballDirectory, tarball), versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not deflate %s", tarball))
	}
	// chown

	if err := configuration.recreateWebSymbolicLink(application, versionPath); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not handle symbolic link for %s", application))
	}

	// remove directories that are not part of configured backup delay

	return nil
}

func (configuration *Configuration) RemoveApplication(tarball string) error {
	log.Println(fmt.Sprintf("Remove for %s", tarball))

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

// helper

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
