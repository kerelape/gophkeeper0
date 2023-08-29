package main

import (
	"fmt"
	"os"

	"github.com/kerelape/gophkeeper/cmd/server/config"
	"github.com/kerelape/gophkeeper/internal/server"
	"github.com/kerelape/gophkeeper/internal/server/rest"
	"github.com/pior/runnable"
)

func main() {
	var configuration config.Config
	if err := config.Read(&configuration); err != nil {
		fmt.Println(configuration.Description())
		os.Exit(1)
	}
	var gophkeeper = server.Server{
		Rest: rest.Rest{
			Address:  configuration.Rest.Address,
			CertFile: configuration.Rest.CertFile,
			KeyFile:  configuration.Rest.KeyFile,
		},
	}
	runnable.Run(&gophkeeper)
}
