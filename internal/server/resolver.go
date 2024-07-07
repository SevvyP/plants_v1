package server

import "fmt"

type Server struct {
	db *DB
}

func ResolveServer() *Server {
	db := ResolveDB()
	return &Server{db: db}
	
}

func ResolveDB() *DB {
	return NewDB()
}

func (s *Server) Start() {
	fmt.Println("hello")
}