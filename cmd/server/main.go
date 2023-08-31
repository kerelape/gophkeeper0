package main

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/kerelape/gophkeeper/cmd/server/config"
	"github.com/kerelape/gophkeeper/internal/server"
	"github.com/kerelape/gophkeeper/internal/server/domain"
	"github.com/pior/runnable"
)

func main() {
	log.SetPrefix("[GOPHKEEPER] ")
	var configuration config.Config
	if err := config.Read(&configuration); err != nil {
		log.Println(configuration.Description())
		os.Exit(1)
	}
	var secret, decodeSecretError = base64.RawStdEncoding.DecodeString(configuration.Token.Secret)
	if decodeSecretError != nil {
		log.Fatalf("failed to parse token secret: %s", decodeSecretError.Error())
	}
	var gophkeeper = server.Server{
		RestAddress:  configuration.Rest.Address,
		RestCertFile: configuration.Rest.CertFile,
		RestKeyFile:  configuration.Rest.KeyFile,

		Repository: &domain.PostgresRepository{
			PasswordEncoding: base64.RawStdEncoding,

			DSN: configuration.DatabaseDSN,

			TokenSecret:   secret,
			TokenLifespan: configuration.Token.Lifespan,

			UsernameMinLength: configuration.UsernameMinLength,
			PasswordMinLength: configuration.PasswordMinLength,
		},
	}
	runnable.Run(&gophkeeper)
}
