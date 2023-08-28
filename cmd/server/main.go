package main

import (
	"github.com/kerelape/gophkeeper/internal/server"
	"github.com/pior/runnable"
)

func main() {
	var gophkeeper = server.MakeServer()
	runnable.Run(&gophkeeper)
}
