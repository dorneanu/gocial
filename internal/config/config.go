package config

import (
	"io/ioutil"

	"github.com/dorneanu/gomation/internal/entity"
	"gopkg.in/yaml.v3"
)

// Config is a minimal configuration for this utility
type Config struct {
	ServerPort int                       `yaml:"server_port"`
	Identities []entity.IdentityProvider `yaml:"identities"`
	JWT        JWTConfig                 `yaml:"jwt_config"`
}

type JWTConfig struct {
	Secret    string `yaml:"secret"`
	Algorithm string `yaml:"algorithm"`
}

func Load(file string) (*Config, error) {
	c := Config{}

	// Load config from YAML file
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
