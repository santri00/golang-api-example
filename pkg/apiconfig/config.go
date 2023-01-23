package apiconfig

import (
	"errors"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	_file "github.com/knadh/koanf/providers/file"
)

// Config defines the config struct
type Config struct {
	Store *koanf.Koanf
}

// New initializes the config
func New() *Config {
	return &Config{
		// Using . as the key path delimiter as default
		Store: koanf.New("."),
	}
}

// Load loads the configuration
func (cfg *Config) Load(file string) error {
	if len(file) == 0 {
		return errors.New("missing config file")
	}

	//Load the config file (YAML file) - a simple example
	if err := cfg.Store.Load(_file.Provider(file), yaml.Parser()); err != nil {
		return err
	}

	return nil
}
