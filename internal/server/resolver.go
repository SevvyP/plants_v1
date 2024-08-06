package server

import (
	"github.com/SevvyP/items/internal/db"
	"github.com/SevvyP/items/internal/middleware"
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
	r.GET("/v1/item/:name", s.HandleGetItem)
	r.POST("/v1/item", s.HandleCreateItem)
	r.PUT("/v1/item", s.HandleUpdateItem)
	r.DELETE("/v1/item/:name", s.HandleDeleteItem)
	r.Run()
}
