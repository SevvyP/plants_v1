package main

import (
	"github.com/SevvyP/plants/internal/server"
)

func main() {
	server := server.ResolveServer()
	server.Run()
}
