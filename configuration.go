package tardeploy

import (
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type (
	// Configuration contains the daemon config
	Configuration struct {
		Directories DirectoryConfiguration
		Application ApplicationHandling
		Trigger     TriggerConfiguration
	}

	// Trigger provides a place to configure triggers ( external programs ) called before or after a deployment
	TriggerConfiguration struct {
		Before string
		After  string
	}

	// ApplicationHandling  provides information about handling older versions
	ApplicationHandling struct {
		NumberOfBackups int    // How many old versions to keep
		BatchInterval   int    // How long to wait until the file changes are passed to tardeploy
		LogLevel        string // Log levels included debug -> info -> warn -> error
		TarCommand      string // Set this to a binary to execute external tar command
	}

	// FileSecurity defines the ownership of files
	FileSecurity struct {
		User  string // User or UID for file/directory owner
		Group string // Group or UID for file/directory owner
	}

	// DirectoryConfiguration contains all data required to handle deployments
	DirectoryConfiguration struct {
		TarballDirectory     string       // TarballDirectory denotes the place to put the tarballs
		WebRootDirectory     string       // WebRootDirectory denotes the root for the web
		ApplicationDirectory string       // ApplicationDirectory - where to store the applications
		Security             FileSecurity // Chown information
	}
)

func LoadConfiguration() (*Configuration, error) {
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "Error reading config file")
	}

	// Confirm which config file is used
	log.Debugf("Using config:         [%s]", viper.ConfigFileUsed())

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
