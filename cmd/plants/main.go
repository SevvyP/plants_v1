package main

import (
	"github.com/SevvyP/plants/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {
	server := server.ResolveServer()
	r := gin.Default()
	r.StaticFile("/favicon.ico", "../public/favicon.ico")
	v1 := r.Group("/v1")
	{
		v1.GET("/plant/:name", server.HandleGetPlant)
	}
	r.Run()
}
