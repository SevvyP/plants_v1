package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) HandleGetPlant(c *gin.Context) {
	plant, err := s.db.GetPlant(c.Param("name"))
	if err != nil {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, plant)
}
