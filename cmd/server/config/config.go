package config

// Config is configuration for gophkeeper's server.
type Config struct {
	Rest struct {
		Address  string `env:"ADDRESS" env-required:"true"`
		CertFile string `env:"CERT_FILE" env-required:"true"`
		KeyFile  string `env:"KEY_FILE" env-required:"true"`
	} `env-prefix:"REST_"`
}
