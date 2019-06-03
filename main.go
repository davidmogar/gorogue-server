package main

import (
	"gorogue-server/server"
	"log"
)

func main() {
	httpServer := &server.Server{}

	log.Fatal(httpServer.Start())
}