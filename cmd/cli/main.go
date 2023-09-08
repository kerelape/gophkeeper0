package main

import (
	"flag"
	"log"

	"github.com/kerelape/gophkeeper/internal/cli"
	"github.com/pior/runnable"
)

func main() {
	var server = flag.String("s", "", "Gophkeeper address")
	flag.Parse()
	if server == nil {
		log.Fatal("Missing -s flag")
	}
}
