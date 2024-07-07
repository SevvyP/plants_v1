package server

import (
	"github.com/SevvyP/plants/internal/db"
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
