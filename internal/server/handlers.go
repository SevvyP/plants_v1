package server

import (
	"encoding/json"
	"errors"
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
	err = s.db.CreatePlant(plant, c)
	if err != nil {
		if errors.Is(err, err.(*db.ErrNotFound)) {
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

func (s *Server) HandleGetPlant(c *gin.Context) {
	plant, err := s.db.GetPlant(c.Param("name"), c)
	if err != nil {
		if errors.Is(err, db.ErrNotFound{}) {
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
	err = s.db.UpdatePlant(plant, c)
	if err != nil {
		log.Println(err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, plant)
}

func (s *Server) HandleDeletePlant(c *gin.Context) {
	plant, err := s.db.DeletePlant(c.Param("name"), c)
	if err != nil {
		if errors.Is(err, db.ErrNotFound{}) {
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