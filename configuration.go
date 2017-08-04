package tardeploy

import (
	"fmt"

	"github.com/prometheus/log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type (
	// Configuration contains the daemon config
	Configuration struct {
		Directories DirectoryConfiguration
		Application ApplicationHandling
	}

	// ApplicationHandling  provides information about handling older versions
	ApplicationHandling struct {
		NumberOfBackups int
	}

	// DirectoryConfiguration contains all data required to handle deployments
	DirectoryConfiguration struct {
		TarballDirectory     string // TarballDirectory denotes the place to put the tarballs
		WebRootDirectory     string // WebRootDirectory denotes the root for the web
		ApplicationDirectory string // ApplicationDirectory - where to store the applications
		WebOwner             string // Owner of the files
	}
)

func LoadConfiguration() (*Configuration, error) {
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Error reading config file")
	}

	// Confirm which config file is used
	log.Println(fmt.Sprintf("Using config:         [%s]", viper.ConfigFileUsed()))

	var config Configuration
	err := viper.Unmarshal(&config)

	return &config, err
}

func init() {
	viper.AddConfigPath("/etc/tardeploy")   // look in system config driectory
	viper.AddConfigPath("$HOME/.tardeploy") // maybe user space
	viper.AddConfigPath(".")                // local config
	viper.SetConfigName("tardeploy")        // file is named tardeploy.[yaml|json|toml]

	viper.AutomaticEnv() // read in environment variables that match
}
