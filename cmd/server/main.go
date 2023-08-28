package main

import (
	"github.com/kerelape/gophkeeper/cmd/server/config"
	"github.com/kerelape/gophkeeper/internal/server"
	"github.com/kerelape/gophkeeper/internal/server/rest"
	"github.com/pior/runnable"
)

func main() {
	var configuration config.Config
	if err := config.Read(&configuration); err != nil {
		panic(err)
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
