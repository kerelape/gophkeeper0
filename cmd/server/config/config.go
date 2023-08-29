package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is configuration for gophkeeper's server.
type Config struct {
	Rest struct {
		Address  string `env:"ADDRESS" env-default:":16355" env-description:"Address that REST api listens on."`
		CertFile string `env:"CERT_FILE" env-required:"true" env-description:"Path to HTTPS cert file for REST api."`
		KeyFile  string `env:"KEY_FILE" env-required:"true" env-description:"Path to HTTPS key file for REST api."`
	} `env-prefix:"REST_"`
}

// Read reads the config.
func Read(config *Config) error {
	if err := cleanenv.ReadEnv(&config); err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	return nil
}

// Description returns config description.
func (c *Config) Description() string {
	var description, err = cleanenv.GetDescription(c, nil)
	if err != nil {
		panic(err)
	}
	return description
}
