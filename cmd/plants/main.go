package main

import (
	"github.com/SevvyP/plants/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// load env file if one exists
	godotenv.Load()
	
	server := server.ResolveServer()
	server.Run()
}
