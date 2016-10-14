package orca

import (
	"fmt"
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
	filepath.Join(getUserDir(), ".orca", "config.yml"),
	filepath.Join(getCwd(), "orca.yml"),
}

// LoadConfig the configuration from the ConfigPath
func LoadConfig() *Config {
	config := new(Config)

	for _, path := range ConfigPath {
		config.Read(path)
	}

	return config
}

// Config is read from a YAML file and defines the current configuration of
// the project and can be exported as such.
type Config struct {
	Debug  bool
	Name   string
	Addr   string
	Domain string
	DBPath string `yaml:"dbpath"`
}

// Parse configuration from data
func (conf *Config) Parse(data []byte) error {
	// Unmarshall the YAML data into the config
	if err := yaml.Unmarshal(data, conf); err != nil {
		return err
	}

	// Perform validation and set defaults.
	if conf.Name == "" {

		// If no name specified, use the hostname of the machine
		name, err := os.Hostname()

		if err != nil {
			// If this lookup fails, specify unknown device
			name = "unknown device"
		}
		conf.Name = name
	}

	if conf.DBPath == "" {
		// If no database path specified, store in the home directory
		conf.DBPath = filepath.Join(getUserDir(), ".orca", "orca.db")
	}

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

// String returns a string representation of the configuration
func (conf Config) String() string {
	output := fmt.Sprintf("%s configuration (debug = %t)", conf.Name, conf.Debug)

	addr, err := ResolveAddr(conf.Addr)
	if err != nil {
		addr = "no resolved IP address!"
	}
	output += fmt.Sprintf("\nIP Address: %s", addr)

	if conf.Domain != "" {
		output += fmt.Sprintf(" | Domain: %s", conf.Domain)
	}

	output += fmt.Sprintf("\nDatabase: %s", conf.DBPath)

	return output
}
