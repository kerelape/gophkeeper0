package main

import (
	"fmt"
	"os"

	"github.com/kerelape/gophkeeper/cmd/server/config"
	"github.com/kerelape/gophkeeper/internal/server"
	"github.com/pior/runnable"
)

func main() {
	var configuration config.Config
	if err := config.Read(&configuration); err != nil {
		fmt.Println(configuration.Description())
		os.Exit(1)
	}
	var gophkeeper = server.Server{
		RestAddress:  configuration.Rest.Address,
		RestCertFile: configuration.Rest.CertFile,
		RestKeyFile:  configuration.Rest.KeyFile,

		Repository: nil, // @todo #16 Implement identity repository.
	}
	runnable.Run(&gophkeeper)
}
