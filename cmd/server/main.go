package main

import (
	"github.com/kerelape/gophkeeper/internal/server"
	"github.com/kerelape/gophkeeper/internal/server/rest"
	"github.com/pior/runnable"
)

func main() {
	var gophkeeper = server.Server{
		Rest: rest.Rest{
			Address:  ":8080",     // @todo #5 make Rest Address configurable.
			CertFile: "rest.cert", // @todo #5 make CertFile configurable.
			KeyFile:  "rest.key",  // @todo #5 make KeyFile configurable.
		},
	}
	runnable.Run(&gophkeeper)
}
