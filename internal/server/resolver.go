package server

import (
	"github.com/SevvyP/plants/internal/db"
	"github.com/SevvyP/plants/internal/middleware"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
)

type Server struct {
	db db.DBInterface
}

func ResolveServer() *Server {
	return &Server{db: ResolveDB()}
}


func ResolveDB() db.DBInterface {
	return db.NewDB()
}

func (s *Server) Run() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(adapter.Wrap(middleware.EnsureValidToken()))
	r.GET("/v1/plant/:name", s.HandleGetPlant)
	r.POST("/v1/plant", s.HandleCreatePlant)
	r.PUT("/v1/plant", s.HandleUpdatePlant)
	r.DELETE("/v1/plant/:name", s.HandleDeletePlant)
	r.Run()
}
