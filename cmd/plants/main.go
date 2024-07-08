package main

import (
	"fmt"

	"github.com/SevvyP/plants/internal/server"
)

func main() {
	fmt.Println("hello from main")
	server := server.ResolveServer()
	server.Run()
}
