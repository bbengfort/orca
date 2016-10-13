package orca

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getUserDir() string {
	// HACK Ignore errors!?!
	usr, _ := user.Current()
	return usr.HomeDir
}

func getCwd() string {
	// HACK Ignore errors!?!
	cwd, _ := os.Getwd()
	return cwd
}

// ConfigPath specifies the locations to look up configurations
var ConfigPath = []string{
	"/etc/orca.yml",
	filepath.Join(getUserDir(), ".orca.yml"),
	filepath.Join(getCwd(), "orca.yml"),
	filepath.Join(getCwd(), "conf", "orca.yml"),
}

// Config is read from a YAML file and defines the current configuration of
// the project and can be exported as such.
type Config struct {
	Debug  bool   `default:"true"`
	Name   string `required:"true"`
	Addr   string `required:"true"`
	Domain string
}

// Parse configuration from data
func (conf *Config) Parse(data []byte) error {
	// Unmarshall the YAML data into the config
	if err := yaml.Unmarshal(data, conf); err != nil {
		return err
	}

	// Perform validation and set defaults.
	// TODO: Use reflection to validate the YAML and only set new values

	// Return nil if there was no error
	return nil
}

// Read a configuration from a path
func (conf *Config) Read(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = conf.Parse(data); err != nil {
		return err
	}

	return nil
}

// LoadConfig the configuration from the ConfigPath
func LoadConfig() *Config {
	config := new(Config)

	for _, path := range ConfigPath {
		config.Read(path)
	}

	return config
}
