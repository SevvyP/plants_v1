package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SevvyP/plants/internal/db"
	"github.com/SevvyP/plants/pkg"
	"github.com/gin-gonic/gin"
)

func (s *Server) HandleCreatePlant(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var plant pkg.Plant
    err := decoder.Decode(&plant)
    if err != nil {
        log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
    }
	if plant.Name == "" || plant.Description == "" {
		log.Println("create request missing name or description")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.db.CreatePlant(plant, c)
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, plant)
}

func (s *Server) HandleGetPlant(c *gin.Context) {
	if c.Param("name") == "" {
		log.Println("create request missing name or description")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	plant, err := s.db.GetPlant(c.Param("name"), c)
	if err != nil {
		if err.Error() == db.ErrNotFound {
			log.Println(err)
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, plant)
}

func (s *Server) HandleUpdatePlant(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var plant pkg.Plant
    err := decoder.Decode(&plant)
    if err != nil {
        log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
    }
	if plant.Name == "" || plant.Description == "" {
		log.Println("create request missing name or description")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.db.UpdatePlant(plant, c)
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, plant)
}

func (s *Server) HandleDeletePlant(c *gin.Context) {
	if c.Param("name") == "" {
		log.Println("create request missing name or description")
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	plant, err := s.db.DeletePlant(c.Param("name"), c)
	if err != nil {
		if err.Error() == db.ErrNotFound {
			log.Println(err)
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, plant)
}