package main

import "github.com/medibloc/panacea-data-market-validator/server"

func main() {
	// TODO: graceful shutdown
	server.Run()
}
