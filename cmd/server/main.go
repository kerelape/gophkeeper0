package main

import (
	"encoding/base64"
	"log"
	"os"
	"path"

	"github.com/kerelape/gophkeeper/cmd/server/config"
	"github.com/kerelape/gophkeeper/internal/server"
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
	var wd, wdError = os.Getwd()
	if wdError != nil {
		log.Fatalf(wdError.Error())
	}
	var gophkeeper = server.Server{
		RestAddress:       configuration.Rest.Address,
		RestUseTLS:        configuration.Rest.UseTLS,
		RestHostWhilelist: configuration.Rest.HostWhilelist,

		DatabaseDSN: configuration.DatabaseDSN,
		BlobsDir:    path.Join(wd, "blobs"),

		TokenSecret:   secret,
		TokenLifespan: configuration.Token.Lifespan,

		UsernameMinLength: configuration.UsernameMinLength,
		PasswordMinLength: configuration.PasswordMinLength,
	}
	runnable.Run(&gophkeeper)
}
