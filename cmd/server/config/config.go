package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is configuration for gophkeeper's server.
type Config struct {
	Rest struct {
		Address  string `env:"ADDRESS" env-default:":16355"`
		CertFile string `env:"CERT_FILE" env-required:"true"`
		KeyFile  string `env:"KEY_FILE" env-required:"true"`
	} `env-prefix:"REST_"`
}

// Read reads the config.
func Read(config *Config) error {
	if err := cleanenv.ReadEnv(&config); err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	return nil
}
