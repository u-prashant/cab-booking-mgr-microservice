package config

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// DefaultEnvVar is the default environment variable which points to the config file
const DefaultEnvVar = "BOOKING_MGR_CONFIG"

// App is the application config
var App *Config

// Config defines the yaml format for the config file
type Config struct {
	// DSN is the data source name
	DSN string

	// Address is the IP address and port to bind this REST to
	Address string
}

func init() {
	filename, found := os.LookupEnv(DefaultEnvVar)
	if !found {
		log.Errorf("failed to located file specified by %s", DefaultEnvVar)
		return
	}

	_ = load(filename)
}

func load(filename string) error {
	App = &Config{}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("failed to read config file. err: %s", err)
		return err
	}

	err = yaml.Unmarshal(bytes, App)
	if err != nil {
		log.Errorf("failed to parse config file. err: %s", err)
		return err
	}

	return nil
}
