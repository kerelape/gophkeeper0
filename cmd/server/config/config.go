package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is configuration for gophkeeper's server.
type Config struct {
	Rest struct {
		Address       string   `env:"ADDRESS" env-default:":16355" env-description:"Address that REST api listens on."`
		UseTLS        bool     `env:"USE_TLS" env-default:"true" env-description:"Use TLS or not"`
		HostWhilelist []string `env:"HOST_WHITELIST" env-default:"" env-description:""`
	} `env-prefix:"REST_"`
	Token struct {
		Lifespan time.Duration `env:"LIFESPAN" env-description:"JWT Token lifespan in milliseconds" env-default:"15m"`
		Secret   string        `env:"SECRET" env-description:"Base64 encoded JWT Token secret" env-required:"true"`
	} `env-prefix:"TOKEN_"`
	UsernameMinLength uint   `env:"USERNAME_MIN_LENGTH" env-description:"Username minimum length" env-default:"0"`
	PasswordMinLength uint   `env:"PASSWORD_MIN_LENGTH" env-description:"Password minimum length" env-default:"0"`
	DatabaseDSN       string `env:"DATABASE_DSN" env-description:"Database connection URL" env-required:"true"`
}

// Read reads the config.
func Read(config *Config) error {
	if err := cleanenv.ReadEnv(config); err != nil {
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
