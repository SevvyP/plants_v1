package server

import (
	"github.com/SevvyP/plants/internal/db"
	"github.com/gin-gonic/gin"
)

type Server struct {
	db *db.DB
}

func ResolveServer() *Server {
	return &Server{db: ResolveDB()}
}


func ResolveDB() *db.DB {
	return db.NewDB()
}

func (s *Server) Run() {
	r := gin.Default()
	r.GET("/v1/plant/:name", s.HandleGetPlant)
	r.POST("/v1/plant", s.HandleCreatePlant)
	r.PUT("/v1/plant", s.HandleUpdatePlant)
	r.DELETE("/v1/plant/:name", s.HandleDeletePlant)
	r.Run()
}
